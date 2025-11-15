package db

import (
	"database/sql"
	"encoding/json"
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

type MatchEvent struct {
	Type string
	AggregateVersion uint32
	CreatedOn string
	Payload json.RawMessage
}

func QueryMatchEvents(db *sql.DB, matchId uuid.UUID) ([]MatchEvent, error) {
	rows, err := db.Query(`
SELECT type, aggregate_version, created_on, payload
FROM hearts.event
WHERE aggregate_id = $1
ORDER BY aggregate_version
`, matchId)
	if err != nil {
		return nil, err
	}
	result := make([]MatchEvent, 0, 5)
	for rows.Next() {
		var row MatchEvent
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
	timestamp := time.Now().UTC().Format(time.RFC3339)
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

type JoinMatchPayload struct {
	UserId uuid.UUID `json:"user_id"`
}

func JoinMatch(db *sql.DB, userId uuid.UUID, matchId uuid.UUID, version uint32) (*MatchEvent, error) {
	payload, err := json.Marshal(JoinMatchPayload{UserId: userId})
	if err != nil {
		return nil, err
	}
	timestamp := time.Now().UTC().Format(time.RFC3339)
	_, err = db.Exec(`
INSERT INTO hearts.event (
	type,
	version,
	created_on,
	aggregate_id,
	aggregate_version,
	payload
)
VALUES (
	'player-joined',
	'1.0.0',
	$1,
	$2,
	$3,
	$4
)
`, timestamp, matchId, version, payload)
	return &MatchEvent{
		Type: "player-joined",
		AggregateVersion: version,
		CreatedOn: timestamp,
		Payload: payload,
	}, err
}
