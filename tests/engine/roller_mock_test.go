package engine_test

import (
	"testing"

	"cachon-casino/internal/engine"
	"github.com/stretchr/testify/require"
)

func TestFixedDiceRoller_OneLineMock(t *testing.T) {
	roller := engine.NewFixedDiceRoller(engine.DiceResult{6, 6, 6})
	require.Equal(t, engine.DiceResult{6, 6, 6}, roller.Roll3())
}
