package strategies

import "cachon-casino/internal/engine"

type exactTotalStrategy struct{}

var totalOdds = map[int]int64{
	4:  60,
	5:  20,
	6:  18,
	7:  12,
	8:  8,
	9:  7,
	10: 6,
	11: 6,
	12: 7,
	13: 8,
	14: 12,
	15: 18,
	16: 20,
	17: 60,
}

func NewExactTotalStrategy() engine.BetStrategy { return exactTotalStrategy{} }

func (exactTotalStrategy) Type() engine.BetType { return engine.BetExactTotal }

func (exactTotalStrategy) Validate(b engine.Bet) error {
	if err := validateStake(b); err != nil {
		return err
	}
	if _, ok := totalOdds[b.TargetValue]; !ok {
		return engine.ErrInvalidTarget
	}
	return nil
}

func (exactTotalStrategy) IsWin(d engine.DiceResult, b engine.Bet) bool {
	return engine.NewDiceEvaluator(d).Sum() == b.TargetValue
}

func (exactTotalStrategy) Odds(b engine.Bet, _ engine.DiceResult) int64 {
	return totalOdds[b.TargetValue]
}
