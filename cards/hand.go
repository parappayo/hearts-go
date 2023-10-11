package cards

type Hand struct {
	Cards []Card
}

func (hand *Hand) String() string {
	return CardsToString(hand.Cards)
}

func (hand *Hand) Contains(card Card) bool {
	return Contains(hand.Cards, card)
}

func (hand *Hand) FindCards(predicate func(Card) bool) []Card {
	result := make([]Card, 0)
	for i := range hand.Cards {
		card := hand.Cards[i]
		if predicate(card) {
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
