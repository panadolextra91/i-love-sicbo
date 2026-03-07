package engine_test

import (
	"testing"

	"cachon-casino/internal/engine"
	"github.com/stretchr/testify/require"
)

func TestTriple111_SmallAndOddMustLoseBlank(t *testing.T) {
	reg := newRegistry()
	dice := []int{1, 1, 1}

	small := engine.CalculatePayout(reg, dice, engine.Bet{Type: engine.BetSmall, Stake: 100})
	odd := engine.CalculatePayout(reg, dice, engine.Bet{Type: engine.BetOdd, Stake: 100})

	require.Equal(t, int64(0), small)
	require.Equal(t, int64(0), odd)
}
