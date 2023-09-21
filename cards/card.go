package cards

var Suits = []string{"♠", "♣", "♥", "♦"}
var Ranks = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

type Card struct {
	Rank string
	Suit string
}

func Contains(cards []Card, card Card) bool {
	for i := range cards {
		if cards[i] == card {
			return true
		}
	}
	return false
}
