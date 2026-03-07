package engine_test

import (
	"testing"

	"cachon-casino/internal/engine"
	"github.com/stretchr/testify/require"
)

func TestAcceptBet_ValidateImmediately(t *testing.T) {
	reg := newRegistry()
	store := &memBetStore{}

	err := engine.AcceptBet(engine.Bet{PlayerID: "p1", Type: engine.BetExactTriple, Stake: 100, TargetValue: 9}, reg, store)
	require.Error(t, err)
	require.Len(t, store.bets, 0)

	err = engine.AcceptBet(engine.Bet{PlayerID: "p1", Type: engine.BetExactTriple, Stake: 100, TargetValue: 6}, reg, store)
	require.NoError(t, err)
	require.Len(t, store.bets, 1)
}
