package engine_test

import (
	"testing"

	"cachon-casino/internal/engine"
	"github.com/stretchr/testify/require"
)

func TestSettleRoundAtomic_GrossAndCommitOrder(t *testing.T) {
	reg := newRegistry()
	wallet := &walletMock{}
	repo := &roundRepoMock{}
	roller := engine.NewFixedDiceRoller(engine.DiceResult{6, 6, 6})

	bets := []engine.Bet{
		{PlayerID: "p1", Type: engine.BetAnyTriple, Stake: 100},
		{PlayerID: "p2", Type: engine.BetBig, Stake: 100},
	}

	settlement, err := engine.SettleRoundAtomic("r1", bets, roller, reg, wallet, repo)
	require.NoError(t, err)
	require.True(t, repo.called)
	require.Equal(t, int64(3100), wallet.applied["p1"])
	require.Equal(t, int64(0), wallet.applied["p2"])
	require.Equal(t, engine.DiceResult{6, 6, 6}, settlement.Dice)
}

func TestSettleRoundAtomic_WalletFailNoRoundCommit(t *testing.T) {
	reg := newRegistry()
	wallet := &walletMock{err: errWallet}
	repo := &roundRepoMock{}
	roller := engine.NewFixedDiceRoller(engine.DiceResult{2, 2, 2})

	_, err := engine.SettleRoundAtomic("r2", []engine.Bet{{PlayerID: "p1", Type: engine.BetEven, Stake: 100}}, roller, reg, wallet, repo)
	require.Error(t, err)
	require.False(t, repo.called)
}
