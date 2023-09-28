package game

import (
	"errors"
	"fmt"
	"hearts/cards"
)

type Player struct {
	Hand *cards.Hand
}

type Table struct {
	Players            []Player
	CurrentPlayersTurn int
	CardsPlayed        []cards.Card
}

func (table *Table) CurrentPlayer() Player {
	return table.Players[table.CurrentPlayersTurn]
}

func (table *Table) IsRoundComplete() bool {
	return len(table.CurrentPlayer().Hand.Cards) == 0
}

func (table *Table) ValidCardsToPlay(hand *cards.Hand) []cards.Card {
	// the first play in a round must be the two of clubs
	if len(table.CardsPlayed) < 1 {
		twoOfClubs := cards.Card{Rank: "2", Suit: "♣"}
		if hand.Contains(twoOfClubs) {
			return []cards.Card{twoOfClubs}
		}
		return []cards.Card{}
	}

	// otherwise follow suit if possible
	matchesSuit := hand.FindCardsWithSuit(table.CardsPlayed[0].Suit)
	if len(matchesSuit) > 0 {
		return matchesSuit
	}

	// otherwise everything is valid
	return hand.Cards
}

func (table *Table) PlayCard(card cards.Card) error {
	currentPlayerHand := table.CurrentPlayer().Hand

	if !currentPlayerHand.Contains(card) {
		return fmt.Errorf("player %d cannot play card %s because it is not in their hand",
			table.CurrentPlayersTurn,
			card)
	}

	validPlays := table.ValidCardsToPlay(currentPlayerHand)
	if !cards.Contains(validPlays, card) {
		return fmt.Errorf("player %d cannot play card %s because it is not a valid play",
			table.CurrentPlayersTurn,
			card)
	}

	table.CardsPlayed = append(table.CardsPlayed, card)
	currentPlayerHand.Remove(card)
	table.CurrentPlayersTurn = (table.CurrentPlayersTurn + 1) % len(table.Players)
	return nil
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
