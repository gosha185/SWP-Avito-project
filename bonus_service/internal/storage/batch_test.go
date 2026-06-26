package storage_test

import (
	"context"
	"testing"
	"time"

	"bonus-service/internal/models"
	"bonus-service/internal/storage"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newBatch(userID uuid.UUID, amount int64, expiresIn time.Duration) *models.BonusBatch {
	return &models.BonusBatch{
		UserID:    userID,
		Amount:    amount,
		Remaining: amount,
		ExpiresAt: time.Now().Add(expiresIn),
	}
}

func TestBatchRepo_CreateBatch(t *testing.T) {
	db := newTestDB(t)
	repo := storage.NewBatchRepo(db)
	ctx := context.Background()
	userID := uuid.New()

	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)

	batch := newBatch(userID, 100, 24*time.Hour)
	require.NoError(t, repo.CreateBatch(ctx, tx, batch))
	require.NoError(t, tx.Commit())

	assert.NotZero(t, batch.ID)
	assert.False(t, batch.CreatedAt.IsZero())
}

func TestBatchRepo_GetExpiringBatches_FEFOOrder(t *testing.T) {
	db := newTestDB(t)
	repo := storage.NewBatchRepo(db)
	ctx := context.Background()
	userID := uuid.New()

	// insert batches out of order: far expiry first, near expiry second
	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	b1 := newBatch(userID, 200, 48*time.Hour)
	b2 := newBatch(userID, 100, 12*time.Hour)
	require.NoError(t, repo.CreateBatch(ctx, tx, b1))
	require.NoError(t, repo.CreateBatch(ctx, tx, b2))
	require.NoError(t, tx.Commit())

	batches, err := repo.GetExpiringBatches(ctx, userID)
	require.NoError(t, err)
	require.Len(t, batches, 2)
	// FEFO: nearest expiry first
	assert.True(t, batches[0].ExpiresAt.Before(batches[1].ExpiresAt))
}

func TestBatchRepo_DecreaseBatchRemaining(t *testing.T) {
	db := newTestDB(t)
	repo := storage.NewBatchRepo(db)
	ctx := context.Background()
	userID := uuid.New()

	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	batch := newBatch(userID, 100, 24*time.Hour)
	require.NoError(t, repo.CreateBatch(ctx, tx, batch))
	require.NoError(t, tx.Commit())

	tx, err = db.BeginTx(ctx, nil)
	require.NoError(t, err)
	require.NoError(t, repo.DecreaseBatchRemaining(ctx, tx, batch.ID, 60))
	require.NoError(t, tx.Commit())

	batches, err := repo.GetExpiringBatches(ctx, userID)
	require.NoError(t, err)
	require.Len(t, batches, 1)
	assert.Equal(t, int64(40), batches[0].Remaining)
}

func TestBatchRepo_DecreaseBatchRemaining_Insufficient(t *testing.T) {
	db := newTestDB(t)
	repo := storage.NewBatchRepo(db)
	ctx := context.Background()
	userID := uuid.New()

	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	batch := newBatch(userID, 50, 24*time.Hour)
	require.NoError(t, repo.CreateBatch(ctx, tx, batch))
	require.NoError(t, tx.Commit())

	tx, err = db.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()
	err = repo.DecreaseBatchRemaining(ctx, tx, batch.ID, 100)
	assert.ErrorIs(t, err, storage.ErrBatchNotFound)
}

func TestBatchRepo_IncreaseBatchRemaining(t *testing.T) {
	db := newTestDB(t)
	repo := storage.NewBatchRepo(db)
	ctx := context.Background()
	userID := uuid.New()

	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	batch := newBatch(userID, 100, 24*time.Hour)
	require.NoError(t, repo.CreateBatch(ctx, tx, batch))
	require.NoError(t, repo.DecreaseBatchRemaining(ctx, tx, batch.ID, 100))
	require.NoError(t, tx.Commit())

	tx, err = db.BeginTx(ctx, nil)
	require.NoError(t, err)
	require.NoError(t, repo.IncreaseBatchRemaining(ctx, tx, batch.ID, 100))
	require.NoError(t, tx.Commit())

	batches, err := repo.GetExpiringBatches(ctx, userID)
	require.NoError(t, err)
	require.Len(t, batches, 1)
	assert.Equal(t, int64(100), batches[0].Remaining)
}

func TestBatchRepo_GetExpiringSum(t *testing.T) {
	db := newTestDB(t)
	repo := storage.NewBatchRepo(db)
	ctx := context.Background()
	userID := uuid.New()

	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	// expiring in 2 days
	require.NoError(t, repo.CreateBatch(ctx, tx, newBatch(userID, 100, 2*24*time.Hour)))
	// expiring in 10 days
	require.NoError(t, repo.CreateBatch(ctx, tx, newBatch(userID, 200, 10*24*time.Hour)))
	require.NoError(t, tx.Commit())

	sum, err := repo.GetExpiringSum(ctx, userID, 5)
	require.NoError(t, err)
	assert.Equal(t, int64(100), sum)
}

func TestBatchRepo_GetAllBatches_IncludesSpent(t *testing.T) {
	db := newTestDB(t)
	repo := storage.NewBatchRepo(db)
	ctx := context.Background()
	userID := uuid.New()

	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	batch := newBatch(userID, 50, 24*time.Hour)
	require.NoError(t, repo.CreateBatch(ctx, tx, batch))
	require.NoError(t, repo.DecreaseBatchRemaining(ctx, tx, batch.ID, 50))
	require.NoError(t, tx.Commit())

	// GetExpiringBatches should exclude spent
	expiring, err := repo.GetExpiringBatches(ctx, userID)
	require.NoError(t, err)
	assert.Empty(t, expiring)

	// GetAllBatches should include it
	all, err := repo.GetAllBatches(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, all, 1)
}
