package strategies

import "cachon-casino/internal/engine"

type exactDoubleStrategy struct{}

func NewExactDoubleStrategy() engine.BetStrategy { return exactDoubleStrategy{} }

func (exactDoubleStrategy) Type() engine.BetType { return engine.BetExactDouble }

func (exactDoubleStrategy) Validate(b engine.Bet) error { return validateTargetRange(b, 1, 6) }

func (exactDoubleStrategy) IsWin(d engine.DiceResult, b engine.Bet) bool {
	e := engine.NewDiceEvaluator(d)
	return !e.IsTriple() && e.CountOccurrences(b.TargetValue) == 2
}

func (exactDoubleStrategy) Odds(_ engine.Bet, _ engine.DiceResult) int64 { return 10 }
