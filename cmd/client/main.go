package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"cachon-casino/internal/network"
	"cachon-casino/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := network.DefaultConfig()

	log.Println("Searching for Cachon Casino server...")
	var serverAddr string
	endpoint, err := network.DiscoverEndpoint(ctx, cfg.ServiceName)
	if err != nil {
		log.Printf("mDNS discovery failed: %v. Falling back to localhost...", err)
		serverAddr = "ws://localhost:8080/ws"
	} else {
		serverAddr = "ws://" + endpoint.Host + ":8080/ws"
		log.Printf("Found server via mDNS: %s", serverAddr)
	}

	log.Printf("Connecting to %s...", serverAddr)
	conn, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	defer conn.Close()

	inbound := make(chan network.Envelope, 32)
	outbound := make(chan network.Envelope, 32)

	// Receiver loop
	go func() {
		for {
			var env network.Envelope
			err := conn.ReadJSON(&env)
			if err != nil {
				return
			}
			inbound <- env
		}
	}()

	// Sender loop
	go func() {
		for env := range outbound {
			conn.WriteJSON(env)
		}
	}()

	sessionID := uuid.NewString()

	// Initial Join
	join, _ := network.NewEnvelope(network.MsgJoinRoom, "", 0, network.JoinRoomPayload{
		Fingerprint: "tui-" + sessionID[:8],
		Name:        "Dân Chơi Hệ Hồng",
	})
	outbound <- join

	m := tui.Model{
		State:     tui.StateJoining,
		SessionID: sessionID,
		Config:    cfg,
		Inbound:   inbound,
		Outbound:  outbound,
		Ctx:       ctx,
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
