package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"

	"hearts/internal/agg"
	"hearts/internal/db"
	"hearts/pkg/cards"
	"hearts/pkg/game"
)

type StartMatchRequest struct {
	MatchId uuid.UUID `json:"match_id"`
}

func (r *StartMatchRequest) CreateEvent(aggVersion uint32) (*agg.MatchEvent, error) {
	table := game.Table{}
	table.AddSeats(4)
	deck := cards.CreateDeck()
	deck.Shuffle()
	table.Deal(deck)

	payload, err := json.Marshal(agg.StartMatchPayload{
		CurrentPlayersTurn: table.CurrentPlayersTurn,
		Hands: table.Hands(),
	})
	if err != nil {
		return nil, err
	}
	return &agg.MatchEvent{
		Type: "match-started",
		AggregateVersion: aggVersion,
		CreatedOn: agg.Timestamp(),
		Payload: payload,
	}, nil
}

func StartMatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return 
	}

	var request StartMatchRequest
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

	if len(agg.Players) != 4 {
		http.Error(w, "cannot start match unless there are four players", http.StatusBadRequest)
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
		http.Error(w, "failed to start match", http.StatusBadRequest)
		log.Println("ERROR: failed to start match", err)
		return
	}

	err = db.InsertEvent(dbConn, request.MatchId, event)
	if err != nil {
		log.Println("ERROR: failed to start match", err)
		http.Error(w, "failed to start match", http.StatusInternalServerError)
		return
	}

	WriteResponse(w, agg)
}
