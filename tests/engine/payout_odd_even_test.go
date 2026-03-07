package engine_test

import (
	"testing"

	"cachon-casino/internal/engine"
	"github.com/stretchr/testify/require"
)

func TestPayoutOddEven(t *testing.T) {
	reg := newRegistry()

	grossEven := engine.CalculatePayout(reg, []int{2, 2, 2}, engine.Bet{Type: engine.BetEven, Stake: 100})
	require.Equal(t, int64(0), grossEven)

	grossOdd := engine.CalculatePayout(reg, []int{2, 2, 1}, engine.Bet{Type: engine.BetOdd, Stake: 100})
	require.Equal(t, int64(200), grossOdd)
}
