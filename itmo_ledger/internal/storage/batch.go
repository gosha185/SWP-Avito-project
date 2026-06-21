package storage

import (
	"bonus-service/internal/models"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type BatchRepo struct {
	db *sql.DB
}

func NewBatchRepo(db *sql.DB) *BatchRepo {
	return &BatchRepo{db: db}
}

// scanBatches scans rows and returns a slice of BonusBatch.
// Used internally by all batch read methods.
func scanBatches(rows *sql.Rows) ([]models.BonusBatch, error) {
	defer rows.Close()

	var batches []models.BonusBatch
	for rows.Next() {
		var b models.BonusBatch
		err := rows.Scan(
			&b.ID,
			&b.UserID,
			&b.Amount,
			&b.Remaining,
			&b.ExpiresAt,
			&b.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan batch: %w", err)
		}
		batches = append(batches, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return batches, nil
}

// GetExpiringBatches returns all batches with remaining > 0, sorted by expires_at (FEFO).
// No locking. Used for read-only balance breakdowns.
func (r *BatchRepo) GetExpiringBatches(ctx context.Context, userID uuid.UUID) ([]models.BonusBatch, error) {
	query := `
        SELECT  id, user_id, amount, remaining, expires_at, created_at
        FROM batches
        WHERE user_id = $1 AND remaining > 0
        ORDER BY expires_at
    `

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query batches for user %s: %w", userID, err)
	}

	return scanBatches(rows)
}

// GetExpiringBatchesForUpdate returns batches with SELECT FOR UPDATE row lock.
// Sorted by expires_at (FEFO). Must be called within a transaction.
// Used by: POST /hold (locks batches before spending).
func (r *BatchRepo) GetExpiringBatchesForUpdate(ctx context.Context, tx *sql.Tx, userID uuid.UUID) ([]models.BonusBatch, error) {
	query := `
		SELECT id, user_id, amount, remaining, expires_at, created_at
		FROM batches
		WHERE user_id = $1 AND remaining > 0
		ORDER BY expires_at
		FOR UPDATE
	`

	rows, err := tx.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query batches with lock for user %s: %w", userID, err)
	}
	return scanBatches(rows)
}

// GetExpiringSum calculates total points expiring in the next N days.
// Used by: GET /balance?days=N (aggregated value).
func (r *BatchRepo) GetExpiringSum(ctx context.Context, userID uuid.UUID, days int) (int64, error) {
	query := `
		SELECT COALESCE(SUM(remaining), 0)
		FROM batches
		WHERE user_id = $1 
		  AND remaining > 0 
		  AND expires_at BETWEEN NOW() AND NOW() + ($2 || ' days')::INTERVAL
	`

	var sum int64
	err := r.db.QueryRowContext(ctx, query, userID, days).Scan(&sum)
	if err != nil {
		return 0, fmt.Errorf("failed to get expiring sum for user %s: %w", userID, err)
	}

	return sum, nil
}

// GetExpiringBreakdown returns detailed breakdown of expiring points by day.
// Used by: GET /balance?days=N&breakdown=true.
// Returns []models.ExpiryBreakdown: [{days_left: 3, amount: 100}, ...]
func (r *BatchRepo) GetExpiringBreakdown(ctx context.Context, userID uuid.UUID, days int) ([]models.ExpiryBreakdown, error) {
	query := `
		SELECT 
			EXTRACT(DAY FROM (expires_at - NOW()))::INT AS days_left,
			SUM(remaining) AS amount
		FROM batches
		WHERE user_id = $1 
		  AND remaining > 0 
		  AND expires_at BETWEEN NOW() AND NOW() + ($2 || ' days')::INTERVAL
		GROUP BY days_left
		ORDER BY days_left
	`

	rows, err := r.db.QueryContext(ctx, query, userID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get expiry breakdown for user %s: %w", userID, err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var breakdown []models.ExpiryBreakdown
	for rows.Next() {
		var b models.ExpiryBreakdown
		err := rows.Scan(&b.DaysLeft, &b.Amount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan breakdown: %w", err)
		}
		breakdown = append(breakdown, b)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return breakdown, nil
}

// CreateBatch creates a new batch of points.
// batch.ID and batch.CreatedAt are populated by the database via RETURNING.
// Must be called within a transaction.
func (r *BatchRepo) CreateBatch(ctx context.Context, tx *sql.Tx, batch *models.BonusBatch) error {
	query := `
		INSERT INTO batches (user_id, amount, remaining, expires_at, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, created_at
	`

	err := tx.QueryRowContext(ctx, query,
		batch.UserID,
		batch.Amount,
		batch.Remaining,
		batch.ExpiresAt,
	).Scan(&batch.ID, &batch.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create batch for user %s: %w", batch.UserID, err)
	}

	return nil
}

func (r *BatchRepo) IncreaseBatchRemaining(ctx context.Context, tx *sql.Tx, batchID int64, amount int64) error {
	query := `
        UPDATE batches
        SET remaining = remaining + $1
        WHERE id = $2
    `

	result, err := tx.ExecContext(ctx, query, amount, batchID)
	if err != nil {
		return fmt.Errorf("failed to increase batch %d: %w", batchID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrBatchNotFound
	}

	return nil
}

// DecreaseBatchRemaining subtracts amount from batch remaining.
// Returns ErrBatchNotFound if batch doesn't exist or remaining < amount.
// Must be called within a transaction (batch should already be locked).
func (r *BatchRepo) DecreaseBatchRemaining(ctx context.Context, tx *sql.Tx, batchID int64, amount int64) error {
	query := `
		UPDATE batches
		SET remaining = remaining - $1
		WHERE id = $2 AND remaining >= $1
	`

	result, err := tx.ExecContext(ctx, query, amount, batchID)
	if err != nil {
		return fmt.Errorf("failed to decrease batch %d: %w", batchID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrBatchNotFound
	}

	return nil
}

// GetAllBatches returns all batches for a user, including those with remaining = 0.
// No locking. Used for admin/debug purposes.
func (r *BatchRepo) GetAllBatches(ctx context.Context, userID uuid.UUID) ([]models.BonusBatch, error) {
	query := `
        SELECT id, user_id, amount, remaining, expires_at, created_at
        FROM batches
        WHERE user_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query all batches for user %s: %w", userID, err)
	}
	return scanBatches(rows)
}

// SetBatchRemaining sets the remaining amount for a batch directly.
// Use with caution — this bypasses the usual remaining >= amount check.
// Must be called within a transaction.
// Returns ErrBatchNotFound if batch doesn't exist.
func (r *BatchRepo) SetBatchRemaining(ctx context.Context, tx *sql.Tx, batchID int64, newRemaining int64) error {
	query := `
        UPDATE batches
        SET remaining = $1
        WHERE id = $2
    `

	result, err := tx.ExecContext(ctx, query, newRemaining, batchID)
	if err != nil {
		return fmt.Errorf("failed to set remaining for batch %d: %w", batchID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrBatchNotFound
	}

	return nil
}
