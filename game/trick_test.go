package game_test

import (
	"hearts/cards"
	"hearts/game"
	"testing"
)

func TestWinnerFirstPlayer(t *testing.T) {
	trick := game.Trick{}

	trick.CardsPlayed = []cards.Card{
		{Rank: "2", Suit: "♣"},
		{Rank: "3", Suit: "♠"},
		{Rank: "4", Suit: "♠"},
		{Rank: "Q", Suit: "♦"},
	}
	expected := 0
	actual := trick.Winner()
	if expected != actual {
		t.Fatalf("wrong winner, expected %d but actual is %d", expected, actual)
	}

	for i := 1; i < 4; i++ {
		trick.StartPlayer = i
		expected = i
		actual = trick.Winner()
		if expected != actual {
			t.Fatalf("wrong winner, expected %d but actual is %d", expected, actual)
		}
	}
}

func TestWinnerThirdPlayer(t *testing.T) {
	trick := game.Trick{}

	trick.CardsPlayed = []cards.Card{
		{Rank: "2", Suit: "♠"},
		{Rank: "3", Suit: "♠"},
		{Rank: "4", Suit: "♠"},
		{Rank: "Q", Suit: "♦"},
	}

	for i := 0; i < 4; i++ {
		trick.StartPlayer = i
		expected := (i + 2) % 4
		actual := trick.Winner()
		if expected != actual {
			t.Fatalf("wrong winner, expected %d but actual is %d", expected, actual)
		}
	}
}

func TestScoreOfZero(t *testing.T) {
	trick := game.Trick{}

	trick.CardsPlayed = []cards.Card{
		{Rank: "2", Suit: "♣"},
		{Rank: "3", Suit: "♠"},
		{Rank: "4", Suit: "♠"},
		{Rank: "Q", Suit: "♦"},
	}
	expected := 0
	actual := trick.Score()
	if expected != actual {
		t.Fatalf("wrong score, expected %d but actual is %d", expected, actual)
	}
}

func TestScoreOfOne(t *testing.T) {
	trick := game.Trick{}

	trick.CardsPlayed = []cards.Card{
		{Rank: "2", Suit: "♣"},
		{Rank: "3", Suit: "♥"},
		{Rank: "4", Suit: "♠"},
		{Rank: "Q", Suit: "♦"},
	}
	expected := 1
	actual := trick.Score()
	if expected != actual {
		t.Fatalf("wrong score, expected %d but actual is %d", expected, actual)
	}
}

func TestScoreOfThirteen(t *testing.T) {
	trick := game.Trick{}

	trick.CardsPlayed = []cards.Card{
		{Rank: "2", Suit: "♣"},
		{Rank: "Q", Suit: "♣"},
		{Rank: "Q", Suit: "♠"},
		{Rank: "K", Suit: "♦"},
	}
	expected := 13
	actual := trick.Score()
	if expected != actual {
		t.Fatalf("wrong score, expected %d but actual is %d", expected, actual)
	}
}

func TestScoreOfFourteen(t *testing.T) {
	trick := game.Trick{}

	trick.CardsPlayed = []cards.Card{
		{Rank: "2", Suit: "♣"},
		{Rank: "3", Suit: "♥"},
		{Rank: "Q", Suit: "♠"},
		{Rank: "Q", Suit: "♦"},
	}
	expected := 14
	actual := trick.Score()
	if expected != actual {
		t.Fatalf("wrong score, expected %d but actual is %d", expected, actual)
	}
}

func TestScoreOfSixteen(t *testing.T) {
	trick := game.Trick{}

	trick.CardsPlayed = []cards.Card{
		{Rank: "Q", Suit: "♠"},
		{Rank: "K", Suit: "♥"},
		{Rank: "Q", Suit: "♥"},
		{Rank: "A", Suit: "♥"},
	}
	expected := 16
	actual := trick.Score()
	if expected != actual {
		t.Fatalf("wrong score, expected %d but actual is %d", expected, actual)
	}
}
