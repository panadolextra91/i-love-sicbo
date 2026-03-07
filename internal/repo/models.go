package repo

import (
	"time"

	"cachon-casino/internal/engine"
)

type Player struct {
	ID          string
	Name        string
	Fingerprint string
	Chips       int64
	LastBonusAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type GameRound struct {
	ID        string
	RoundNo   int64
	Dice      engine.DiceResult
	Total     int
	StartedAt time.Time
	SettledAt time.Time
}

type BetRecord struct {
	ID          string
	RoundID     string
	PlayerID    string
	BetType     string
	Amount      int64
	Payout      int64
	IsWin       bool
	CreatedAt   time.Time
	TargetValue int
}

type ChipTransaction struct {
	ID        string
	PlayerID  string
	Amount    int64
	Reason    string
	RoundID   *string
	CreatedAt time.Time
}
