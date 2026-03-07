package strategies

import "cachon-casino/internal/engine"

type bigStrategy struct{}
type smallStrategy struct{}

func NewBigStrategy() engine.BetStrategy { return bigStrategy{} }
func NewSmallStrategy() engine.BetStrategy { return smallStrategy{} }

func (bigStrategy) Type() engine.BetType { return engine.BetBig }
func (smallStrategy) Type() engine.BetType { return engine.BetSmall }

func (bigStrategy) Validate(b engine.Bet) error { return validateStake(b) }
func (smallStrategy) Validate(b engine.Bet) error { return validateStake(b) }

func (bigStrategy) IsWin(d engine.DiceResult, _ engine.Bet) bool {
	e := engine.NewDiceEvaluator(d)
	if e.IsTriple() {
		return false
	}
	s := e.Sum()
	return s >= 11 && s <= 17
}

func (smallStrategy) IsWin(d engine.DiceResult, _ engine.Bet) bool {
	e := engine.NewDiceEvaluator(d)
	if e.IsTriple() {
		return false
	}
	s := e.Sum()
	return s >= 4 && s <= 10
}

func (bigStrategy) Odds(_ engine.Bet, _ engine.DiceResult) int64 { return 1 }
func (smallStrategy) Odds(_ engine.Bet, _ engine.DiceResult) int64 { return 1 }
