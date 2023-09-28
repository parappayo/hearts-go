package cards

type Hand struct {
	Cards []Card
}

func (hand *Hand) Contains(card Card) bool {
	return Contains(hand.Cards, card)
}

func (hand *Hand) FindCardsWithSuit(suit string) []Card {
	result := make([]Card, 0)
	for i := range hand.Cards {
		card := hand.Cards[i]
		if card.Suit == suit {
			result = append(result, card)
		}
	}
	return result
}

func (hand *Hand) Remove(card Card) {
	i := FindIndex(hand.Cards, card)
	if i < 0 {
		return
	}
	hand.Cards = Remove(hand.Cards, i)
}
