package engine_test

import (
	"testing"

	"cachon-casino/internal/engine"
	"github.com/stretchr/testify/require"
)

func TestPayoutExactDoubleAndCombo(t *testing.T) {
	reg := newRegistry()

	doubleWin := engine.CalculatePayout(reg, []int{2, 2, 5}, engine.Bet{Type: engine.BetExactDouble, Stake: 100, TargetValue: 2})
	require.Equal(t, int64(1100), doubleWin)

	doubleLoseOnTriple := engine.CalculatePayout(reg, []int{2, 2, 2}, engine.Bet{Type: engine.BetExactDouble, Stake: 100, TargetValue: 2})
	require.Equal(t, int64(0), doubleLoseOnTriple)

	comboWin := engine.CalculatePayout(reg, []int{1, 2, 6}, engine.Bet{Type: engine.BetTwoNumberCombo, Stake: 100, TargetValue: 12})
	require.Equal(t, int64(700), comboWin)
}
