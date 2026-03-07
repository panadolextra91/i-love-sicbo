package engine_test

import (
	"errors"
	"testing"

	"cachon-casino/internal/engine"

	"github.com/stretchr/testify/require"
)

func TestSettleRoundAtomic_GrossAndCommitOrder(t *testing.T) {
	reg := newRegistry()
	repo := &roundRepoMock{}
	roller := engine.NewFixedDiceRoller(engine.DiceResult{6, 6, 6})

	bets := []engine.Bet{
		{PlayerID: "p1", Type: engine.BetAnyTriple, Stake: 100},
		{PlayerID: "p2", Type: engine.BetBig, Stake: 100},
	}

	settlement, err := engine.SettleRoundAtomic("r1", 1, 0, 0, bets, roller, reg, repo)
	require.NoError(t, err)
	require.True(t, repo.called)
	// wallet is not used in engine.SettleRoundAtomic anymore, it's inside SettlementStore.SettleRound (repo)
	// since we are testing engine.SettleRoundAtomic, we check the settlement result
	require.Equal(t, engine.DiceResult{6, 6, 6}, settlement.Dice)
	require.Equal(t, int64(3100), settlement.PlayerGross["p1"])
}

func TestSettleRoundAtomic_WalletFailNoRoundCommit(t *testing.T) {
	reg := newRegistry()
	repo := &roundRepoMock{err: errors.New("storage failed")}
	roller := engine.NewFixedDiceRoller(engine.DiceResult{2, 2, 2})

	_, err := engine.SettleRoundAtomic("r2", 2, 0, 0, []engine.Bet{{PlayerID: "p1", Type: engine.BetEven, Stake: 100}}, roller, reg, repo)
	require.Error(t, err)
	require.False(t, repo.called) // actual implementation sets called=true only on success
}
