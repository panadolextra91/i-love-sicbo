package engine

type DiceEvaluator struct {
	Dice DiceResult
}

func NewDiceEvaluator(d DiceResult) DiceEvaluator {
	return DiceEvaluator{Dice: d}
}

func (e DiceEvaluator) Sum() int {
	return e.Dice[0] + e.Dice[1] + e.Dice[2]
}

func (e DiceEvaluator) IsTriple() bool {
	return e.Dice[0] == e.Dice[1] && e.Dice[1] == e.Dice[2]
}

func (e DiceEvaluator) CountOccurrences(number int) int {
	count := 0
	for _, v := range e.Dice {
		if v == number {
			count++
		}
	}
	return count
}

func (e DiceEvaluator) HasBoth(a, b int) bool {
	return e.CountOccurrences(a) > 0 && e.CountOccurrences(b) > 0
}

func DecodeTwoNumberCombo(v int) (int, int) {
	a := v / 10
	b := v % 10
	return a, b
}
