package game_test

import (
	"hearts/cards"
	"hearts/game"
	"testing"
)

func TestCurrentTrick(t *testing.T) {
	table := game.Table{}
	table.Deal(cards.CreateDeck(), 4)

	table.CardsPlayed = []cards.Card{}
	trick := table.CurrentTrick()
	expected := len(table.CardsPlayed)
	actual := len(trick.CardsPlayed)
	if expected != actual {
		t.Fatalf("expected trick len %d but actual is %d", expected, actual)
	}

	playedTrick, err := table.PlayCard(cards.Card{Rank: "2", Suit: "♣"})
	if err != nil {
		t.Fatalf(err.Error())
	}

	trick = table.CurrentTrick()
	expected = len(table.CardsPlayed)
	actual = len(trick.CardsPlayed)
	if expected != actual {
		t.Fatalf("expected trick len %d but actual is %d", expected, actual)
	}

	if playedTrick.CardsPlayed[0] != trick.CardsPlayed[0] {
		t.Fatalf("trick returned from PlayCard does not match CurrentTrick")
	}
}

func TestPlayCard(t *testing.T) {
	table := game.Table{}
	table.Deal(cards.CreateDeck(), 4)

	player := table.CurrentPlayer()
	oldHandLen := len(player.Hand.Cards)
	expectedHandLen := 52 / 4
	if oldHandLen != expectedHandLen {
		t.Fatalf("expected dealt hand len %d, found dealt hand len %d", expectedHandLen, oldHandLen)
	}

	validPlays := table.ValidCardsToPlay(player.Hand)
	trick, err := table.PlayCard(validPlays[0])

	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(table.CardsPlayed) != 1 {
		t.Fatalf("expected 1 table card played, found %d cards played", len(table.CardsPlayed))
	}

	if len(trick.CardsPlayed) != 1 {
		t.Fatalf("expected 1 trick card played, found %d cards played", len(trick.CardsPlayed))
	}

	// player hand len should decrease by one
	newHandLen := len(player.Hand.Cards)
	if newHandLen != oldHandLen-1 {
		t.Fatalf("expected hand len after play to be %d, actual hand len %d", oldHandLen-1, newHandLen)
	}

	// current player should have changed
	if table.CurrentPlayer() == player {
		t.Fatalf("expected turn to pass to the next player, but current player is the same")
	}
}
