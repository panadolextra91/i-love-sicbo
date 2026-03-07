package strategies

import "cachon-casino/internal/engine"

func validateStake(b engine.Bet) error {
	if b.Stake <= 0 {
		return engine.ErrInvalidStake
	}
	return nil
}

func validateTargetRange(b engine.Bet, min int, max int) error {
	if err := validateStake(b); err != nil {
		return err
	}
	if b.TargetValue < min || b.TargetValue > max {
		return engine.ErrInvalidTarget
	}
	return nil
}
