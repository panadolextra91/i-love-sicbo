package engine

import "errors"

var (
	ErrUnknownBetType  = errors.New("unknown bet type")
	ErrInvalidStake    = errors.New("stake must be > 0")
	ErrInvalidTarget   = errors.New("invalid target value")
	ErrInvalidDiceSize = errors.New("dice must contain exactly 3 values")
)
