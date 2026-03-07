package hub

import (
	"context"

	"cachon-casino/internal/network"
)

type HandlerFunc func(context.Context, *Client, network.Envelope)

type Dispatcher struct {
	handlers map[network.MessageType]HandlerFunc
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{handlers: map[network.MessageType]HandlerFunc{}}
}

func (d *Dispatcher) Register(t network.MessageType, h HandlerFunc) {
	d.handlers[t] = h
}

func (d *Dispatcher) Dispatch(ctx context.Context, c *Client, msg network.Envelope) {
	h, ok := d.handlers[msg.Type]
	if !ok {
		return
	}
	h(ctx, c, msg)
}
