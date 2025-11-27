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
	// TODO: don't pass around json data, instead use any type here and serialize when ready to write
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
	Hands []*cards.Hand
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

type PlayCardPayload struct {
	Card cards.Card
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
		playerCount := len(state.Players)
		if playerCount != 4 {
			return errors.New("cannot start match unless there are four players")
		}
		var payload StartMatchPayload
		err := json.Unmarshal(event.Payload, &payload)
		if err != nil {
			return err
		}
		state.Hands = make([]*cards.Hand, 0, 4)
		for i := range payload.Hands {
			state.Hands = append(state.Hands, &payload.Hands[i])
		}

	case "card-played":
		var payload PlayCardPayload
		err := json.Unmarshal(event.Payload, &payload)
		if err != nil {
			return err
		}
		// TODO: update the Table here
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
