package storage

import (
	"bonus-service/internal/models"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type HoldBatchRepo struct {
	db *sql.DB
}

func NewHoldBatchRepo(db *sql.DB) *HoldBatchRepo {
	return &HoldBatchRepo{db: db}
}

// CreateHoldBatch links a hold to a batch with the amount taken from that batch.
// Used during Hold operation to track which batches were locked.
func (r *HoldBatchRepo) CreateHoldBatch(ctx context.Context, tx *sql.Tx, holdID int64, batchUserID uuid.UUID, batchID int64, amount int64) error {
	query := `
        INSERT INTO hold_batches (hold_id, batch_user_id, batch_id, amount)
        VALUES ($1, $2, $3, $4)
    `
	_, err := tx.ExecContext(ctx, query, holdID, batchUserID, batchID, amount)
	if err != nil {
		return fmt.Errorf("failed to link hold %d to batch %d (user %s): %w", holdID, batchID, batchUserID, err)
	}
	return nil
}

// GetHoldBatchesByHoldID returns all hold-batch links for a given hold.
// Used during Cancel to return points to the original batches.
func (r *HoldBatchRepo) GetHoldBatchesByHoldID(ctx context.Context, tx *sql.Tx, holdID int64) ([]models.HoldBatch, error) {
	query := `
		SELECT hold_id, batch_user_id, batch_id, amount
		FROM hold_batches
		WHERE hold_id = $1
	`
	rows, err := tx.QueryContext(ctx, query, holdID)
	if err != nil {
		return nil, fmt.Errorf("failed to get hold_batches for hold %d: %w", holdID, err)
	}
	defer rows.Close()

	var result []models.HoldBatch
	for rows.Next() {
		var hb models.HoldBatch
		if err := rows.Scan(&hb.HoldID, &hb.BatchUserID, &hb.BatchID, &hb.Amount); err != nil {
			return nil, fmt.Errorf("failed to scan hold_batch: %w", err)
		}
		result = append(result, hb)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return result, nil
}
