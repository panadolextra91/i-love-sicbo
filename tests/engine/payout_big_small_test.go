package engine_test

import (
	"testing"

	"cachon-casino/internal/engine"
	"github.com/stretchr/testify/require"
)

func TestPayoutBigAndSmall(t *testing.T) {
	reg := newRegistry()

	grossBig := engine.CalculatePayout(reg, []int{6, 5, 2}, engine.Bet{Type: engine.BetBig, Stake: 100})
	require.Equal(t, int64(200), grossBig)

	grossSmall := engine.CalculatePayout(reg, []int{1, 2, 2}, engine.Bet{Type: engine.BetSmall, Stake: 100})
	require.Equal(t, int64(200), grossSmall)
}
