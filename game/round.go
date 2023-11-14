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

	// check for "shoot the moon"
	firstPlayerWithAnyPoints := -1
	numberOfPlayersWithAnyPoints := 0
	for i := range result {
		if result[i] > 0 {
			numberOfPlayersWithAnyPoints += 1
			if firstPlayerWithAnyPoints == -1 {
				firstPlayerWithAnyPoints = i
			}
		}
	}
	maxPoints := 26
	if numberOfPlayersWithAnyPoints == 1 {
		for i := range result {
			if i == numberOfPlayersWithAnyPoints {
				result[i] = 0
			} else {
				result[i] = maxPoints
			}
		}
	}

	return result
}
