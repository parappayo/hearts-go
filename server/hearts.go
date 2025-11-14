package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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
	table := game.Table{}
	table.AddSeats(4)

	deck := cards.CreateDeck()
	deck.Shuffle()
	table.Deal(deck)

	// TODO: we're serving the entire game state but we need to filter it for the requesting player's view
	writeResponse(w, table)
}

type CreateMatchResponse struct {
	MatchId string `json:"match_id"`
}

func createMatchHandler(w http.ResponseWriter, r *http.Request) {
	id, err := db.CreateMatch(dbConn)
	if err != nil {
		log.Println("ERROR: failed to create match", err)
		http.Error(w, "failed to create match", http.StatusInternalServerError)
		return
	}
	log.Println("INFO: created match with id", id)

	writeResponse(w, CreateMatchResponse{MatchId: id.String()})
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
	http.Handle("/create-match", commonHeaders(http.HandlerFunc(createMatchHandler)))

	fmt.Println("listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
