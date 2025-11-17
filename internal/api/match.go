package api

import (
	"log"
	"net/http"

	"github.com/google/uuid"

	"hearts/internal/agg"
	"hearts/internal/db"
)

func MatchStateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	matchId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "could not parse match id", http.StatusBadRequest)
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

	result, err := agg.GetAggregate(events)
	if err != nil {
		log.Println("ERROR: failed to query match events:", err)
		http.Error(w, "failed to query match", http.StatusInternalServerError)
		return
	}
	if result == nil {
		http.Error(w, "not found", http.StatusNotFound)
	}

	WriteResponse(w, result)
}
