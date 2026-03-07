package strategies

import "cachon-casino/internal/engine"

type singleNumberStrategy struct{}

func NewSingleNumberStrategy() engine.BetStrategy { return singleNumberStrategy{} }

func (singleNumberStrategy) Type() engine.BetType { return engine.BetSingleNumber }

func (singleNumberStrategy) Validate(b engine.Bet) error { return validateTargetRange(b, 1, 6) }

func (singleNumberStrategy) IsWin(d engine.DiceResult, b engine.Bet) bool {
	return engine.NewDiceEvaluator(d).CountOccurrences(b.TargetValue) > 0
}

func (singleNumberStrategy) Odds(b engine.Bet, d engine.DiceResult) int64 {
	return int64(engine.NewDiceEvaluator(d).CountOccurrences(b.TargetValue))
}
