package agg

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"hearts/pkg/cards"
)

type MatchEvent struct {
	Type string
	AggregateVersion uint32
	CreatedOn string
	Payload json.RawMessage
}

// TODO: should mix-in the game.Player type?
type Player struct {
	ID uuid.UUID
	Name string
}

// TODO: need a func that creates a game.Table from one of these
type MatchState struct {
	Version uint32
	CreatedOn string
	StartedOn string
	Players []Player
}

func (state *MatchState) ContainsPlayer(playerId uuid.UUID) bool {
	for i := range state.Players {
		if state.Players[i].ID == playerId {
			return true
		}
	}
	return false
}

func (state *MatchState) AddPlayer(playerId uuid.UUID, playerName string) error {
	if len(state.Players) > 3 {
		return errors.New("player cannot join, table full")
	}
	if state.ContainsPlayer(playerId) {
		return errors.New("player cannot join, already joined")
	}
	state.Players = append(
		state.Players,
		Player{
			ID: playerId,
			Name: playerName,
		})
	return nil
}

type JoinMatchPayload struct {
	// TODO: rename to player id
	UserId uuid.UUID `json:"user_id"`
}

type StartMatchPayload struct {
	CurrentPlayersTurn int
	Hands []cards.Hand
}

func Timestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func (state *MatchState) ApplyEvent(event *MatchEvent) error {
	if event.AggregateVersion != state.Version + 1 {
		return errors.New("cannot apply event with unexpected revision num")
	}
	if event.AggregateVersion == 1 && event.Type != "match-created" {
		return errors.New("first event must be to create a match")
	}

	switch event.Type {

	case "match-created":
		if state.CreatedOn != "" {
			return errors.New("match already created")
		}
		state.CreatedOn = event.CreatedOn

	case "player-joined":
		var joinMatchPayload JoinMatchPayload
		err := json.Unmarshal(event.Payload, &joinMatchPayload)
		if err != nil {
			return err
		}
		playerName := fmt.Sprintf("Seat %d", len(state.Players)+1)
		err = state.AddPlayer(joinMatchPayload.UserId, playerName)
		if err != nil {
			return err
		}

	case "player-left":
		// TODO: error if player is not in the game
		return errors.New("event not implemented")

	case "match-started":
		// TODO: error if player count is not 4
		log.Println(string(event.Payload))
		return errors.New("event not implemented")

	case "card-played":
		return errors.New("event not implemented")

	case "round-finished":
		return errors.New("event not implemented")

	case "match-finished":
		return errors.New("event not implemented")

	default:
		log.Println("WARN: unrecognized event type", event.Type)
	}

	state.Version = event.AggregateVersion
	return nil
}

func GetAggregate(events[] MatchEvent) (*MatchState, error) {
	if len(events) < 1 {
		return nil, nil
	}
	result := MatchState {
		Version: 0,
		CreatedOn: "",
		Players: make([]Player, 0, 4),
	}
	for _, event := range events {
		err := result.ApplyEvent(&event)
		if err != nil {
			return nil, err
		}
	}
	return &result, nil
}
