package cards

var Suits = []string{"♠", "♣", "♥", "♦"}
var Ranks = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

type Card struct {
	Rank string
	Suit string
}

func Contains(cards []Card, card Card) bool {
	return FindIndex(cards, card) != -1
}

func FindIndex(cards []Card, card Card) int {
	for i := range cards {
		if cards[i] == card {
			return i
		}
	}
	return -1
}

func Remove(cards []Card, i int) []Card {
	// fast slice removal: overwrite the removed element and drop the last element
	cards[i] = cards[len(cards)-1]
	return cards[:len(cards)-1]
}
