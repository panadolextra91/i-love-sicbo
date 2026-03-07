package engine_test

import (
	"errors"

	"cachon-casino/internal/engine"
)

type memBetStore struct {
	bets []engine.Bet
}

func (s *memBetStore) Append(b engine.Bet) error {
	s.bets = append(s.bets, b)
	return nil
}

type walletMock struct {
	applied map[string]int64
	err     error
}

func (w *walletMock) ApplyGrossBatch(playerGross map[string]int64) error {
	if w.err != nil {
		return w.err
	}
	w.applied = map[string]int64{}
	for k, v := range playerGross {
		w.applied[k] = v
	}
	return nil
}

type roundRepoMock struct {
	called bool
	err    error
}

func (r *roundRepoMock) MarkSettled(_ string, _ engine.DiceResult, _ []engine.PayoutResult) error {
	if r.err != nil {
		return r.err
	}
	r.called = true
	return nil
}

var errWallet = errors.New("wallet failed")
