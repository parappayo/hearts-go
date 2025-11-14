package db

import (
	"errors"
	"fmt"
	"log"
)

type Player struct {
	Name string
}

type MatchState struct {
	Version uint32
	CreatedOn string
	Players []Player
}

func GetAggregate(events[] MatchEvent) (*MatchState, error) {
	result := MatchState {
		Version: 0,
		CreatedOn: "",
		Players: make([]Player, 0, 4),
	}
	for i, event := range events {
		if event.AggregateVersion != uint32(i+1) {
			return nil, errors.New("events not given in order")
		}
		result.Version = event.AggregateVersion

		switch event.Type {

		case "match-created":
			result.CreatedOn = event.CreatedOn

		case "player-joined":
			playerCount := len(result.Players)
			result.Players = append(result.Players, Player{Name: fmt.Sprintf("Seat %d", playerCount+1)})

		default:
			log.Println("WARN: unrecognized event type", event.Type)
		}
	}
	return &result, nil
}
