package cards

import (
	"math/rand"
)

type Deck struct {
	Cards []Card
}

func CreateCustomDeck(suits []string, ranks []string) Deck {
	deck := Deck{}
	for _, suit := range suits {
		for _, rank := range ranks {
			deck.Cards = append(deck.Cards, Card{Suit: suit, Rank: rank})
		}
	}
	return deck
}

func CreateDeck() Deck {
	return CreateCustomDeck(Suits, Ranks)
}

func (deck *Deck) Contains(card Card) bool {
	return Contains(deck.Cards, card)
}

func (deck *Deck) Shuffle() {
	for i := range deck.Cards {
		j := rand.Intn(i + 1)
		deck.Cards[i], deck.Cards[j] = deck.Cards[j], deck.Cards[i]
	}
}

func (deck *Deck) Draw() Card {
	result := deck.Cards[0]
	deck.Cards = deck.Cards[1:]
	return result
}

func (deck *Deck) Deal(handCount uint8) []Hand {
	hands := make([]Hand, handCount)
	i := uint8(0)
	for len(deck.Cards) > 0 {
		hands[i].Cards = append(hands[i].Cards, deck.Draw())
		i = (i + 1) % handCount
	}
	return hands
}
