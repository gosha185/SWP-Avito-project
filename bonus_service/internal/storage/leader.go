package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type LeaderRepo struct {
	db *sql.DB
}

func NewLeaderRepo(db *sql.DB) *LeaderRepo {
	return &LeaderRepo{db: db}
}

// TryAcquireLock attempts to become the leader for the given role.
// Returns true if this instance is now the leader.
func (r *LeaderRepo) TryAcquireLock(ctx context.Context, roleName, instanceID string, ttlSeconds int) (bool, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	var leaderID string
	query := `
			INSERT INTO leaders (role_name, leader_id, updated_at) 
			VALUES ($1, $2, NOW())
			ON CONFLICT (role_name) DO UPDATE 
			SET leader_id = $2, 
			    updated_at = NOW()
			WHERE leaders.role_name = $1
				AND (leaders.updated_at < NOW() - ($3 || ' seconds')::INTERVAL
			    	OR leaders.leader_id = $2)
	`

	err = tx.QueryRowContext(ctx, query, roleName, instanceID, ttlSeconds).Scan(&leaderID)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to acquire lock: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("failed to commit: %w", err)
	}

	return leaderID == instanceID, nil
}

// RenewLock updates the updated_at timestamp for the current leader.
// Returns true if the lock was successfully renewed.
func (r *LeaderRepo) RenewLock(ctx context.Context, roleName, instanceID string) (bool, error) {
	query := `
			UPDATE leaders 
			SET updated_at = NOW()
			WHERE leaders.role_name = $1 AND leader_id = $2
			RETURNING leader_id
			`

	var leaderID string
	err := r.db.QueryRowContext(ctx, query, roleName, instanceID).Scan(&leaderID)
	if errors.Is(err, sql.ErrNoRows) {
		return false, fmt.Errorf("leader record not found for role %s and instance %s", roleName, instanceID)
	}
	if err != nil {
		return false, fmt.Errorf("failed to renew lock: %w", err)
	}

	return true, nil
}
