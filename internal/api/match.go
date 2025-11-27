package api

import (
	"log"
	"net/http"

	"github.com/google/uuid"

	"hearts/internal/agg"
	"hearts/internal/db"
	"hearts/pkg/cards"
	"hearts/pkg/game"
)

type MatchStateResponse struct {
	Hand cards.Hand `json:"Hand"`
	CurrentPlayersTurn int `json:"CurrentPlayersTurn"`
	Seat int `json:"Seat"`
	Trick game.Trick `json:"Trick"`
}

func MatchStateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	matchId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "match id not provided", http.StatusBadRequest)
		return
	}

	userId, err := uuid.Parse(r.Header.Get("User-Id"))
	if err != nil {
		http.Error(w, "user id not provided", http.StatusBadRequest)
		return
	}

	dbConn, err := GetDatabaseOrFail(r.Context(), w)
	if err != nil || dbConn == nil {
		return
	}

	events, err := db.QueryMatchEvents(dbConn, matchId)
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
	}

	seat, _, tablePlayer := agg.FindPlayer(userId)
	response := MatchStateResponse{
		Hand: *tablePlayer.Hand,
		CurrentPlayersTurn: agg.Table.CurrentPlayersTurn,
		Seat: seat,
		Trick: agg.Table.CurrentTrick(),
	}
	WriteResponse(w, response)
}
