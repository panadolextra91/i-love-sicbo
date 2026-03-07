package server

import (
	"sync"
	"time"

	"cachon-casino/internal/engine"
)

type RoundState struct {
	mu          sync.RWMutex
	RoundID     string
	HardCloseAt time.Time
	Open        bool
}

func (s *RoundState) Start(roundID string, hardCloseAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.RoundID = roundID
	s.HardCloseAt = hardCloseAt
	s.Open = true
}

func (s *RoundState) Close(roundID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.RoundID == roundID {
		s.Open = false
	}
}

func (s *RoundState) Snapshot() (roundID string, hardCloseAt time.Time, open bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.RoundID, s.HardCloseAt, s.Open
}

type BetBuffer struct {
	mu   sync.Mutex
	bets map[string][]engine.Bet
}

func NewBetBuffer() *BetBuffer {
	return &BetBuffer{bets: map[string][]engine.Bet{}}
}

func (b *BetBuffer) Add(roundID string, bet engine.Bet) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.bets[roundID] = append(b.bets[roundID], bet)
}

func (b *BetBuffer) Drain(roundID string) []engine.Bet {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := b.bets[roundID]
	delete(b.bets, roundID)
	return out
}
