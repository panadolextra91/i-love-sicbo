package server

import (
	"context"
	"fmt"
	"time"

	"cachon-casino/internal/engine"
	"cachon-casino/internal/network"
	"cachon-casino/internal/repo"
)

func RunRoundLoop(ctx context.Context, deps Deps, betWindowSec int) {
	roundNo := int64(1)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if err := deps.SettleStore.AuditChipLedger(ctx); err != nil {
			alert, _ := network.NewEnvelope(network.MsgActivityLog, deps.SessionID, deps.Seq.Next(), network.ActivityPayload{Message: "[ALERT] Ledger audit mismatch, betting is paused"})
			deps.Hub.Broadcast <- alert
			time.Sleep(1 * time.Second)
			continue
		}

		roundID := fmt.Sprintf("r-%d", roundNo)
		startedAt := time.Now()
		hard := startedAt.Add(time.Duration(betWindowSec) * time.Second).Add(deps.Config.LatencyBuffer)
		deps.State.Start(roundID, hard)

		for sec := betWindowSec; sec >= 0; sec-- {
			env, _ := network.NewEnvelope(network.MsgCountdownTick, deps.SessionID, deps.Seq.Next(), network.CountdownPayload{RoundID: roundID, SecondsLeft: sec, BettingOpen: sec > 0})
			deps.Hub.Broadcast <- env
			time.Sleep(1 * time.Second)
		}

		time.Sleep(deps.Config.LatencyBuffer)
		deps.State.Close(roundID)

		bets := deps.BetBuffer.Drain(roundID)
		settledAt := time.Now()
		settlement, err := engine.SettleRoundAtomic(roundID, roundNo, startedAt.UnixMilli(), settledAt.UnixMilli(), bets, deps.Roller, deps.Registry, settleStoreAdapter{repo: deps.SettleStore})
		if err == nil {
			payload := network.RoundResultPayload{RoundID: settlement.RoundID, Dice: settlement.Dice, Settlements: settlement.Details}
			res, _ := network.NewEnvelope(network.MsgRoundResult, deps.SessionID, deps.Seq.Next(), payload)
			deps.Hub.Broadcast <- res

			for _, line := range deps.Activity.Build(settlement) {
				act, _ := network.NewEnvelope(network.MsgActivityLog, deps.SessionID, deps.Seq.Next(), network.ActivityPayload{Message: line})
				deps.Hub.Broadcast <- act
			}
		}

		roundNo++
	}
}

type settleStoreAdapter struct {
	repo *repo.Repository
}

func (s settleStoreAdapter) SettleRound(roundID string, roundNo int64, startedAt, settledAt int64, dice engine.DiceResult, results []engine.PayoutResult) error {
	return s.repo.SettleRound(context.Background(), roundID, roundNo, time.UnixMilli(startedAt), time.UnixMilli(settledAt), dice, results)
}
