package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"cachon-casino/internal/network"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := network.DefaultConfig()

	log.Println("Searching for Cachon Casino server...")
	endpoint, err := network.DiscoverEndpoint(ctx, cfg.ServiceName)
	if err != nil {
		log.Fatalf("mDNS discovery failed: %v", err)
	}

	serverAddr := fmt.Sprintf("ws://%s:%d/ws", endpoint.Host, endpoint.Port)
	log.Printf("Found server: %s. Connecting...", serverAddr)

	conn, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	defer conn.Close()

	fingerprint := "p-cli-" + uuid.NewString()[:8]
	seq := uint64(1)

	// Receiver loop
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Disconnected: %v", err)
				return
			}
			var env network.Envelope
			if err := json.Unmarshal(message, &env); err != nil {
				continue
			}

			switch env.Type {
			case network.MsgCountdownTick:
				var p network.CountdownPayload
				_ = network.DecodePayloadTo(env, &p)
				if p.SecondsLeft%5 == 0 || p.SecondsLeft <= 3 {
					fmt.Printf("\r[Round %s] Countdown: %d  ", p.RoundID, p.SecondsLeft)
				}
			case network.MsgRoundResult:
				var p network.RoundResultPayload
				_ = network.DecodePayloadTo(env, &p)
				fmt.Printf("\n>>> Result for %s: Dice %v (Total: %d)\n", p.RoundID, p.Dice, p.Dice[0]+p.Dice[1]+p.Dice[2])
			case network.MsgActivityLog:
				var p network.ActivityPayload
				_ = network.DecodePayloadTo(env, &p)
				fmt.Printf("\n[Activity] %s\n", p.Message)
			case network.MsgJoinAck:
				var p network.JoinAckPayload
				_ = network.DecodePayloadTo(env, &p)
				fmt.Printf("\nJoined successfully! Name: %s, Chips: %d\n", p.Name, p.Chips)
			case network.MsgBetAccepted:
				fmt.Println("\nBet accepted!")
			case network.MsgError:
				fmt.Printf("\n[Error] %v\n", string(env.Payload))
			}
		}
	}()

	// Automatic Join
	join, _ := network.NewEnvelope(network.MsgJoinRoom, "", seq, network.JoinRoomPayload{Fingerprint: fingerprint, Name: "Dân Chơi CLI"})
	_ = conn.WriteJSON(join)
	seq++

	// CLI Input Loop
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Commands: 'b <type> <stake> <target>' (e.g. 'b big 100 0') or 'q' to quit")
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "q" {
			return
		}

		parts := strings.Fields(input)
		if len(parts) >= 3 && parts[0] == "b" {
			betType := parts[1]
			stake := 0
			fmt.Sscanf(parts[2], "%d", &stake)
			target := 0
			if len(parts) > 3 {
				fmt.Sscanf(parts[3], "%d", &target)
			}

			bet, _ := network.NewEnvelope(network.MsgPlaceBet, "", seq, network.PlaceBetPayload{
				BetType:     betType,
				Stake:       int64(stake),
				TargetValue: target,
			})
			_ = conn.WriteJSON(bet)
			seq++
		}
	}
}
