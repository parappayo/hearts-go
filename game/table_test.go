package game_test

import (
	"hearts/cards"
	"hearts/game"
	"testing"
)

func TestFullGame(t *testing.T) {
	maxTrickScore := 26
	table := game.Table{}
	table.AddSeats(4)

	deck := cards.CreateDeck()
	deck.Shuffle()
	table.Deal(deck)

	for !table.IsRoundComplete() {
		validPlays := table.ValidCardsToPlay(table.CurrentPlayer().Hand)
		if len(validPlays) < 1 {
			t.Fatalf("round is incomplete but player has no valid plays")
		}

		currentTrick := table.CurrentTrick()
		if len(currentTrick.CardsPlayed) >= len(table.Players) {
			t.Fatalf("trick has more cards than there are players")
		}

		trick, err := table.PlayCard(validPlays[0])
		if err != nil {
			t.Fatalf(err.Error())
			return
		}

		if len(trick.CardsPlayed) == int(len(table.Players)) {
			if trick.Score() < 0 {
				t.Fatalf("trick score is negative")
			}
			if trick.Score() > maxTrickScore {
				t.Fatalf("trick score greater than max allowed")
			}
			winner := trick.Winner()
			if winner < 0 {
				t.Fatalf("trick winner is negative")
			}
			if winner > len(table.Players)-1 {
				t.Fatalf("tricker winner index is too high")
			}
		}
	}

	totalScore := 0
	for i := range table.Players {
		totalScore += table.Players[i].Score
		if table.Players[i].Score < 0 {
			t.Fatalf("player has negative score")
		}
	}
	if totalScore != maxTrickScore {
		t.Fatalf("final scores do not sum to the max trick score")
	}
}

func TestDeal(t *testing.T) {
	table := game.Table{}
	table.AddSeats(4)
	table.Deal(cards.CreateDeck())

	if len(table.Players) != 4 {
		t.Fatalf("player count does not match specified seat count")
	}
	for _, player := range table.Players {
		if player.Score != 0 {
			t.Fatalf("player is starting with non-zero score")
		}
	}

	countPlayersWithTwoOfClubs := 0
	playerWithTwoOfClubs := -1
	for i, player := range table.Players {
		expected := 52 / 4
		actual := len(player.Hand.Cards)
		if expected != actual {
			t.Fatalf("player starting hand size is wrong: expected %d, actual %d", expected, actual)
		}

		if player.Hand.Contains(cards.Card{Rank: "2", Suit: "♣"}) {
			countPlayersWithTwoOfClubs += 1
			playerWithTwoOfClubs = i
		}
	}
	if countPlayersWithTwoOfClubs != 1 {
		t.Fatalf("unexpected number of players who have two of clubs: %d", countPlayersWithTwoOfClubs)
	}
	if playerWithTwoOfClubs != table.CurrentPlayersTurn {
		t.Fatalf("player holding two of clubs is not the first player to play")
	}
	if table.AreHeartsBroken() {
		t.Fatalf("hearts are broken after new deal")
	}
}

