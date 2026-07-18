package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"bonus-service/internal/models"

	"github.com/google/uuid"
)

type BalanceRepo struct {
	db *sql.DB
	// TODO(post-mvp): Add shard resolver (by user_id) for horizontal scaling.
}

func NewBalanceRepo(db *sql.DB) *BalanceRepo {
	return &BalanceRepo{db: db}
}

// GetBalance reads current balance without locking.
// Creates and returns zero balance if user not found. No transaction required.
func (r *BalanceRepo) GetBalance(ctx context.Context, userID uuid.UUID) (*models.Balance, error) {
	query := `
		SELECT user_id, available, held, updated_at 
		FROM balances 
		WHERE user_id = $1
	`

	var balance models.Balance

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&balance.UserID,
		&balance.Available,
		&balance.Held,
		&balance.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return &models.Balance{
			UserID:    userID,
			Available: 0,
			Held:      0,
			UpdatedAt: time.Now(),
		}, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get balance for user %s: %w", userID, err)
	}

	return &balance, nil
}

// GetBalanceForUpdate reads balance with SELECT FOR UPDATE row lock.
// Creates balance if not exists. Must be called within a transaction.
// Used by: Accrue, Hold, Confirm, Cancel (all write operations).
func (r *BalanceRepo) GetBalanceForUpdate(ctx context.Context, tx *sql.Tx, userID uuid.UUID) (*models.Balance, error) {
	query := `
		SELECT user_id, available, held, updated_at 
		FROM balances 
		WHERE user_id = $1 
		FOR UPDATE
	`

	var balance models.Balance
	err := tx.QueryRowContext(ctx, query, userID).Scan(
		&balance.UserID,
		&balance.Available,
		&balance.Held,
		&balance.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		_, err := tx.ExecContext(
			ctx,
			`INSERT INTO balances (user_id, available, held, updated_at)
         VALUES ($1, 0, 0, NOW())
         ON CONFLICT (user_id) DO NOTHING`,
			userID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create balance for user %s: %w", userID, err)
		}

		err = tx.QueryRowContext(ctx, query, userID).Scan(
			&balance.UserID,
			&balance.Available,
			&balance.Held,
			&balance.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to reload balance for user %s: %w", userID, err)
		}

		return &balance, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get balance with lock for user %s: %w", userID, err)
	}

	return &balance, nil
}

// UpdateBalance applies delta changes to available and held.
// Must be called within a transaction after GetBalanceForUpdate.
// Deltas: positive = add, negative = subtract.
// Returns ErrBalanceNotUpdated if no rows affected.
func (r *BalanceRepo) UpdateBalance(ctx context.Context, tx *sql.Tx, userID uuid.UUID, deltaAvailable, deltaHeld int64) error {
	query := `
        UPDATE balances
        SET available = available + $1,
            held = held + $2,
            updated_at = NOW()
        WHERE user_id = $3
    `
	result, err := tx.ExecContext(ctx, query, deltaAvailable, deltaHeld, userID)
	if err != nil {
		return fmt.Errorf("failed to update balance for user %s: %w", userID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrBalanceNotUpdated
	}

	return nil
}

// UpdateBalancesForExpiredBatches decreases available balances
// for users whose batches have expired.
// Must be called within a transaction.
// Returns number of updated rows.
//
// Deprecated: now done atomically inside BatchRepo.ExpireAllBatches.
// Calling it separately reintroduces the race. Kept until the cleanup PR.
func (r *BalanceRepo) UpdateBalancesForExpiredBatches(ctx context.Context, tx *sql.Tx) (int64, error) {
	query := `
        UPDATE balances
        SET available = available - sub.expired_sum,
            updated_at = NOW()
        FROM (
            SELECT user_id, COALESCE(SUM(remaining), 0) AS expired_sum
            FROM batches
            WHERE remaining > 0 AND expires_at < NOW()
            GROUP BY user_id
        ) AS sub
        WHERE balances.user_id = sub.user_id
          AND sub.expired_sum > 0
    `

	result, err := tx.ExecContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to update balances for expired batches: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// UpdateBalancesForExpiredHolds restores available and decreases held
// for users whose holds have expired.
// Must be called within a transaction.
// Returns number of updated rows.
//
// Deprecated: now done atomically inside HoldRepo.CancelAllExpiredHolds.
// Calling it separately reintroduces the race. Kept until the cleanup PR.
func (r *BalanceRepo) UpdateBalancesForExpiredHolds(ctx context.Context, tx *sql.Tx) (int64, error) {
	query := `
        UPDATE balances
        SET available = available + sub.hold_sum,
            held = held - sub.hold_sum,
            updated_at = NOW()
        FROM (
            SELECT user_id, COALESCE(SUM(amount), 0) AS hold_sum
            FROM holds
            WHERE status = 'active' AND expires_at < NOW()
            GROUP BY user_id
        ) AS sub
        WHERE balances.user_id = sub.user_id
          AND sub.hold_sum > 0
    `

	result, err := tx.ExecContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to update balances for expired holds: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}
