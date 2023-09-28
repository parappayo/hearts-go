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

func main() {
	rand.Seed(int64(time.Now().Nanosecond()))
	table := game.Table{}
	table.Deal(cards.CreateDeck())
	for i := range table.Players {
		fmt.Println(table.Players[i].Hand)
	}

	for !table.IsRoundComplete() {
		validPlays := table.ValidCardsToPlay(table.CurrentPlayer().Hand)
		fmt.Printf("player %d can play %d\n", table.CurrentPlayersTurn, validPlays)

		err := table.PlayCard(validPlays[0])
		if err != nil {
			fmt.Println("failed to play card: ", err)
			return
		}
	}
}
