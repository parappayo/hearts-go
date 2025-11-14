package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"

	"hearts/cards"
	"hearts/game"
	"hearts/db"
)

var (
	dbConn *sql.DB
	err error
)

func writeResponse[T any](w http.ResponseWriter, data T) {
	body, err := json.Marshal(data)
	if err != nil {
		log.Println("ERROR: failed to marshal response", err)
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(body))
}

type HealthResponse struct {
	Status string `json:"status"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	writeResponse(w, HealthResponse{Status: "pass"})
}

func gameStateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	table := game.Table{}
	table.AddSeats(4)

	deck := cards.CreateDeck()
	deck.Shuffle()
	table.Deal(deck)

	// TODO: we're serving the entire game state but we need to filter it for the requesting player's view
	writeResponse(w, table)
}

func matchStateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	matchId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "could not parse match id", http.StatusBadRequest)
		return
	}

	events, err := db.QueryMatchEvents(dbConn, matchId)
	if err != nil {
		log.Println("ERROR: failed to query match events", err)
		http.Error(w, "failed to query match", http.StatusInternalServerError)
		return
	}

	result, err := db.GetAggregate(events)
	if err != nil {
		log.Println("ERROR: failed to query match events", err)
		http.Error(w, "failed to query match", http.StatusInternalServerError)
		return
	}

	writeResponse(w, result)
}

type CreateMatchResponse struct {
	MatchId uuid.UUID `json:"match_id"`
}

func createMatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := db.CreateMatch(dbConn)
	if err != nil {
		log.Println("ERROR: failed to create match", err)
		http.Error(w, "failed to create match", http.StatusInternalServerError)
		return
	}
	log.Println("INFO: created match with id", id)

	writeResponse(w, CreateMatchResponse{MatchId: id})
}

type JoinMatchRequest struct {
	UserId uuid.UUID `json:"user_id"`
	MatchId uuid.UUID `json:"match_id"`
}

func joinMatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return 
	}

	var request JoinMatchRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	// TODO: read the match aggregate and get the actual next version
	version := uint32(2)

	err = db.JoinMatch(dbConn, request.UserId, request.MatchId, version+1)
	if err != nil {
		log.Println("ERROR: failed to join match", err)
		http.Error(w, "failed to join match", http.StatusInternalServerError)
		return
	}

	// TODO: return the new match aggregate
	writeResponse(w, db.MatchState{})
}

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func main() {
	dbConn, err = db.Open(os.Getenv("DB_CONN"))
	if err != nil {
		panic(err)
	}
	err = db.CreateSchema(dbConn)
	if err != nil {
		panic(err)
	}

	http.Handle("/health", commonHeaders(http.HandlerFunc(healthHandler)))
	http.Handle("/game-state", commonHeaders(http.HandlerFunc(gameStateHandler)))
	http.Handle("/match/{id}", commonHeaders(http.HandlerFunc(matchStateHandler)))
	http.Handle("/create-match", commonHeaders(http.HandlerFunc(createMatchHandler)))

	fmt.Println("listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
