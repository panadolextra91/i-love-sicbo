package server

import (
	"context"
	"fmt"
	"time"

	"cachon-casino/internal/engine"
	"cachon-casino/internal/network"
)

func RunRoundLoop(ctx context.Context, deps Deps, betWindowSec int) {
	roundNo := 1
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		roundID := fmt.Sprintf("r-%d", roundNo)
		now := time.Now()
		hard := now.Add(time.Duration(betWindowSec) * time.Second).Add(deps.Config.LatencyBuffer)
		deps.State.Start(roundID, hard)

		for sec := betWindowSec; sec >= 0; sec-- {
			env, _ := network.NewEnvelope(network.MsgCountdownTick, deps.SessionID, deps.Seq.Next(), network.CountdownPayload{RoundID: roundID, SecondsLeft: sec, BettingOpen: sec > 0})
			deps.Hub.Broadcast <- env
			time.Sleep(1 * time.Second)
		}

		time.Sleep(deps.Config.LatencyBuffer)
		deps.State.Close(roundID)

		bets := deps.BetBuffer.Drain(roundID)
		settlement, err := engine.SettleRoundAtomic(roundID, bets, deps.Roller, deps.Registry, deps.Wallet, deps.RoundRepo)
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
