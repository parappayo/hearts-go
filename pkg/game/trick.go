package game

import (
	"fmt"
	"hearts/cards"
)

type Trick struct {
	StartPlayer int
	CardsPlayed []cards.Card
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func (trick *Trick) String() string {
	return fmt.Sprintf("StartPlayer %d CardsPlayed %s", trick.StartPlayer, trick.CardsPlayed)
}

func (trick *Trick) Winner() int {
	if len(trick.CardsPlayed) < 1 {
		return 0
	}

	leadSuit := trick.CardsPlayed[0].Suit
	highestRank := trick.CardsPlayed[0].RankValue(cards.Ranks)
	winnerIndex := trick.StartPlayer
	numPlayers := max(trick.StartPlayer, len(trick.CardsPlayed))

	for i := 1; i < len(trick.CardsPlayed); i++ {
		card := trick.CardsPlayed[i]
		value := card.RankValue(cards.Ranks)
		if card.Suit == leadSuit && value > highestRank {
			highestRank = value
			winnerIndex = (trick.StartPlayer + i) % numPlayers
		}
	}

	return winnerIndex
}

func (trick *Trick) Score() int {
	result := 0
	for i := range trick.CardsPlayed {
		card := trick.CardsPlayed[i]
		if card.Suit == "♥" {
			result += 1
		} else if card.Suit == "♠" && card.Rank == "Q" {
			result += 13
		}
	}
	return result
}
