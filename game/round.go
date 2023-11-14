package game

type Round struct {
	Tricks []Trick
}

func (round *Round) Scores() []int {
	if len(round.Tricks) < 1 {
		return nil
	}
	playerCount := len(round.Tricks[0].CardsPlayed)
	result := make([]int, playerCount)

	for _, trick := range round.Tricks {
		result[trick.Winner()] += trick.Score()
	}

	// TODO: check for "shoot the moon"
	return result
}
