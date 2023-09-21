package game

import (
	"errors"
	"hearts/cards"
)

type Player struct {
	Hand *cards.Hand
}

type Table struct {
	Players            []Player
	CurrentPlayersTurn int
}

func (table *Table) Deal(deck cards.Deck) error {
	var err error
	const seatCount = 4
	if len(table.Players) < seatCount {
		table.Players = make([]Player, seatCount)
	}
	deck.Shuffle()
	hands := deck.Deal(seatCount)
	for i := range table.Players {
		table.Players[i].Hand = &hands[i]
	}
	table.CurrentPlayersTurn, err = table.PlayerWhoGoesFirst()
	return err
}

func (table *Table) PlayerWhoHasCard(card cards.Card) (int, bool) {
	for i := range table.Players {
		if table.Players[i].Hand.Contains(card) {
			return i, true
		}
	}
	return 0, false
}

func (table *Table) PlayerWhoGoesFirst() (int, error) {
	playerIndex, playerFound := table.PlayerWhoHasCard(cards.Card{Rank: "2", Suit: "♣"})
	if playerFound {
		return playerIndex, nil
	}
	return 0, errors.New("player hands are not in a valid game start state")
}
