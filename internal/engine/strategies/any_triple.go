package strategies

import "cachon-casino/internal/engine"

type anyTripleStrategy struct{}

func NewAnyTripleStrategy() engine.BetStrategy { return anyTripleStrategy{} }

func (anyTripleStrategy) Type() engine.BetType { return engine.BetAnyTriple }

func (anyTripleStrategy) Validate(b engine.Bet) error { return validateStake(b) }

func (anyTripleStrategy) IsWin(d engine.DiceResult, _ engine.Bet) bool {
	return engine.NewDiceEvaluator(d).IsTriple()
}

func (anyTripleStrategy) Odds(_ engine.Bet, _ engine.DiceResult) int64 { return 30 }
