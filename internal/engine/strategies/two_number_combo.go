package strategies

import "cachon-casino/internal/engine"

type twoNumberComboStrategy struct{}

func NewTwoNumberComboStrategy() engine.BetStrategy { return twoNumberComboStrategy{} }

func (twoNumberComboStrategy) Type() engine.BetType { return engine.BetTwoNumberCombo }

func (twoNumberComboStrategy) Validate(b engine.Bet) error {
	if err := validateStake(b); err != nil {
		return err
	}
	a, c := engine.DecodeTwoNumberCombo(b.TargetValue)
	if a < 1 || a > 6 || c < 1 || c > 6 || a >= c {
		return engine.ErrInvalidTarget
	}
	return nil
}

func (twoNumberComboStrategy) IsWin(d engine.DiceResult, b engine.Bet) bool {
	a, c := engine.DecodeTwoNumberCombo(b.TargetValue)
	e := engine.NewDiceEvaluator(d)
	return e.HasBoth(a, c)
}

func (twoNumberComboStrategy) Odds(_ engine.Bet, _ engine.DiceResult) int64 { return 6 }
