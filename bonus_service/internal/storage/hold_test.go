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

func TestHoldRepo_CreateHold(t *testing.T) {
	db := newTestDB(t)
	batchRepo := storage.NewBatchRepo(db)
	holdRepo := storage.NewHoldRepo(db)
	ctx := context.Background()
	userID := uuid.New()
	orderID := uuid.New()

	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	batch := newBatch(userID, 100, 24*time.Hour)
	require.NoError(t, batchRepo.CreateBatch(ctx, tx, batch))
	require.NoError(t, tx.Commit())

	tx, err = db.BeginTx(ctx, nil)
	require.NoError(t, err)
	hold := &models.Hold{
		UserID:    userID,
		OrderID:   orderID,
		Amount:    100,
		Status:    "active",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	require.NoError(t, holdRepo.CreateHold(ctx, tx, hold))
	require.NoError(t, tx.Commit())

	assert.NotZero(t, hold.ID)
}

func TestHoldRepo_GetHoldByOrderID(t *testing.T) {
	db := newTestDB(t)
	batchRepo := storage.NewBatchRepo(db)
	holdRepo := storage.NewHoldRepo(db)
	ctx := context.Background()
	userID := uuid.New()
	orderID := uuid.New()

	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	batch := newBatch(userID, 100, 24*time.Hour)
	require.NoError(t, batchRepo.CreateBatch(ctx, tx, batch))
	hold := &models.Hold{
		UserID:    userID,
		OrderID:   orderID,
		Amount:    100,
		Status:    "active",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	require.NoError(t, holdRepo.CreateHold(ctx, tx, hold))
	require.NoError(t, tx.Commit())

	tx, err = db.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()

	found, err := holdRepo.GetHoldByOrderID(ctx, tx, userID, orderID)
	require.NoError(t, err)
	assert.Equal(t, hold.ID, found.ID)
	assert.Equal(t, int64(100), found.Amount)
	assert.Equal(t, "active", found.Status)
}

func TestHoldRepo_GetHoldByOrderID_NotFound(t *testing.T) {
	db := newTestDB(t)
	holdRepo := storage.NewHoldRepo(db)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()

	_, err = holdRepo.GetHoldByOrderID(ctx, tx, uuid.New(), uuid.New())
	assert.ErrorIs(t, err, storage.ErrHoldNotFound)
}

func TestHoldRepo_UpdateHoldStatus(t *testing.T) {
	db := newTestDB(t)
	batchRepo := storage.NewBatchRepo(db)
	holdRepo := storage.NewHoldRepo(db)
	ctx := context.Background()
	userID := uuid.New()
	orderID := uuid.New()

	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	batch := newBatch(userID, 100, 24*time.Hour)
	require.NoError(t, batchRepo.CreateBatch(ctx, tx, batch))
	hold := &models.Hold{
		UserID:    userID,
		OrderID:   orderID,
		Amount:    100,
		Status:    "active",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	require.NoError(t, holdRepo.CreateHold(ctx, tx, hold))
	require.NoError(t, tx.Commit())

	tx, err = db.BeginTx(ctx, nil)
	require.NoError(t, err)
	require.NoError(t, holdRepo.UpdateHoldStatus(ctx, tx, hold.ID, "confirmed"))
	require.NoError(t, tx.Commit())

	found, err := holdRepo.GetHoldByID(ctx, hold.ID)
	require.NoError(t, err)
	assert.Equal(t, "confirmed", found.Status)
}

func TestHoldRepo_UniqueOrderConstraint(t *testing.T) {
	db := newTestDB(t)
	batchRepo := storage.NewBatchRepo(db)
	holdRepo := storage.NewHoldRepo(db)
	ctx := context.Background()
	userID := uuid.New()
	orderID := uuid.New()

	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	batch := newBatch(userID, 200, 24*time.Hour)
	require.NoError(t, batchRepo.CreateBatch(ctx, tx, batch))
	hold := &models.Hold{
		UserID:    userID,
		OrderID:   orderID,
		Amount:    100,
		Status:    "active",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	require.NoError(t, holdRepo.CreateHold(ctx, tx, hold))
	require.NoError(t, tx.Commit())

	// second hold with same order should fail
	tx, err = db.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()
	hold2 := &models.Hold{
		UserID:    userID,
		OrderID:   orderID,
		Amount:    100,
		Status:    "active",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	err = holdRepo.CreateHold(ctx, tx, hold2)
	assert.Error(t, err)
}

func TestHoldBatchRepo_CreateAndGet(t *testing.T) {
	db := newTestDB(t)
	batchRepo := storage.NewBatchRepo(db)
	holdRepo := storage.NewHoldRepo(db)
	holdBatchRepo := storage.NewHoldBatchRepo(db)
	ctx := context.Background()
	userID := uuid.New()

	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	batch := newBatch(userID, 100, 24*time.Hour)
	require.NoError(t, batchRepo.CreateBatch(ctx, tx, batch))
	hold := &models.Hold{
		UserID:    userID,
		OrderID:   uuid.New(),
		Amount:    100,
		Status:    "active",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	require.NoError(t, holdRepo.CreateHold(ctx, tx, hold))
	require.NoError(t, holdBatchRepo.CreateHoldBatch(ctx, tx, hold.ID, userID, batch.ID, 100))
	require.NoError(t, tx.Commit())

	tx, err = db.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()

	hbs, err := holdBatchRepo.GetHoldBatchesByHoldID(ctx, tx, hold.ID)
	require.NoError(t, err)
	require.Len(t, hbs, 1)
	assert.Equal(t, batch.ID, hbs[0].BatchID)
	assert.Equal(t, int64(100), hbs[0].Amount)
}
