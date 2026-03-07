package engine_test

import (
	"testing"

	"cachon-casino/internal/engine"
	"cachon-casino/internal/engine/strategies"
	"github.com/stretchr/testify/require"
)

func newRegistry() *engine.StrategyRegistry {
	return engine.NewStrategyRegistry(
		strategies.NewBigStrategy(),
		strategies.NewSmallStrategy(),
		strategies.NewOddStrategy(),
		strategies.NewEvenStrategy(),
		strategies.NewAnyTripleStrategy(),
		strategies.NewExactTripleStrategy(),
		strategies.NewExactDoubleStrategy(),
		strategies.NewTwoNumberComboStrategy(),
		strategies.NewExactTotalStrategy(),
		strategies.NewSingleNumberStrategy(),
	)
}

func TestRegistryGet(t *testing.T) {
	reg := newRegistry()
	require.NotNil(t, reg.Get(engine.BetSmall))
	require.NotNil(t, reg.Get(engine.BetExactTriple))
}
