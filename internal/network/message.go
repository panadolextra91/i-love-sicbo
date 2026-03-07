package network

import (
	"encoding/json"
	"sync/atomic"
	"time"
)

type MessageType string

const (
	MsgJoinRoom            MessageType = "JOIN_ROOM"
	MsgPlaceBet            MessageType = "PLACE_BET"
	MsgPing                MessageType = "PING"
	MsgPong                MessageType = "PONG"
	MsgJoinAck             MessageType = "JOIN_ACK"
	MsgBetAccepted         MessageType = "BET_ACCEPTED"
	MsgBetRejected         MessageType = "BET_REJECTED"
	MsgCountdownTick       MessageType = "COUNTDOWN_TICK"
	MsgRoundResult         MessageType = "ROUND_RESULT"
	MsgLeaderboardSnapshot MessageType = "LEADERBOARD_SNAPSHOT"
	MsgActivityLog         MessageType = "ACTIVITY_LOG"
	MsgServerShutdown      MessageType = "SERVER_SHUTDOWN"
	MsgError               MessageType = "ERROR"
)

type Envelope struct {
	Type      MessageType     `json:"type"`
	SessionID string          `json:"session_id"`
	Sequence  uint64          `json:"sequence"`
	RequestID string          `json:"request_id,omitempty"`
	SentAt    int64           `json:"sent_at"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

type Sequence struct{ v atomic.Uint64 }

func (s *Sequence) Next() uint64 { return s.v.Add(1) }

func NewEnvelope(t MessageType, sessionID string, seq uint64, payload any) (Envelope, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return Envelope{}, err
	}
	return Envelope{Type: t, SessionID: sessionID, Sequence: seq, SentAt: time.Now().UnixMilli(), Payload: raw}, nil
}

type CountdownPayload struct {
	RoundID     string `json:"round_id"`
	SecondsLeft int    `json:"seconds_left"`
	BettingOpen bool   `json:"betting_open"`
}

type JoinRoomPayload struct {
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
}

type JoinAckPayload struct {
	PlayerID string `json:"player_id"`
	Name     string `json:"name"`
	Chips    int64  `json:"chips"`
	Bonus    int64  `json:"bonus"`
	RoundID  string `json:"round_id,omitempty"`
}

type PlaceBetPayload struct {
	RoundID     string `json:"round_id"`
	BetType     string `json:"bet_type"`
	Stake       int64  `json:"stake"`
	TargetValue int    `json:"target_value"`
}

type RoundResultPayload struct {
	RoundID       string      `json:"round_id"`
	Dice          [3]int      `json:"dice"`
	Settlements   interface{} `json:"settlements"`
	Leaderboard   interface{} `json:"leaderboard,omitempty"`
	TopActivities []string    `json:"top_activities,omitempty"`
}

type ShutdownPayload struct {
	Reason     string `json:"reason"`
	ShutdownAt int64  `json:"shutdown_at"`
}

type ActivityPayload struct {
	Message string `json:"message"`
}
