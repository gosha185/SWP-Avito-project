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
