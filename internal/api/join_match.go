package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"

	"hearts/internal/agg"
	"hearts/internal/db"
)

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

func JoinMatchHandler(w http.ResponseWriter, r *http.Request) {
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

	dbConn, err := GetDatabaseOrFail(r.Context(), w)
	if err != nil || dbConn == nil {
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

	WriteResponse(w, agg)
}
