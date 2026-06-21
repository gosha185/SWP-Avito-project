package storage

import (
	"bonus-service/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type LedgerRepo struct {
	db *sql.DB
}

func NewLedgerRepo(db *sql.DB) *LedgerRepo {
	return &LedgerRepo{db: db}
}

// Insert appends an immutable audit record to the ledger.
// Returns ErrLedgerDuplicate if external_key already exists.
// Must be called within a transaction.
func (r *LedgerRepo) Insert(ctx context.Context, tx *sql.Tx, entry *models.LedgerEntry) error {
	query := `
		INSERT INTO ledger (user_id, operation_type, amount, batch_id, external_key, created_at, metadata)
		VALUES ($1, $2, $3, $4, $5, NOW(), $6)
		RETURNING id, created_at
	`

	var metadata any = nil
	if entry.Metadata != nil {
		metadata = entry.Metadata
	}

	var batchID any = nil
	if entry.BatchID.Valid {
		batchID = entry.BatchID.Int64
	}

	err := tx.QueryRowContext(ctx, query,
		entry.UserID,
		entry.OperationType,
		entry.Amount,
		batchID,
		entry.ExternalKey,
		metadata,
	).Scan(&entry.ID, &entry.CreatedAt)

	if err != nil {
		if isDuplicateKeyError(err) {
			return ErrLedgerDuplicate
		}
		return fmt.Errorf("failed to insert ledger entry for user %s, key %s: %w", entry.UserID, entry.ExternalKey, err)
	}

	return nil
}

// GetByExternalKey checks if an operation with this external_key already exists.
// Used for idempotency checks. Returns nil if not found.
func (r *LedgerRepo) GetByExternalKey(ctx context.Context, externalKey string) (*models.LedgerEntry, error) {
	query := `
		SELECT id, user_id, operation_type, amount, batch_id, external_key, created_at, metadata
		FROM ledger
		WHERE external_key = $1
	`

	var entry models.LedgerEntry
	var batchID sql.NullInt64
	var metadata []byte

	err := r.db.QueryRowContext(ctx, query, externalKey).Scan(
		&entry.ID,
		&entry.UserID,
		&entry.OperationType,
		&entry.Amount,
		&batchID,
		&entry.ExternalKey,
		&entry.CreatedAt,
		&metadata,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get ledger entry by key %s: %w", externalKey, err)
	}

	entry.BatchID = batchID
	if metadata != nil {
		entry.Metadata = metadata
	}

	return &entry, nil
}

// GetHistory returns paginated transaction history for a user.
// Used by: GET /balance/:user_id/history.
func (r *LedgerRepo) GetHistory(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.LedgerEntry, error) {
	query := `
		SELECT id, user_id, operation_type, amount, batch_id, external_key, created_at, metadata
		FROM ledger
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query ledger history for user %s: %w", userID, err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var entries []models.LedgerEntry
	for rows.Next() {
		var e models.LedgerEntry
		var batchID sql.NullInt64
		var metadata []byte

		err := rows.Scan(
			&e.ID,
			&e.UserID,
			&e.OperationType,
			&e.Amount,
			&batchID,
			&e.ExternalKey,
			&e.CreatedAt,
			&metadata,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ledger entry: %w", err)
		}

		e.BatchID = batchID
		if metadata != nil {
			e.Metadata = metadata
		}
		entries = append(entries, e)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return entries, nil
}

// GetHistoryCount returns total number of transactions for a user.
// Used for pagination.
func (r *LedgerRepo) GetHistoryCount(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM ledger
		WHERE user_id = $1
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count ledger entries for user %s: %w", userID, err)
	}

	return count, nil
}

// GetLedgerByOrderID finds all operations related to a specific order.
// Used for admin/debug purposes.
func (r *LedgerRepo) GetLedgerByOrderID(ctx context.Context, orderID uuid.UUID) ([]models.LedgerEntry, error) {
	query := `
	SELECT id, user_id, operation_type, amount, batch_id, external_key, created_at, metadata
	FROM ledger
	WHERE metadata @> $1
	ORDER BY created_at DESC
    `
	searchJSON := fmt.Sprintf(`{"order_id":"%s"}`, orderID)

	rows, err := r.db.QueryContext(ctx, query, searchJSON)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.LedgerEntry
	for rows.Next() {
		var e models.LedgerEntry
		var batchID sql.NullInt64
		var metadata []byte

		err := rows.Scan(
			&e.ID,
			&e.UserID,
			&e.OperationType,
			&e.Amount,
			&batchID,
			&e.ExternalKey,
			&e.CreatedAt,
			&metadata,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ledger row: %w", err)
		}

		if batchID.Valid {
			e.BatchID = batchID
		}

		if metadata != nil {
			e.Metadata = metadata
		}

		entries = append(entries, e)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return entries, nil
}

// isDuplicateKeyError checks if the error is a PostgreSQL unique violation (23505).
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	if pgErr, ok := errors.AsType[*pq.Error](err); ok {
		return pgErr.Code == "23505"
	}
	return false
}
