package hub

import (
	"context"
	"time"

	"cachon-casino/internal/network"
)

type Inbound struct {
	Client *Client
	Msg    network.Envelope
}

type Hub struct {
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan network.Envelope
	Inbound    chan Inbound

	clients    map[*Client]struct{}
	SlowGrace  time.Duration
	dispatcher *Dispatcher
}

func New(grace time.Duration, dispatcher *Dispatcher) *Hub {
	return &Hub{
		Register:   make(chan *Client, 64),
		Unregister: make(chan *Client, 64),
		Broadcast:  make(chan network.Envelope, 256),
		Inbound:    make(chan Inbound, 256),
		clients:    map[*Client]struct{}{},
		SlowGrace:  grace,
		dispatcher: dispatcher,
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case c := <-h.Register:
			h.clients[c] = struct{}{}
		case c := <-h.Unregister:
			h.remove(c)
		case m := <-h.Broadcast:
			h.broadcast(m)
		case in := <-h.Inbound:
			h.dispatcher.Dispatch(ctx, in.Client, in.Msg)
		case <-ctx.Done():
			for c := range h.clients {
				close(c.Send)
			}
			return
		}
	}
}

func (h *Hub) remove(c *Client) {
	if _, ok := h.clients[c]; !ok {
		return
	}
	delete(h.clients, c)
	close(c.Send)
	if c.Conn != nil {
		_ = c.Conn.Close()
	}
}

func (h *Hub) broadcast(msg network.Envelope) {
	now := time.Now()
	for c := range h.clients {
		if enqueueOrPolicy(c, msg, now, h.SlowGrace) {
			h.remove(c)
		}
	}
}

func enqueueOrPolicy(c *Client, msg network.Envelope, now time.Time, grace time.Duration) bool {
	if tryEnqueue(c, msg) {
		c.SetSlowSince(time.Time{})
		return false
	}

	if msg.Type == network.MsgCountdownTick && dropOldestCountdownAndRetry(c, msg) {
		c.SetSlowSince(time.Time{})
		return false
	}

	slow := c.SlowSince()
	if slow.IsZero() {
		c.SetSlowSince(now)
		return false
	}

	if msg.Type == network.MsgRoundResult || msg.Type == network.MsgServerShutdown {
		return now.Sub(slow) > grace
	}
	return now.Sub(slow) > grace
}

func tryEnqueue(c *Client, msg network.Envelope) bool {
	select {
	case c.Send <- msg:
		return true
	default:
		return false
	}
}

func dropOldestCountdownAndRetry(c *Client, msg network.Envelope) bool {
	select {
	case old := <-c.Send:
		if old.Type != network.MsgCountdownTick {
			select {
			case c.Send <- old:
			default:
			}
			return false
		}
		return tryEnqueue(c, msg)
	default:
		return false
	}
}
