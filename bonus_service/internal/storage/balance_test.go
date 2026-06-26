package storage_test

import (
	"context"
	"testing"

	"bonus-service/internal/storage"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBalanceRepo_GetBalance_NotFound(t *testing.T) {
	db := newTestDB(t)
	repo := storage.NewBalanceRepo(db)

	bal, err := repo.GetBalance(context.Background(), uuid.New())
	require.NoError(t, err)
	assert.Equal(t, int64(0), bal.Available)
	assert.Equal(t, int64(0), bal.Held)
}

func TestBalanceRepo_GetBalance_Existing(t *testing.T) {
	db := newTestDB(t)
	repo := storage.NewBalanceRepo(db)
	ctx := context.Background()
	userID := uuid.New()

	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	_, err = repo.GetBalanceForUpdate(ctx, tx, userID)
	require.NoError(t, err)
	require.NoError(t, repo.UpdateBalance(ctx, tx, userID, 500, 0))
	require.NoError(t, tx.Commit())

	bal, err := repo.GetBalance(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(500), bal.Available)
	assert.Equal(t, int64(0), bal.Held)
}

func TestBalanceRepo_UpdateBalance_Deltas(t *testing.T) {
	db := newTestDB(t)
	repo := storage.NewBalanceRepo(db)
	ctx := context.Background()
	userID := uuid.New()

	// create balance
	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	_, err = repo.GetBalanceForUpdate(ctx, tx, userID)
	require.NoError(t, err)
	require.NoError(t, repo.UpdateBalance(ctx, tx, userID, 1000, 0))
	require.NoError(t, tx.Commit())

	// move 300 from available to held
	tx, err = db.BeginTx(ctx, nil)
	require.NoError(t, err)
	require.NoError(t, repo.UpdateBalance(ctx, tx, userID, -300, 300))
	require.NoError(t, tx.Commit())

	bal, err := repo.GetBalance(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(700), bal.Available)
	assert.Equal(t, int64(300), bal.Held)
}

func TestBalanceRepo_UpdateBalance_NotUpdated(t *testing.T) {
	db := newTestDB(t)
	repo := storage.NewBalanceRepo(db)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()

	err = repo.UpdateBalance(ctx, tx, uuid.New(), 100, 0)
	assert.ErrorIs(t, err, storage.ErrBalanceNotUpdated)
}

func TestBalanceRepo_GetBalanceForUpdate_CreatesIfMissing(t *testing.T) {
	db := newTestDB(t)
	repo := storage.NewBalanceRepo(db)
	ctx := context.Background()
	userID := uuid.New()

	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()

	bal, err := repo.GetBalanceForUpdate(ctx, tx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), bal.Available)
	assert.Equal(t, userID, bal.UserID)
}
