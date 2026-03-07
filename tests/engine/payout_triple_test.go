package engine_test

import (
	"testing"

	"cachon-casino/internal/engine"
	"github.com/stretchr/testify/require"
)

func TestPayoutAnyAndExactTriple(t *testing.T) {
	reg := newRegistry()

	anyTriple := engine.CalculatePayout(reg, []int{6, 6, 6}, engine.Bet{Type: engine.BetAnyTriple, Stake: 100})
	require.Equal(t, int64(3100), anyTriple)

	exactTriple := engine.CalculatePayout(reg, []int{6, 6, 6}, engine.Bet{Type: engine.BetExactTriple, Stake: 100, TargetValue: 6})
	require.Equal(t, int64(18100), exactTriple)

	wrongExact := engine.CalculatePayout(reg, []int{6, 6, 6}, engine.Bet{Type: engine.BetExactTriple, Stake: 100, TargetValue: 5})
	require.Equal(t, int64(0), wrongExact)
}
