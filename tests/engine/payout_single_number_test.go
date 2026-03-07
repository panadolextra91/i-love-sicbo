package engine_test

import (
	"testing"

	"cachon-casino/internal/engine"
	"github.com/stretchr/testify/require"
)

func TestPayoutSingleNumber(t *testing.T) {
	reg := newRegistry()

	once := engine.CalculatePayout(reg, []int{3, 1, 2}, engine.Bet{Type: engine.BetSingleNumber, Stake: 100, TargetValue: 3})
	require.Equal(t, int64(200), once)

	twice := engine.CalculatePayout(reg, []int{3, 3, 2}, engine.Bet{Type: engine.BetSingleNumber, Stake: 100, TargetValue: 3})
	require.Equal(t, int64(300), twice)

	three := engine.CalculatePayout(reg, []int{3, 3, 3}, engine.Bet{Type: engine.BetSingleNumber, Stake: 100, TargetValue: 3})
	require.Equal(t, int64(400), three)
}
