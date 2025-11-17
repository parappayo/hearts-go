package api

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"

	"hearts/internal/agg"
	"hearts/internal/db"
)

// TODO: instead of passing around an sql connection, use an abstract interface
func WithDatabase(dbConn *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "db_conn", dbConn)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func GetDatabaseOrFail(ctx context.Context, w http.ResponseWriter) (*sql.DB, error) {
	dbConn, ok := ctx.Value("db_conn").(*sql.DB)
	if !ok {
		http.Error(w, "database not configured", http.StatusInternalServerError)
		return nil, errors.New("failed to get db_conn from context")
	}
	return dbConn, nil
}

func QueryAggregate(dbConn *sql.DB, matchId uuid.UUID, w http.ResponseWriter) (*agg.MatchState, error) {
	events, err := db.QueryMatchEvents(dbConn, matchId)
	if err != nil {
		log.Println("ERROR: failed to query match events:", err)
		http.Error(w, "failed to query match", http.StatusInternalServerError)
		return nil, err
	}
	agg, err := agg.GetAggregate(events)
	if err != nil {
		log.Println("ERROR: failed to query match events:", err)
		http.Error(w, "failed to query match", http.StatusInternalServerError)
		return nil, err
	}
	if agg == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return nil, errors.New("aggregate not found")
	}
	return agg, nil
}
