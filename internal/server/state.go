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

func NewRoundState() *RoundState {
	return &RoundState{}
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

type ReadyBarrier struct {
	mu      sync.Mutex
	roundID string
	ready   map[string]struct{}
}

func NewReadyBarrier() *ReadyBarrier {
	return &ReadyBarrier{ready: map[string]struct{}{}}
}

func (b *ReadyBarrier) Reset(roundID string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.roundID = roundID
	b.ready = map[string]struct{}{}
}

func (b *ReadyBarrier) Mark(roundID, playerID string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.roundID == roundID && playerID != "" {
		b.ready[playerID] = struct{}{}
	}
}

func (b *ReadyBarrier) ReadyCount() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.ready)
}
