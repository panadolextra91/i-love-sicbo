package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cachon-casino/internal/engine"
	"cachon-casino/internal/engine/strategies"
	"cachon-casino/internal/hub"
	"cachon-casino/internal/network"
	"cachon-casino/internal/repo"
	"cachon-casino/internal/server"
	"cachon-casino/internal/storage"

	"github.com/google/uuid"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := network.DefaultConfig()
	sessionID := uuid.NewString()

	// 1. Storage
	db, err := storage.OpenSQLite(ctx, "casino.db")
	if err != nil {
		log.Fatalf("failed to open sqlite: %v", err)
	}
	defer db.Close()

	repository := repo.New(db)

	// 2. Engine
	reg := engine.NewStrategyRegistry(
		strategies.NewBigStrategy(),
		strategies.NewSmallStrategy(),
		strategies.NewOddStrategy(),
		strategies.NewEvenStrategy(),
		strategies.NewAnyTripleStrategy(),
		strategies.NewExactTripleStrategy(),
		strategies.NewExactDoubleStrategy(),
		strategies.NewTwoNumberComboStrategy(),
		strategies.NewExactTotalStrategy(),
		strategies.NewSingleNumberStrategy(),
	)
	roller := engine.CryptoDiceRoller{}

	// 3. Hub & Dispatcher
	dispatcher := hub.NewDispatcher()
	h := hub.New(cfg.SlowClientGrace, dispatcher)
	go h.Run(ctx)

	seq := &network.Sequence{}

	deps := server.Deps{
		Registry:     reg,
		Roller:       roller,
		SettleStore:  repository,
		PlayerRepo:   repository,
		State:        server.NewRoundState(),
		BetBuffer:    server.NewBetBuffer(),
		SessionID:    sessionID,
		Seq:          seq,
		Hub:          h,
		Activity:     server.NewActivityEngine(),
		Config:       cfg,
		CurrentRound: 1,
	}

	server.RegisterHandlers(dispatcher, deps)
	// dispatcher.Run is not needed as Hub.Run calls dispatcher.Dispatch

	// 4. mDNS
	stopMDNS, err := network.StartMDNSAdvertise(cfg)
	if err != nil {
		log.Printf("failed to start mDNS advertise: %v", err)
	} else {
		defer stopMDNS.Stop()
		log.Printf("mDNS advertising started: %s on port %d", cfg.ServiceName, cfg.WSPort)
	}

	// 5. HTTP/WS Server
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.ServeWS(h, cfg, w, r)
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.WSPort),
		Handler: mux,
	}

	go func() {
		log.Printf("WebSocket server starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 6. Game Loop
	go server.RunRoundLoop(ctx, deps, 15)

	log.Println("Server is running. Press Ctrl+C to stop.")
	<-ctx.Done()

	log.Println("Shutting down...")
	server.GracefulShutdown(deps, srv, 5*time.Second)
	log.Println("Server stopped.")
}
