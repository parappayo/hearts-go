package cards

type Hand struct {
	Cards []Card
}

func (hand *Hand) Contains(card Card) bool {
	return Contains(hand.Cards, card)
}
