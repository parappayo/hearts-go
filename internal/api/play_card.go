package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"

	"hearts/internal/agg"
	"hearts/internal/db"
	"hearts/pkg/cards"
)

type PlayCardRequest struct {
	MatchId uuid.UUID `json:"match_id"`
	PlayerId uuid.UUID `json:"player_id"`
	Rank string `json:"rank"`
	Suit string `json:"suit"`
}

func (r *PlayCardRequest) CreateEvent(aggVersion uint32) (*agg.MatchEvent, error) {
	card := cards.Card{}
	card.Rank = r.Rank
	card.Suit = r.Suit

	payload, err := json.Marshal(agg.PlayCardPayload{
		Card: card,
	})
	if err != nil {
		return nil, err
	}
	return &agg.MatchEvent{
		Type: "card-played",
		AggregateVersion: aggVersion,
		CreatedOn: agg.Timestamp(),
		Payload: payload,
	}, nil
}

func PlayCardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return 
	}

	var request PlayCardRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	dbConn, err := GetDatabaseOrFail(r.Context(), w)
	if err != nil || dbConn == nil {
		return
	}

	agg, err := QueryAggregate(dbConn, request.MatchId, w)
	if err != nil || agg == nil {
		return
	}

	// TODO: check current player's turn against request player uuid

	event, err := request.CreateEvent(agg.Version+1)
	if err != nil {
		log.Println("ERROR: failed to marshal event:", err)
		http.Error(w, "failed to marshal event", http.StatusInternalServerError)
		return
	}

	err = agg.ApplyEvent(event)
	if err != nil {
		http.Error(w, "failed to play card", http.StatusBadRequest)
		log.Println("ERROR: failed to play card:", err)
		return
	}

	err = db.InsertEvent(dbConn, request.MatchId, event)
	if err != nil {
		http.Error(w, "failed to play card", http.StatusInternalServerError)
		log.Println("ERROR: failed to play card:", err)
		return
	}

	WriteResponse(w, agg)
}