func TestCurrentTrick(t *testing.T) {
	table := game.Table{}
	table.AddSeats(4)
	table.Deal(cards.CreateDeck())

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

func testPlayCardHelper(t *testing.T, table *game.Table, card cards.Card) *game.Trick {
	oldPlayer := table.CurrentPlayersTurn
	trick, err := table.PlayCard(card)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !table.IsTrickFinished() && table.CurrentPlayersTurn != (oldPlayer+1)%len(table.Players) {
		t.Fatalf("played a card, but it didn't increment the player turn, %d -> %d",
			oldPlayer, table.CurrentPlayersTurn)
	}
	return trick
}

func TestValidCardsToPlay(t *testing.T) {
	table := game.Table{}
	table.AddSeats(4)
	table.Deal(cards.CreateDeck())

	// the only valid card on first lead is the two of clubs
	validPlays := table.ValidCardsToPlay(table.CurrentPlayer().Hand)
	if len(validPlays) != 1 {
		t.Fatalf("expected only one valid first play (two of clubs)")
	}
	twoOfClubs := cards.Card{Rank: "2", Suit: "♣"}
	if validPlays[0] != twoOfClubs {
		t.Fatalf("valid first play is not two of clubs, instead %s", validPlays[0])
	}
	if table.CurrentPlayersTurn != 1 {
		t.Fatalf("expected player 1 to start the round but player %d did", table.CurrentPlayersTurn)
	}

	_, err := table.PlayCard(twoOfClubs)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if table.CurrentPlayersTurn != 2 {
		t.Fatalf("played a card, but it didn't increment the player turn, %d", table.CurrentPlayersTurn)
	}

	// player must follow lead suit when possible
	validPlays = table.ValidCardsToPlay(table.CurrentPlayer().Hand)
	expectedCards := []cards.Card{
		{Rank: "3", Suit: "♣"},
		{Rank: "7", Suit: "♣"},
		{Rank: "J", Suit: "♣"},
	}
	if !cards.Equals(validPlays, expectedCards) {
		t.Fatalf("valid plays %s does not match expected %s", cards.ToString(validPlays), cards.ToString(expectedCards))
	}

	testPlayCardHelper(t, &table, cards.Card{Rank: "7", Suit: "♣"})          // player 2
	testPlayCardHelper(t, &table, cards.Card{Rank: "8", Suit: "♣"})          // player 3
	trick := testPlayCardHelper(t, &table, cards.Card{Rank: "9", Suit: "♣"}) // player 0

	if !table.IsTrickFinished() {
		t.Fatal("expected trick to be finished but it is not")
	}
	if trick.Winner() != 0 {
		t.Fatalf("expected player 0 to have won the trick, but player %d did", trick.Winner())
	}
	if table.CurrentPlayersTurn != 0 {
		t.Fatalf("expected player 0 to lead the next trick, but player %d did", table.CurrentPlayersTurn)
	}

	// player cannot lead hearts when hearts are not broken
	validPlays = table.ValidCardsToPlay(table.CurrentPlayer().Hand)
	for _, card := range validPlays {
		if card.Suit == "♥" {
			t.Fatalf("hearts are not broken but player can play hearts")
		}
	}

	testPlayCardHelper(t, &table, cards.Card{Rank: "K", Suit: "♣"})  // player 0
	testPlayCardHelper(t, &table, cards.Card{Rank: "10", Suit: "♣"}) // player 1
	testPlayCardHelper(t, &table, cards.Card{Rank: "J", Suit: "♣"})  // player 2
	testPlayCardHelper(t, &table, cards.Card{Rank: "Q", Suit: "♣"})  // player 3

	testPlayCardHelper(t, &table, cards.Card{Rank: "5", Suit: "♣"}) // player 0
	testPlayCardHelper(t, &table, cards.Card{Rank: "A", Suit: "♣"}) // player 1
	testPlayCardHelper(t, &table, cards.Card{Rank: "3", Suit: "♣"}) // player 2
	testPlayCardHelper(t, &table, cards.Card{Rank: "4", Suit: "♣"}) // player 3

	// now player 1 will lead clubs, and player 2 has a void in clubs
	testPlayCardHelper(t, &table, cards.Card{Rank: "6", Suit: "♣"}) // player 1

	// player can play anything if they have a void in the lead suit
	if table.CurrentPlayersTurn != 2 {
		t.Fatalf("something went wrong, expected player turn 2 but it is %d", table.CurrentPlayersTurn)
	}
	validCards := table.ValidCardsToPlay(table.CurrentPlayer().Hand)
	if len(validCards) != len(table.CurrentPlayer().Hand.Cards) {
		t.Fatalf("expected player to be able to play anything in their hand")
	}

	// player 2 will break hearts
	testPlayCardHelper(t, &table, cards.Card{Rank: "2", Suit: "♥"})
	// player 3
	// player 4

	// figuring out what the current hand is
	//t.Fatalf("test %s", table.ValidCardsToPlay(table.CurrentPlayer().Hand))
	// player can lead hearts if hearts are broken
	// player can lead hearts if hearts are not broken but that's all the player has
}

func TestPlayCard(t *testing.T) {
	table := game.Table{}
	table.AddSeats(4)
	table.Deal(cards.CreateDeck())

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
