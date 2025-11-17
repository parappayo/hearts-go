package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"

	"hearts/internal/agg"
	"hearts/internal/api"
	"hearts/internal/db"
)

func getDatabaseOrFail(ctx context.Context, w http.ResponseWriter) (*sql.DB, error) {
	dbConn, ok := ctx.Value("db_conn").(*sql.DB)
	if !ok {
		http.Error(w, "database not configured", http.StatusInternalServerError)
		return nil, errors.New("failed to get db_conn from context")
	}
	return dbConn, nil
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

	dbConn, err := getDatabaseOrFail(r.Context(), w)
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

	dbConn, err := getDatabaseOrFail(r.Context(), w)
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

	dbConn, err := getDatabaseOrFail(r.Context(), w)
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

	api.WriteResponse(w, agg)
}

// TODO: instead of passing around an sql connection, use an abstract interface
func WithDatabase(dbConn *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "db_conn", dbConn)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func main() {
	dbConn, err := db.Open(os.Getenv("DB_CONN"))
	if err != nil {
		panic(err)
	}
	defer dbConn.Close()
	err = db.CreateSchema(dbConn)
	if err != nil {
		panic(err)
	}

	// TODO: health endpoint should also test db connectivity
	http.Handle("/health",
		api.CommonHeaders(
			http.HandlerFunc(api.HealthHandler)))

	http.Handle("/match/{id}",
		WithDatabase(
			dbConn,
			api.CommonHeaders(
				http.HandlerFunc(matchStateHandler))))

	http.Handle("/create-match",
		WithDatabase(
			dbConn,
			api.CommonHeaders(
				http.HandlerFunc(createMatchHandler))))

	http.Handle(
		"/join-match",
		WithDatabase(
			dbConn,
			api.CommonHeaders(http.HandlerFunc(joinMatchHandler))))

	fmt.Println("listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
