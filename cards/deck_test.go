//
//  Run tests:
//    go test cards/deck_test.go
//

package cards_test

import (
	"hearts/cards"
	"testing"
)

func TestCreateCustomDeck(t *testing.T) {
	suits := []string{"one", "two", "three"}
	ranks := []string{"a", "b"}
	deck := cards.CreateCustomDeck(suits, ranks)

	{
		expected := len(suits) * len(ranks)
		actual := len(deck.Cards)
		if expected != actual {
			t.Fatalf("TestCreateCustomDeck wrong deck.CardsPlayed length, expected %d, actual %d", expected, actual)
		}
	}

	{
		expected := cards.Card{Rank: "a", Suit: "one"}
		actual := deck.Cards[0]
		if expected != actual {
			t.Fatalf("TestCreateCustomDeck wrong first card, expected %s, actual %s", expected, actual)
		}
	}
}

func TestCreateDeck(t *testing.T) {
	deck := cards.CreateDeck()

	{
		expected := 52
		actual := len(deck.Cards)
		if expected != actual {
			t.Fatalf("TestCreateDeck wrong deck.CardsPlayed length, expected %d, actual %d", expected, actual)
		}
	}

	{
		expected := cards.Card{Rank: cards.Ranks[0], Suit: cards.Suits[0]}
		actual := deck.Cards[0]
		if expected != actual {
			t.Fatalf("TestCreateDeck wrong first card, expected %s, actual %s", expected, actual)
		}
	}
}

func TestShuffleDoesNotRemoveCards(t *testing.T) {
	deck := cards.CreateDeck()
	oldFirstCard := deck.Cards[0]

	oldLen := len(deck.Cards)
	deck.Shuffle()
	newLen := len(deck.Cards)

	if newLen != oldLen {
		t.Fatalf("TestShuffle expected %d, actual %d", oldLen-1, newLen)
	}

	if !deck.Contains(oldFirstCard) {
		t.Fatalf("TestShuffle failed to find old first card")
	}
}

func TestDrawRemovesOneCard(t *testing.T) {
	deck := cards.CreateDeck()

	oldLen := len(deck.Cards)
	card := deck.Draw()
	newLen := len(deck.Cards)

	if newLen != oldLen-1 {
		t.Fatalf("TestDraw expected %d, actual %d", oldLen-1, newLen)
	}

	if deck.Contains(card) {
		t.Fatalf("TestDraw expected card to be removed from deck after being drawn")
	}
}

func TestShuffleAndDrawRemovesOneCard(t *testing.T) {
	deck := cards.CreateDeck()
	deck.Shuffle()

	oldLen := len(deck.Cards)
	card := deck.Draw()
	newLen := len(deck.Cards)

	if newLen != oldLen-1 {
		t.Fatalf("TestDraw expected %d, actual %d", oldLen-1, newLen)
	}

	if deck.Contains(card) {
		t.Fatalf("TestDraw expected card to be removed from deck after being drawn")
	}
}
