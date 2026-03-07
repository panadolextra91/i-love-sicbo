package server

import (
	"net/http"
	"time"

	"cachon-casino/internal/hub"
	"cachon-casino/internal/network"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func ServeWS(h *hub.Hub, cfg network.Config, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := hub.NewClient(uuid.NewString(), conn, cfg.ClientSendBuffer)
	h.Register <- client

	go writePump(h, client, cfg)
	go readPump(h, client, cfg)
}

func readPump(h *hub.Hub, c *hub.Client, cfg network.Config) {
	defer func() { h.Unregister <- c }()

	_ = c.Conn.SetReadDeadline(time.Now().Add(cfg.PongWait))
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(cfg.PongWait))
		return nil
	})

	for {
		var env network.Envelope
		if err := c.Conn.ReadJSON(&env); err != nil {
			return
		}
		if env.Type == network.MsgPong {
			c.MarkPong()
		}
		h.Inbound <- hub.Inbound{Client: c, Msg: env}
	}
}

func writePump(h *hub.Hub, c *hub.Client, cfg network.Config) {
	wsTicker := time.NewTicker(cfg.PingInterval)
	appTicker := time.NewTicker(cfg.AppPingInterval)
	defer func() {
		wsTicker.Stop()
		appTicker.Stop()
		h.Unregister <- c
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(cfg.PongWait))
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteJSON(msg); err != nil {
				return
			}
		case <-appTicker.C:
			if time.Since(c.LastPong()) > cfg.HeartbeatTimeout {
				return
			}
			pingEnv, _ := network.NewEnvelope(network.MsgPing, "", 0, map[string]string{"kind": "app"})
			_ = c.Conn.SetWriteDeadline(time.Now().Add(cfg.PongWait))
			if err := c.Conn.WriteJSON(pingEnv); err != nil {
				return
			}
		case <-wsTicker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(cfg.PongWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
