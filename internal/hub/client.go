package hub

import (
	"sync"
	"time"

	"cachon-casino/internal/network"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID    string
	Conn  *websocket.Conn
	Send  chan network.Envelope
	mu    sync.RWMutex
	pong  time.Time
	slow  time.Time
}

func NewClient(id string, conn *websocket.Conn, buffer int) *Client {
	return &Client{ID: id, Conn: conn, Send: make(chan network.Envelope, buffer), pong: time.Now()}
}

func (c *Client) MarkPong() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.pong = time.Now()
}

func (c *Client) LastPong() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.pong
}

func (c *Client) SetSlowSince(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.slow = t
}

func (c *Client) SlowSince() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.slow
}
