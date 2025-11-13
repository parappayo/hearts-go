package db

import (
	"database/sql"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)


func Open(connectionInfo string) (*sql.DB, error) {
	return sql.Open("postgres", connectionInfo)
}

func CreateSchema(db *sql.DB) (error) {
	queryText, err := os.ReadFile("db/schema.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(queryText))
	return err
}

func CreateMatch(db *sql.DB) (uuid.UUID, error) {
	var matchId uuid.UUID
	timestamp := time.Now().UTC().Format(time.RFC3339)
	aggregateId := uuid.New()
	err := db.QueryRow(`
INSERT INTO hearts.event (
	type,
	version,
	created_on,
	aggregate_id,
	aggregate_version,
	payload
)
VALUES (
	'match-created',
	'1.0.0',
	$1,
	$2,
	1,
	'{}'
)
RETURNING aggregate_id
`, timestamp, aggregateId).Scan(&matchId)
	return matchId, err
}
