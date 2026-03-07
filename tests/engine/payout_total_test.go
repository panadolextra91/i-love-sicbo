package engine_test

import (
	"testing"

	"cachon-casino/internal/engine"
	"github.com/stretchr/testify/require"
)

func TestPayoutExactTotal(t *testing.T) {
	reg := newRegistry()

	p4 := engine.CalculatePayout(reg, []int{1, 1, 2}, engine.Bet{Type: engine.BetExactTotal, Stake: 100, TargetValue: 4})
	require.Equal(t, int64(6100), p4)

	p11 := engine.CalculatePayout(reg, []int{4, 3, 4}, engine.Bet{Type: engine.BetExactTotal, Stake: 100, TargetValue: 11})
	require.Equal(t, int64(700), p11)
}
