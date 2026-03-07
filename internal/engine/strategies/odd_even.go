package strategies

import "cachon-casino/internal/engine"

type evenStrategy struct{}
type oddStrategy struct{}

func NewEvenStrategy() engine.BetStrategy { return evenStrategy{} }
func NewOddStrategy() engine.BetStrategy  { return oddStrategy{} }

func (evenStrategy) Type() engine.BetType { return engine.BetEven }
func (oddStrategy) Type() engine.BetType  { return engine.BetOdd }

func (evenStrategy) Validate(b engine.Bet) error { return validateStake(b) }
func (oddStrategy) Validate(b engine.Bet) error  { return validateStake(b) }

func (evenStrategy) IsWin(d engine.DiceResult, _ engine.Bet) bool {
	e := engine.NewDiceEvaluator(d)
	return !e.IsTriple() && e.Sum()%2 == 0
}

func (oddStrategy) IsWin(d engine.DiceResult, _ engine.Bet) bool {
	e := engine.NewDiceEvaluator(d)
	return !e.IsTriple() && e.Sum()%2 == 1
}

func (evenStrategy) Odds(_ engine.Bet, _ engine.DiceResult) int64 { return 1 }
func (oddStrategy) Odds(_ engine.Bet, _ engine.DiceResult) int64  { return 1 }
