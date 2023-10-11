//
//  Run:
//    go run hearts.go
//

package main

import (
	"fmt"
	"hearts/cards"
	"hearts/game"
	"math/rand"
	"time"
)

func printHands(table *game.Table) {
	for i := range table.Players {
		fmt.Printf("seat %d hand %s\n", i, table.Players[i].Hand)
	}
}

func printScores(table *game.Table) {
	for i := range table.Players {
		fmt.Printf("seat %d: score %d\n", i, table.Players[i].Score)
	}
}

func main() {
	seatCount := uint8(4)
	rand.Seed(int64(time.Now().Nanosecond()))
	table := game.Table{}
	table.Deal(cards.CreateDeck(), seatCount)
	printHands(&table)

	for !table.IsRoundComplete() {
		validPlays := table.ValidCardsToPlay(table.CurrentPlayer().Hand)
		currentTrick := table.CurrentTrick()
		fmt.Printf("trick %s\n", currentTrick.String())
		fmt.Printf("player %d can play %s\n", table.CurrentPlayersTurn, cards.CardsToString(validPlays))

		trick, err := table.PlayCard(validPlays[0])
		if err != nil {
			fmt.Println("failed to play card: ", err)
			return
		}

		if len(trick.CardsPlayed) == int(seatCount) {
			fmt.Printf("trick worth %d points taken by seat %d\n", trick.Score(), trick.Winner())
		}
	}

	fmt.Println("Round Scores")
	printScores(&table)
}
