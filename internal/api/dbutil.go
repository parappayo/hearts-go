package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
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
