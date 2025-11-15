package db

import (
	"database/sql"
	"os"

	"hearts/internal/agg"

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

func QueryMatchEvents(db *sql.DB, matchId uuid.UUID) ([]agg.MatchEvent, error) {
	rows, err := db.Query(`
SELECT type, aggregate_version, created_on, payload
FROM hearts.event
WHERE aggregate_id = $1
ORDER BY aggregate_version
`, matchId)
	if err != nil {
		return nil, err
	}
	result := make([]agg.MatchEvent, 0, 5)
	for rows.Next() {
		var row agg.MatchEvent
		err := rows.Scan(&row.Type, &row.AggregateVersion, &row.CreatedOn, &row.Payload)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func CreateMatch(db *sql.DB) (uuid.UUID, error) {
	timestamp := agg.Timestamp()
	matchId := uuid.New()
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
`, timestamp, matchId).Scan(&matchId)
	return matchId, err
}

func InsertEvent(db *sql.DB, aggId uuid.UUID, event *agg.MatchEvent) error {
	_, err := db.Exec(`
INSERT INTO hearts.event (
	type,
	version,
	created_on,
	aggregate_id,
	aggregate_version,
	payload
)
VALUES (
	$1,
	'1.0.0',
	$2,
	$3,
	$4,
	$5
)
`,
		event.Type,
		event.CreatedOn,
		aggId,
		event.AggregateVersion,
		event.Payload)
	return err
}
