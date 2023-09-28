package game_test

import (
	"hearts/cards"
	"hearts/game"
	"testing"
)

func TestPlayCard(t *testing.T) {
	table := game.Table{}
	table.Deal(cards.CreateDeck())

	player := table.CurrentPlayer()
	oldHandLen := len(player.Hand.Cards)
	expectedHandLen := 52 / 4
	if oldHandLen != expectedHandLen {
		t.Fatalf("expected dealt hand len %d, found dealt hand len %d", expectedHandLen, oldHandLen)
	}

	validPlays := table.ValidCardsToPlay(player.Hand)
	err := table.PlayCard(validPlays[0])

	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(table.CardsPlayed) != 1 {
		t.Fatalf("expected 1 card played, found %d cards played", len(table.CardsPlayed))
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
