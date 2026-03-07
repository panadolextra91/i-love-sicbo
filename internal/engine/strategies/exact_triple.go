package strategies

import "cachon-casino/internal/engine"

type exactTripleStrategy struct{}

func NewExactTripleStrategy() engine.BetStrategy { return exactTripleStrategy{} }

func (exactTripleStrategy) Type() engine.BetType { return engine.BetExactTriple }

func (exactTripleStrategy) Validate(b engine.Bet) error { return validateTargetRange(b, 1, 6) }

func (exactTripleStrategy) IsWin(d engine.DiceResult, b engine.Bet) bool {
	e := engine.NewDiceEvaluator(d)
	return e.IsTriple() && d[0] == b.TargetValue
}

func (exactTripleStrategy) Odds(_ engine.Bet, _ engine.DiceResult) int64 { return 180 }
