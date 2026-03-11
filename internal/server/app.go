package server

import (
	"context"
	"time"

	"cachon-casino/internal/engine"
	"cachon-casino/internal/hub"
	"cachon-casino/internal/network"
	"cachon-casino/internal/repo"
)

type Deps struct {
	Registry     *engine.StrategyRegistry
	Roller       engine.DiceRoller
	SettleStore  *repo.Repository
	PlayerRepo   *repo.Repository
	State        *RoundState
	BetBuffer    *BetBuffer
	Barrier      *ReadyBarrier
	SessionID    string
	Seq          *network.Sequence
	Hub          *hub.Hub
	Activity     *ActivityEngine
	Config       network.Config
	CurrentRound int64
}

func RegisterHandlers(d *hub.Dispatcher, deps Deps) {
	d.Register(network.MsgPing, func(_ context.Context, c *hub.Client, _ network.Envelope) {
		c.MarkPong()
		pong, _ := network.NewEnvelope(network.MsgPong, deps.SessionID, deps.Seq.Next(), map[string]string{"ok": "1"})
		select {
		case c.Send <- pong:
		default:
		}
	})

	d.Register(network.MsgPong, func(_ context.Context, c *hub.Client, _ network.Envelope) {
		c.MarkPong()
	})

	d.Register(network.MsgJoinRoom, func(ctx context.Context, c *hub.Client, env network.Envelope) {
		payload, err := network.DecodePayload[network.JoinRoomPayload](env)
		if err != nil || payload.Fingerprint == "" {
			return
		}
		name := payload.Name
		if name == "" {
			name = "anonymous"
		}
		player, bonus, err := deps.PlayerRepo.RegisterOrLoadPlayer(ctx, payload.Fingerprint, name, time.Now())
		if err != nil {
			errEnv, _ := network.NewEnvelope(network.MsgError, deps.SessionID, deps.Seq.Next(), map[string]string{"reason": err.Error()})
			select {
			case c.Send <- errEnv:
			default:
			}
			return
		}
		c.PlayerID = player.ID
		c.DisplayName = player.Name
		c.Fingerprint = player.Fingerprint
		roundID, _, _ := deps.State.Snapshot()
		ack, _ := network.NewEnvelope(network.MsgJoinAck, deps.SessionID, deps.Seq.Next(), network.JoinAckPayload{PlayerID: player.ID, Name: player.Name, Chips: player.Chips, Bonus: bonus, RoundID: roundID})
		select {
		case c.Send <- ack:
		default:
		}
	})

	d.Register(network.MsgPlaceBet, func(_ context.Context, c *hub.Client, env network.Envelope) {
		payload, err := network.DecodePayload[network.PlaceBetPayload](env)
		if err != nil {
			return
		}
		if c.PlayerID == "" {
			rej, _ := network.NewEnvelope(network.MsgBetRejected, deps.SessionID, deps.Seq.Next(), map[string]string{"reason": "NOT_JOINED"})
			select {
			case c.Send <- rej:
			default:
			}
			return
		}

		rid, hard, open := deps.State.Snapshot()
		if !open || payload.RoundID != rid || time.Now().After(hard) {
			rej, _ := network.NewEnvelope(network.MsgBetRejected, deps.SessionID, deps.Seq.Next(), map[string]string{"reason": "HARD_CLOSED"})
			select {
			case c.Send <- rej:
			default:
			}
			return
		}

		bet := engine.Bet{PlayerID: c.PlayerID, Type: engine.BetType(payload.BetType), Stake: payload.Stake, TargetValue: payload.TargetValue}
		if err := engine.AcceptBet(bet, deps.Registry, betStoreAdapter{buf: deps.BetBuffer, roundID: rid}); err != nil {
			rej, _ := network.NewEnvelope(network.MsgBetRejected, deps.SessionID, deps.Seq.Next(), map[string]string{"reason": err.Error()})
			select {
			case c.Send <- rej:
			default:
			}
			return
		}

		ack, _ := network.NewEnvelope(network.MsgBetAccepted, deps.SessionID, deps.Seq.Next(), map[string]string{"round_id": rid})
		select {
		case c.Send <- ack:
		default:
		}
	})

	d.Register(network.MsgRoundReady, func(_ context.Context, c *hub.Client, env network.Envelope) {
		payload, err := network.DecodePayload[network.RoundReadyPayload](env)
		if err != nil {
			return
		}
		deps.Barrier.Mark(payload.RoundID, c.PlayerID)
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
