package cards

import (
	"fmt"
	"strings"
)

var Suits = []string{"♠", "♣", "♥", "♦"}
var Ranks = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

type Card struct {
	Rank string
	Suit string
}

func (card Card) String() string {
	return fmt.Sprintf("[%s %s]", card.Suit, card.Rank)
}

func ToString(cards []Card) string {
	cardStrings := make([]string, len(cards))
	for i := range cards {
		cardStrings[i] = cards[i].String()
	}
	//slices.Sort(cardStrings)
	return strings.Join(cardStrings, " ")
}

func Equals(cardsA []Card, cardsB []Card) bool {
	if len(cardsA) != len(cardsB) {
		return false
	}
	for i := range cardsA {
		if cardsA[i] != cardsB[i] {
			return false
		}
	}
	return true
}

func (card Card) RankValue(ranks []string) int {
	for i := 0; i < len(Ranks); i++ {
		if ranks[i] == card.Rank {
			return i
		}
	}
	return -1
}

func Contains(cards []Card, card Card) bool {
	return FindIndex(cards, card) != -1
}

func FindIndex(cards []Card, card Card) int {
	for i := range cards {
		if cards[i] == card {
			return i
		}
	}
	return -1
}

func Remove(cards []Card, i int) []Card {
	// fast slice removal: overwrite the removed element and drop the last element
	cards[i] = cards[len(cards)-1]
	return cards[:len(cards)-1]
}
