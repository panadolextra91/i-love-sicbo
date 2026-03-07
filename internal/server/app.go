package server

import (
	"context"
	"time"

	"cachon-casino/internal/engine"
	"cachon-casino/internal/hub"
	"cachon-casino/internal/network"
)

type Deps struct {
	Registry  *engine.StrategyRegistry
	Roller    engine.DiceRoller
	Wallet    engine.WalletStore
	RoundRepo engine.RoundRepo
	State     *RoundState
	BetBuffer *BetBuffer
	SessionID string
	Seq       *network.Sequence
	Hub       *hub.Hub
	Activity  *ActivityEngine
	Config    network.Config
}

func RegisterHandlers(d *hub.Dispatcher, deps Deps) {
	d.Register(network.MsgPing, func(_ context.Context, c *hub.Client, _ network.Envelope) {
		c.MarkPong()
		pong, _ := network.NewEnvelope(network.MsgPong, deps.SessionID, deps.Seq.Next(), map[string]string{"ok": "1"})
		select { case c.Send <- pong: default: }
	})

	d.Register(network.MsgPong, func(_ context.Context, c *hub.Client, _ network.Envelope) {
		c.MarkPong()
	})

	d.Register(network.MsgPlaceBet, func(_ context.Context, c *hub.Client, env network.Envelope) {
		payload, err := network.DecodePayload[network.PlaceBetPayload](env)
		if err != nil {
			return
		}

		rid, hard, open := deps.State.Snapshot()
		if !open || payload.RoundID != rid || time.Now().After(hard) {
			rej, _ := network.NewEnvelope(network.MsgBetRejected, deps.SessionID, deps.Seq.Next(), map[string]string{"reason": "HARD_CLOSED"})
			select { case c.Send <- rej: default: }
			return
		}

		bet := engine.Bet{PlayerID: c.ID, Type: engine.BetType(payload.BetType), Stake: payload.Stake, TargetValue: payload.TargetValue}
		if err := engine.AcceptBet(bet, deps.Registry, betStoreAdapter{buf: deps.BetBuffer, roundID: rid}); err != nil {
			rej, _ := network.NewEnvelope(network.MsgBetRejected, deps.SessionID, deps.Seq.Next(), map[string]string{"reason": err.Error()})
			select { case c.Send <- rej: default: }
			return
		}

		ack, _ := network.NewEnvelope(network.MsgBetAccepted, deps.SessionID, deps.Seq.Next(), map[string]string{"round_id": rid})
		select { case c.Send <- ack: default: }
	})
}

type betStoreAdapter struct {
	buf     *BetBuffer
	roundID string
}

func (b betStoreAdapter) Append(bet engine.Bet) error {
	b.buf.Add(b.roundID, bet)
	return nil
}
