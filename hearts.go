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
	fmt.Println(table.CurrentPlayersTurn)
}
