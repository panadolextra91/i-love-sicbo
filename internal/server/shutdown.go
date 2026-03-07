package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"cachon-casino/internal/network"
)

func BroadcastShutdown(deps Deps, reason string, drain time.Duration) {
	env, _ := network.NewEnvelope(network.MsgServerShutdown, deps.SessionID, deps.Seq.Next(), network.ShutdownPayload{
		Reason:     reason,
		ShutdownAt: time.Now().Add(drain).UnixMilli(),
	})
	deps.Hub.Broadcast <- env
}

func GracefulShutdown(deps Deps, srv *http.Server, timeout time.Duration) {
	log.Println("Broadcasting shutdown to clients...")
	BroadcastShutdown(deps, "Server shutting down", timeout/2)

	time.Sleep(timeout / 2) // Wait for clients to receive message

	log.Println("Stopping HTTP server...")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}
}
