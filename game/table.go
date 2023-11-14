package game

import (
	"errors"
	"fmt"
	"hearts/cards"
)

type Player struct {
	Hand  *cards.Hand
	Score int
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

func (table *Table) CurrentTrickCardsPlayedCount() int {
	return len(table.CardsPlayed) % len(table.Players)
}

func (table *Table) IsTrickFinished() bool {
	return table.CurrentTrickCardsPlayedCount() == 0 &&
		len(table.CardsPlayed) > 0
}

func (table* Table) CurrentTrickStartPlayerIndex() int {
	result := table.CurrentPlayersTurn - table.CurrentTrickCardsPlayedCount()
	if result < 0 {
		result += len(table.Players)
	}
	return result
}

func (table *Table) CurrentTrick() Trick {
	result := Trick{
		StartPlayer: table.CurrentTrickStartPlayerIndex(),
	}
	trickCardsPlayedCount := table.CurrentTrickCardsPlayedCount()
	if trickCardsPlayedCount == 0 {
		result.CardsPlayed = []cards.Card{}
	} else {
		result.CardsPlayed = table.CardsPlayed[len(table.CardsPlayed)-trickCardsPlayedCount : len(table.CardsPlayed)]
	}
	return result
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

	trick := table.CurrentTrick()
	if len(trick.CardsPlayed) < 1 {
		// if hearts are broken, a new trick lead can be anything
		if table.AreHeartsBroken() {
			return hand.Cards
		}

		// otherwise, no hearts unless that's the player's entire hand
		nonHeartsCards := hand.FindCards(
			func(card cards.Card) bool {
				return card.Suit != "♥"
			})
		if len(nonHeartsCards) < 1 {
			return hand.Cards
		}
		return nonHeartsCards
	}

	// follow suit if possible
	matchesSuit := hand.FindCards(
		func(card cards.Card) bool {
			return card.Suit == trick.CardsPlayed[0].Suit
		})
	if len(matchesSuit) > 0 {
		return matchesSuit
	}

	// otherwise everything is valid
	return hand.Cards
}

func (table *Table) PlayCard(card cards.Card) (*Trick, error) {
	currentPlayerHand := table.CurrentPlayer().Hand

	if !currentPlayerHand.Contains(card) {
		return nil, fmt.Errorf(
			"player %d cannot play card %s because it is not in their hand %s",
			table.CurrentPlayersTurn,
			card,
			currentPlayerHand)
	}

	validPlays := table.ValidCardsToPlay(currentPlayerHand)
	if !cards.Contains(validPlays, card) {
		return nil, fmt.Errorf(
			"player %d cannot play card %s because it is not a valid play",
			table.CurrentPlayersTurn,
			card)
	}

	trick := table.CurrentTrick()
	table.CardsPlayed = append(table.CardsPlayed, card)
	currentPlayerHand.Remove(card)
	trick.CardsPlayed = append(trick.CardsPlayed, card)

	if len(trick.CardsPlayed) == len(table.Players) {
		winnerIndex := trick.Winner()
		winner := &(table.Players[winnerIndex])
		winner.Score += trick.Score()
		table.CurrentPlayersTurn = winnerIndex
	} else {
		table.CurrentPlayersTurn = (table.CurrentPlayersTurn + 1) % len(table.Players)
	}

	return &trick, nil
}

func (table *Table) AddSeats(seatCount int) {
	table.Players = make([]Player, seatCount)
}

func (table *Table) Deal(deck cards.Deck) error {
	var err error
	table.CardsPlayed = nil
	hands := deck.Deal(len(table.Players))
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

func (table *Table) AreHeartsBroken() bool {
	for i := range table.CardsPlayed {
		if table.CardsPlayed[i].Suit == "♥" {
			return true
		}
	}
	return false
}
