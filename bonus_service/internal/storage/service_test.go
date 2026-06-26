package storage_test

import (
	"context"
	"testing"

	"bonus-service/internal/models"
	"bonus-service/internal/service"
	"bonus-service/internal/storage"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newService(t *testing.T) *service.BonusService {
	t.Helper()
	return service.NewBonusService(newTestDB(t))
}

func accrualEntry(userID uuid.UUID, amount int64) *models.LedgerEntry {
	return &models.LedgerEntry{
		UserID:        userID,
		OperationType: models.OpAccrual,
		Amount:        amount,
		ExternalKey:   uuid.NewString(),
	}
}

func TestService_Accrue_IncreasesBalance(t *testing.T) {
	svc := newService(t)
	ctx := context.Background()
	userID := uuid.New()

	batchID, err := svc.Accrue(ctx, accrualEntry(userID, 500), 30)
	require.NoError(t, err)
	assert.NotZero(t, batchID)

	available, err := svc.GetAvailablePoints(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(500), available)
}

func TestService_Accrue_Idempotent(t *testing.T) {
	svc := newService(t)
	ctx := context.Background()
	userID := uuid.New()
	entry := accrualEntry(userID, 200)

	_, err := svc.Accrue(ctx, entry, 30)
	require.NoError(t, err)

	// same external key — second call must fail
	entry2 := accrualEntry(userID, 200)
	entry2.ExternalKey = entry.ExternalKey
	_, err = svc.Accrue(ctx, entry2, 30)
	assert.ErrorIs(t, err, storage.ErrLedgerDuplicate)

	// balance should remain 200
	available, err := svc.GetAvailablePoints(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(200), available)
}

func TestService_Hold_MovesAvailableToHeld(t *testing.T) {
	svc := newService(t)
	ctx := context.Background()
	userID := uuid.New()
	orderID := uuid.New()

	_, err := svc.Accrue(ctx, accrualEntry(userID, 1000), 30)
	require.NoError(t, err)

	holdEntry := &models.LedgerEntry{
		UserID:        userID,
		OperationType: models.OpHold,
		Amount:        300,
		ExternalKey:   uuid.NewString(),
	}
	holdID, err := svc.Hold(ctx, holdEntry, orderID, 24)
	require.NoError(t, err)
	assert.NotZero(t, holdID)

	available, err := svc.GetAvailablePoints(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(700), available)

	held, err := svc.GetHeld(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(300), held)
}

func TestService_Hold_InsufficientBalance(t *testing.T) {
	svc := newService(t)
	ctx := context.Background()
	userID := uuid.New()

	_, err := svc.Accrue(ctx, accrualEntry(userID, 100), 30)
	require.NoError(t, err)

	holdEntry := &models.LedgerEntry{
		UserID:        userID,
		OperationType: models.OpHold,
		Amount:        500,
		ExternalKey:   uuid.NewString(),
	}
	_, err = svc.Hold(ctx, holdEntry, uuid.New(), 24)
	assert.ErrorIs(t, err, storage.ErrInsufficientBalance)
}

func TestService_ConfirmWithdraw(t *testing.T) {
	svc := newService(t)
	ctx := context.Background()
	userID := uuid.New()
	orderID := uuid.New()

	_, err := svc.Accrue(ctx, accrualEntry(userID, 1000), 30)
	require.NoError(t, err)

	holdEntry := &models.LedgerEntry{
		UserID:        userID,
		OperationType: models.OpHold,
		Amount:        400,
		ExternalKey:   uuid.NewString(),
	}
	_, err = svc.Hold(ctx, holdEntry, orderID, 24)
	require.NoError(t, err)

	confirmEntry := &models.LedgerEntry{
		UserID:        userID,
		OperationType: models.OpConfirm,
		ExternalKey:   uuid.NewString(),
	}
	holdID, err := svc.ConfirmWithdraw(ctx, confirmEntry, orderID)
	require.NoError(t, err)
	assert.NotZero(t, holdID)

	available, err := svc.GetAvailablePoints(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(600), available)

	held, err := svc.GetHeld(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), held)
}

func TestService_CancelHold_ReturnsPointsToBatches(t *testing.T) {
	svc := newService(t)
	ctx := context.Background()
	userID := uuid.New()
	orderID := uuid.New()

	_, err := svc.Accrue(ctx, accrualEntry(userID, 1000), 30)
	require.NoError(t, err)

	holdEntry := &models.LedgerEntry{
		UserID:        userID,
		OperationType: models.OpHold,
		Amount:        300,
		ExternalKey:   uuid.NewString(),
	}
	_, err = svc.Hold(ctx, holdEntry, orderID, 24)
	require.NoError(t, err)

	cancelEntry := &models.LedgerEntry{
		UserID:        userID,
		OperationType: models.OpCancel,
		ExternalKey:   uuid.NewString(),
	}
	holdID, err := svc.CancelHold(ctx, cancelEntry, orderID)
	require.NoError(t, err)
	assert.NotZero(t, holdID)

	available, err := svc.GetAvailablePoints(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(1000), available)

	held, err := svc.GetHeld(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), held)
}

func TestService_Hold_FEFO_SpansMultipleBatches(t *testing.T) {
	svc := newService(t)
	ctx := context.Background()
	userID := uuid.New()
	orderID := uuid.New()

	// two batches: 200 expiring in 10 days, 300 expiring in 30 days
	_, err := svc.Accrue(ctx, accrualEntry(userID, 200), 10)
	require.NoError(t, err)
	_, err = svc.Accrue(ctx, accrualEntry(userID, 300), 30)
	require.NoError(t, err)

	holdEntry := &models.LedgerEntry{
		UserID:        userID,
		OperationType: models.OpHold,
		Amount:        350,
		ExternalKey:   uuid.NewString(),
	}
	_, err = svc.Hold(ctx, holdEntry, orderID, 24)
	require.NoError(t, err)

	available, err := svc.GetAvailablePoints(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(150), available)
}

func TestService_GetHistory(t *testing.T) {
	svc := newService(t)
	ctx := context.Background()
	userID := uuid.New()

	for i := 0; i < 3; i++ {
		_, err := svc.Accrue(ctx, accrualEntry(userID, 100), 30)
		require.NoError(t, err)
	}

	history, err := svc.GetHistory(ctx, userID, 10, 0)
	require.NoError(t, err)
	assert.Len(t, history, 3)
	for _, e := range history {
		assert.Equal(t, models.OpAccrual, e.OperationType)
	}
}
