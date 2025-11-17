package api

import (
	"log"
	"net/http"

	"github.com/google/uuid"

	"hearts/internal/db"
)

type CreateMatchResponse struct {
	MatchId uuid.UUID `json:"match_id"`
}

func CreateMatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dbConn, err := GetDatabaseOrFail(r.Context(), w)
	if err != nil || dbConn == nil {
		return
	}

	id, err := db.CreateMatch(dbConn)
	if err != nil {
		log.Println("ERROR: failed to create match:", err)
		http.Error(w, "failed to create match", http.StatusInternalServerError)
		return
	}
	log.Println("INFO: created match with id", id)

	WriteResponse(w, CreateMatchResponse{MatchId: id})
}
