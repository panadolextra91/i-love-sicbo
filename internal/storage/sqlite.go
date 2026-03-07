package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func OpenSQLite(ctx context.Context, dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	pragmas := []string{
		"PRAGMA journal_mode=WAL;",
		"PRAGMA synchronous=NORMAL;",
		"PRAGMA busy_timeout=5000;",
		"PRAGMA foreign_keys=ON;",
	}
	for _, q := range pragmas {
		if _, err := db.ExecContext(ctx, q); err != nil {
			_ = db.Close()
			return nil, err
		}
	}

	if err := migrate(ctx, db); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func migrate(ctx context.Context, db *sql.DB) error {
	var version int
	if err := db.QueryRowContext(ctx, "PRAGMA user_version;").Scan(&version); err != nil {
		return err
	}

	for version < latestSchemaVersion {
		next := version + 1
		sqls, ok := migrations[next]
		if !ok {
			return fmt.Errorf("missing migration for version %d", next)
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		for _, stmt := range strings.Split(sqls, ";") {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if _, err := tx.ExecContext(ctx, stmt); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		if _, err := tx.ExecContext(ctx, fmt.Sprintf("PRAGMA user_version=%d", next)); err != nil {
			_ = tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		version = next
	}

	_, _ = db.ExecContext(ctx, "PRAGMA wal_checkpoint(PASSIVE);")
	db.SetConnMaxLifetime(10 * time.Minute)
	return nil
}

const latestSchemaVersion = 2

var migrations = map[int]string{
	1: `
CREATE TABLE IF NOT EXISTS players (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  chips BIGINT NOT NULL DEFAULT 0,
  last_bonus_at DATETIME NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);
CREATE TABLE IF NOT EXISTS game_rounds (
  id TEXT PRIMARY KEY,
  round_no BIGINT NOT NULL,
  dice_1 INTEGER NOT NULL,
  dice_2 INTEGER NOT NULL,
  dice_3 INTEGER NOT NULL,
  total INTEGER NOT NULL,
  started_at DATETIME NOT NULL,
  settled_at DATETIME NOT NULL
);
CREATE TABLE IF NOT EXISTS bets (
  id TEXT PRIMARY KEY,
  round_id TEXT NOT NULL REFERENCES game_rounds(id),
  player_id TEXT NOT NULL REFERENCES players(id),
  bet_type TEXT NOT NULL,
  amount BIGINT NOT NULL,
  payout BIGINT NOT NULL DEFAULT 0,
  is_win BOOLEAN NOT NULL DEFAULT FALSE,
  created_at DATETIME NOT NULL,
  UNIQUE(round_id, player_id, bet_type, amount, created_at)
);
CREATE TABLE IF NOT EXISTS chip_transactions (
  id TEXT PRIMARY KEY,
  player_id TEXT NOT NULL REFERENCES players(id),
  amount BIGINT NOT NULL,
  reason TEXT NOT NULL,
  round_id TEXT NULL REFERENCES game_rounds(id),
  created_at DATETIME NOT NULL,
  UNIQUE(round_id, player_id, reason)
);
CREATE INDEX IF NOT EXISTS idx_bets_round_id ON bets(round_id);
CREATE INDEX IF NOT EXISTS idx_bets_player_id ON bets(player_id);
CREATE INDEX IF NOT EXISTS idx_game_rounds_settled_at ON game_rounds(settled_at DESC);
CREATE INDEX IF NOT EXISTS idx_chip_tx_player_created ON chip_transactions(player_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_chip_tx_round_id ON chip_transactions(round_id);
`,
	2: `
ALTER TABLE players ADD COLUMN fingerprint TEXT;
CREATE UNIQUE INDEX IF NOT EXISTS idx_players_fingerprint ON players(fingerprint);
`,
}
