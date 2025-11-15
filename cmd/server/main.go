package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"

	"hearts/internal/agg"
	"hearts/internal/api"
	"hearts/internal/db"
)

var (
	dbConn *sql.DB
	err error
)

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
		log.Println("ERROR: failed to query match events:", err)
		http.Error(w, "failed to query match", http.StatusInternalServerError)
		return
	}

	result, err := agg.GetAggregate(events)
	if err != nil {
		log.Println("ERROR: failed to query match events:", err)
		http.Error(w, "failed to query match", http.StatusInternalServerError)
		return
	}
	if result == nil {
		http.Error(w, "not found", http.StatusNotFound)
	}

	api.WriteResponse(w, result)
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
		log.Println("ERROR: failed to create match:", err)
		http.Error(w, "failed to create match", http.StatusInternalServerError)
		return
	}
	log.Println("INFO: created match with id", id)

	api.WriteResponse(w, CreateMatchResponse{MatchId: id})
}

type JoinMatchRequest struct {
	UserId uuid.UUID `json:"user_id"`
	MatchId uuid.UUID `json:"match_id"`
}

func (r *JoinMatchRequest) CreateEvent(aggVersion uint32) (*agg.MatchEvent, error) {
	payload, err := json.Marshal(agg.JoinMatchPayload{UserId: r.UserId})
	if err != nil {
		return nil, err
	}
	return &agg.MatchEvent{
		Type: "player-joined",
		AggregateVersion: aggVersion,
		CreatedOn: agg.Timestamp(),
		Payload: payload,
	}, nil
}

func joinMatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return 
	}

	var request JoinMatchRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	events, err := db.QueryMatchEvents(dbConn, request.MatchId)
	if err != nil {
		log.Println("ERROR: failed to query match events:", err)
		http.Error(w, "failed to query match", http.StatusInternalServerError)
		return
	}

	agg, err := agg.GetAggregate(events)
	if err != nil {
		log.Println("ERROR: failed to query match events:", err)
		http.Error(w, "failed to query match", http.StatusInternalServerError)
		return
	}
	if agg == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	event, err := request.CreateEvent(agg.Version+1)
	if err != nil {
		log.Println("ERROR: failed to marshal event:", err)
		http.Error(w, "failed to marshal event", http.StatusInternalServerError)
		return
	}

	err = agg.ApplyEvent(event)
	if err != nil {
		http.Error(w, "failed to join match", http.StatusBadRequest)
		return
	}

	err = db.InsertEvent(dbConn, request.MatchId, event)
	if err != nil {
		log.Println("ERROR: failed to join match", err)
		http.Error(w, "failed to join match", http.StatusInternalServerError)
		return
	}

	api.WriteResponse(w, agg)
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

	http.Handle("/health", api.CommonHeaders(http.HandlerFunc(api.HealthHandler)))
	http.Handle("/match/{id}", api.CommonHeaders(http.HandlerFunc(matchStateHandler)))
	http.Handle("/create-match", api.CommonHeaders(http.HandlerFunc(createMatchHandler)))
	http.Handle("/join-match", api.CommonHeaders(http.HandlerFunc(joinMatchHandler)))

	fmt.Println("listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
