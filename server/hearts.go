package main

import (
	"encoding/json"
	"fmt"
	"hearts/cards"
	"hearts/game"
	"log"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func gameStateHandler(w http.ResponseWriter, r *http.Request) {
	table := game.Table{}
	table.AddSeats(4)

	deck := cards.CreateDeck()
	deck.Shuffle()
	table.Deal(deck)

	// TODO: we're serving the entire game state but we need to filter it for the requesting player's view
	tableJson, err := json.Marshal(table)
	if err != nil {
		http.Error(w, "Failed to marshal game state", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(tableJson))
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/game-state", gameStateHandler)

	fmt.Println("listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
