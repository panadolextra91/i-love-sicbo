package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"cachon-casino/internal/engine"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository { return &Repository{db: db} }

func (r *Repository) DB() *sql.DB { return r.db }

func (r *Repository) RegisterOrLoadPlayer(ctx context.Context, fingerprint, name string, now time.Time) (Player, int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Player{}, 0, mapSQLError(err)
	}
	defer func() { _ = tx.Rollback() }()

	p, err := r.getPlayerByFingerprintTx(ctx, tx, fingerprint)
	if err != nil {
		if !errors.Is(err, ErrPlayerNotFound) {
			return Player{}, 0, err
		}
		id := uuid.NewString()
		p = Player{ID: id, Name: name, Fingerprint: fingerprint, Chips: 0, CreatedAt: now, UpdatedAt: now}
		if _, err := tx.ExecContext(ctx, `INSERT INTO players(id,name,fingerprint,chips,last_bonus_at,created_at,updated_at) VALUES(?,?,?,?,?,?,?)`, p.ID, p.Name, p.Fingerprint, p.Chips, nil, now, now); err != nil {
			return Player{}, 0, mapSQLError(err)
		}
	}

	bonus := int64(0)
	if eligibleDailyBonus(p.LastBonusAt, now) {
		bonus = 1000
		if _, err := tx.ExecContext(ctx, `UPDATE players SET chips = chips + ?, last_bonus_at = ?, updated_at = ? WHERE id = ?`, bonus, now, now, p.ID); err != nil {
			return Player{}, 0, mapSQLError(err)
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO chip_transactions(id, player_id, amount, reason, round_id, created_at) VALUES(?,?,?,?,?,?)`, uuid.NewString(), p.ID, bonus, "DAILY_BONUS", nil, now); err != nil {
			return Player{}, 0, mapSQLError(err)
		}
		p.Chips += bonus
		p.LastBonusAt = &now
	}

	if err := tx.Commit(); err != nil {
		return Player{}, 0, mapSQLError(err)
	}

	if name != "" && name != p.Name {
		_, _ = r.db.ExecContext(ctx, `UPDATE players SET name = ?, updated_at = ? WHERE id = ?`, name, now, p.ID)
		p.Name = name
	}

	return p, bonus, nil
}

func (r *Repository) getPlayerByFingerprintTx(ctx context.Context, tx *sql.Tx, fingerprint string) (Player, error) {
	var p Player
	var last sql.NullTime
	err := tx.QueryRowContext(ctx, `SELECT id,name,fingerprint,chips,last_bonus_at,created_at,updated_at FROM players WHERE fingerprint = ?`, fingerprint).Scan(&p.ID, &p.Name, &p.Fingerprint, &p.Chips, &last, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Player{}, ErrPlayerNotFound
		}
		return Player{}, mapSQLError(err)
	}
	if last.Valid {
		t := last.Time
		p.LastBonusAt = &t
	}
	return p, nil
}

func (r *Repository) SettleRound(ctx context.Context, roundID string, roundNo int64, startedAt, settledAt time.Time, dice engine.DiceResult, results []engine.PayoutResult) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return mapSQLError(err)
	}
	defer func() { _ = tx.Rollback() }()

	total := dice[0] + dice[1] + dice[2]
	if _, err := tx.ExecContext(ctx, `INSERT INTO game_rounds(id,round_no,dice_1,dice_2,dice_3,total,started_at,settled_at) VALUES(?,?,?,?,?,?,?,?)`, roundID, roundNo, dice[0], dice[1], dice[2], total, startedAt, settledAt); err != nil {
		if isUniqueErr(err) {
			return ErrRoundAlreadySettled
		}
		return mapSQLError(err)
	}

		for _, pr := range results {
		betID := uuid.NewString()
		if _, err := tx.ExecContext(ctx, `INSERT INTO bets(id,round_id,player_id,bet_type,amount,payout,is_win,created_at) VALUES(?,?,?,?,?,?,?,?)`, betID, roundID, pr.Bet.PlayerID, string(pr.Bet.Type), pr.Bet.Stake, pr.GrossPayout, pr.Win, settledAt); err != nil {
			if isUniqueErr(err) {
				return ErrRoundAlreadySettled
			}
			return mapSQLError(err)
		}
			netChange := pr.GrossPayout - pr.Bet.Stake
			reason := "ROUND_PAYOUT"
			if netChange < 0 {
				reason = "BET_DEBIT"
			}
			if _, err := tx.ExecContext(ctx, `INSERT INTO chip_transactions(id,player_id,amount,reason,round_id,created_at) VALUES(?,?,?,?,?,?)`, uuid.NewString(), pr.Bet.PlayerID, netChange, reason, roundID, settledAt); err != nil {
			if isUniqueErr(err) {
				return ErrRoundAlreadySettled
			}
			return mapSQLError(err)
		}
			if _, err := tx.ExecContext(ctx, `UPDATE players SET chips = chips + ?, updated_at = ? WHERE id = ?`, netChange, settledAt, pr.Bet.PlayerID); err != nil {
			return mapSQLError(err)
		}
	}

	if err := tx.Commit(); err != nil {
		return mapSQLError(err)
	}
	return nil
}

func (r *Repository) AuditChipLedger(ctx context.Context) error {
	var playersSum, ledgerSum sql.NullInt64
	if err := r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(chips),0) FROM players`).Scan(&playersSum); err != nil {
		return mapSQLError(err)
	}
	if err := r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(amount),0) FROM chip_transactions`).Scan(&ledgerSum); err != nil {
		return mapSQLError(err)
	}
	if playersSum.Int64 != ledgerSum.Int64 {
		return fmt.Errorf("audit mismatch players=%d ledger=%d", playersSum.Int64, ledgerSum.Int64)
	}
	return nil
}

func eligibleDailyBonus(last *time.Time, now time.Time) bool {
	if last == nil {
		return true
	}
	y1, m1, d1 := last.In(now.Location()).Date()
	y2, m2, d2 := now.Date()
	return y1 != y2 || m1 != m2 || d1 != d2
}

func mapSQLError(err error) error {
	if err == nil {
		return nil
	}
	var se sqlite3.Error
	if errors.As(err, &se) {
		if se.Code == sqlite3.ErrBusy || se.Code == sqlite3.ErrLocked {
			return ErrStorageBusy
		}
		if se.Code == sqlite3.ErrCorrupt || se.Code == sqlite3.ErrNotADB {
			return ErrStorageCorrupted
		}
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrPlayerNotFound
	}
	if strings.Contains(strings.ToLower(err.Error()), "insufficient") {
		return ErrInsufficientChips
	}
	return err
}

func isUniqueErr(err error) bool {
	var se sqlite3.Error
	if errors.As(err, &se) {
		return se.ExtendedCode == sqlite3.ErrConstraintUnique || se.ExtendedCode == sqlite3.ErrConstraintPrimaryKey
	}
	return false
}
