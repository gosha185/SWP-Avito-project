package storage

import (
	"bonus-service/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type HoldRepo struct {
	db *sql.DB
}

func NewHoldRepo(db *sql.DB) *HoldRepo {
	return &HoldRepo{db: db}
}

// GetHoldByID returns a hold by its ID.
// No locking. Used for read-only operations.
func (r *HoldRepo) GetHoldByID(ctx context.Context, holdID int64) (*models.Hold, error) {
	query := `
		SELECT id, user_id, order_id, amount, status, expires_at, created_at
		FROM holds
		WHERE id = $1
	`

	var hold models.Hold
	err := r.db.QueryRowContext(ctx, query, holdID).Scan(
		&hold.ID,
		&hold.UserID,
		&hold.OrderID,
		&hold.Amount,
		&hold.Status,
		&hold.ExpiresAt,
		&hold.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrHoldNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get hold by id %d: %w", holdID, err)
	}

	return &hold, nil
}

// GetHoldByOrderID finds an active hold by order_id with row lock.
// Must be called within a transaction.
// Used by: POST /confirm, POST /cancel.
func (r *HoldRepo) GetHoldByOrderID(ctx context.Context, tx *sql.Tx, userID uuid.UUID, orderID uuid.UUID) (*models.Hold, error) {
	query := `
		SELECT id, user_id, order_id, amount, status, expires_at, created_at
		FROM holds
		WHERE user_id = $1 AND order_id = $2 AND status = 'active'
		FOR UPDATE
	`

	var hold models.Hold
	err := tx.QueryRowContext(ctx, query, userID, orderID).Scan(
		&hold.ID,
		&hold.UserID,
		&hold.OrderID,
		&hold.Amount,
		&hold.Status,
		&hold.ExpiresAt,
		&hold.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrHoldNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get hold by order_id %s for user %s: %w", orderID, userID, err)
	}

	return &hold, nil
}

// GetActiveHoldsByUser returns all active holds for a user.
// No locking. Used for GET /balance (show held points).
func (r *HoldRepo) GetActiveHoldsByUser(ctx context.Context, userID uuid.UUID) ([]models.Hold, error) {
	query := `
		SELECT id, user_id, order_id, amount, status, expires_at, created_at
		FROM holds
		WHERE user_id = $1 AND status = 'active'
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query active holds for user %s: %w", userID, err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var holds []models.Hold
	for rows.Next() {
		var h models.Hold
		err := rows.Scan(
			&h.ID,
			&h.UserID,
			&h.OrderID,
			&h.Amount,
			&h.Status,
			&h.ExpiresAt,
			&h.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan hold: %w", err)
		}
		holds = append(holds, h)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return holds, nil
}

// CreateHold creates a new hold record.
// hold.ID and hold.CreatedAt are populated by the database via RETURNING.
// Must be called within a transaction.
func (r *HoldRepo) CreateHold(ctx context.Context, tx *sql.Tx, hold *models.Hold) error {
	query := `
		INSERT INTO holds (user_id, order_id, amount, status, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, created_at
	`

	err := tx.QueryRowContext(ctx, query,
		hold.UserID,
		hold.OrderID,
		hold.Amount,
		hold.Status,
		hold.ExpiresAt,
	).Scan(&hold.ID, &hold.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create hold for user %s, order %s: %w", hold.UserID, hold.OrderID, err)
	}

	return nil
}

// UpdateHoldStatus changes hold status to confirmed or cancelled.
// Must be called within a transaction.
// Returns ErrHoldNotFound if hold not found or not active.
func (r *HoldRepo) UpdateHoldStatus(ctx context.Context, tx *sql.Tx, holdID int64, status string) error {
	query := `
		UPDATE holds
		SET status = $1
		WHERE id = $2 AND status = 'active'
	`

	result, err := tx.ExecContext(ctx, query, status, holdID)
	if err != nil {
		return fmt.Errorf("failed to update hold status for hold %d: %w", holdID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrHoldNotFound
	}

	return nil
}

// GetExpiredHolds returns all active holds with expired_at < NOW().
// Must be called within a transaction.
// Used by TTL worker.
func (r *HoldRepo) GetExpiredHolds(ctx context.Context, tx *sql.Tx) ([]models.Hold, error) {
	query := `
        SELECT id, user_id, order_id, amount, status, expires_at, created_at
        FROM holds
        WHERE status = 'active' AND expires_at < NOW()
        FOR UPDATE
    `

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query expired holds: %w", err)
	}
	defer rows.Close()

	var holds []models.Hold
	for rows.Next() {
		var h models.Hold
		err := rows.Scan(
			&h.ID,
			&h.UserID,
			&h.OrderID,
			&h.Amount,
			&h.Status,
			&h.ExpiresAt,
			&h.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan expired hold: %w", err)
		}
		holds = append(holds, h)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return holds, nil
}

// DeleteOldHolds deletes holds with status in ('confirmed', 'cancelled')
// and expires_at < NOW() - daysOld.
// Returns number of deleted rows.
// hold_batches is deleted automatically via ON DELETE CASCADE.
func (r *HoldRepo) DeleteOldHolds(ctx context.Context, daysOld int) (int64, error) {
	query := `
        DELETE FROM holds
        WHERE status IN ('confirmed', 'cancelled')
          AND expires_at < NOW() - ($1 || ' days')::INTERVAL
    `
	result, err := r.db.ExecContext(ctx, query, daysOld)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old holds: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}
	return rowsAffected, nil
}

// CancelAllExpiredHolds cancels all holds with status = 'active' and expires_at < NOW().
// This is a mass operation used by the TTL worker.
// Must be called within a transaction.
// Returns number of updated rows.
func (r *HoldRepo) CancelAllExpiredHolds(ctx context.Context, tx *sql.Tx) (int64, error) {
	query := `
        UPDATE holds
        SET status = 'cancelled'
        WHERE status = 'active'
          AND expires_at < NOW()
    `

	result, err := tx.ExecContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to cancel expired holds: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}
