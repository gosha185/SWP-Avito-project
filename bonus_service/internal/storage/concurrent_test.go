package storage_test

import (
	"context"
	"sync"
	"testing"

	"bonus-service/internal/models"
	"bonus-service/internal/service"
	"bonus-service/internal/storage"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestConcurrent_Accrual_BalanceConsistency(t *testing.T) {
	svc := service.NewBonusService(newTestDB(t))
	ctx := context.Background()
	userID := uuid.New()

	const goroutines = 10
	const amount = 100

	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			entry := &models.LedgerEntry{
				UserID:        userID,
				OperationType: models.OpAccrual,
				Amount:        amount,
				ExternalKey:   uuid.NewString(),
			}
			_, err := svc.Accrue(ctx, entry, 30)
			assert.NoError(t, err)
		}()
	}
	wg.Wait()

	available, err := svc.GetAvailablePoints(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(goroutines*amount), available)
}

func TestConcurrent_Hold_OnlyOneSucceedsWhenBalanceExact(t *testing.T) {
	svc := service.NewBonusService(newTestDB(t))
	ctx := context.Background()
	userID := uuid.New()

	_, err := svc.Accrue(ctx, &models.LedgerEntry{
		UserID:        userID,
		OperationType: models.OpAccrual,
		Amount:        100,
		ExternalKey:   uuid.NewString(),
	}, 30)
	require.NoError(t, err)

	const goroutines = 5
	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		successes int
		errs      []error
	)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			entry := &models.LedgerEntry{
				UserID:        userID,
				OperationType: models.OpHold,
				Amount:        100,
				ExternalKey:   uuid.NewString(),
			}
			_, holdErr := svc.Hold(ctx, entry, uuid.New(), 24)
			mu.Lock()
			if holdErr == nil {
				successes++
			} else {
				errs = append(errs, holdErr)
			}
			mu.Unlock()
		}()
	}
	wg.Wait()

	assert.Equal(t, 1, successes, "ровно один холд должен пройти")
	assert.Equal(t, goroutines-1, len(errs))
	for _, e := range errs {
		assert.ErrorIs(t, e, storage.ErrInsufficientBalance)
	}

	held, err := svc.GetHeld(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(100), held)
}

func TestConcurrent_ConfirmAndCancel_SameHold(t *testing.T) {
	svc := service.NewBonusService(newTestDB(t))
	ctx := context.Background()
	userID := uuid.New()
	orderID := uuid.New()

	_, err := svc.Accrue(ctx, &models.LedgerEntry{
		UserID:        userID,
		OperationType: models.OpAccrual,
		Amount:        500,
		ExternalKey:   uuid.NewString(),
	}, 30)
	require.NoError(t, err)

	_, err = svc.Hold(ctx, &models.LedgerEntry{
		UserID:        userID,
		OperationType: models.OpHold,
		Amount:        500,
		ExternalKey:   uuid.NewString(),
	}, orderID, 24)
	require.NoError(t, err)

	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results []error
	)

	attempt := func(fn func() (int64, error)) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, e := fn()
			mu.Lock()
			results = append(results, e)
			mu.Unlock()
		}()
	}

	attempt(func() (int64, error) {
		return svc.ConfirmWithdraw(ctx, &models.LedgerEntry{
			UserID:        userID,
			OperationType: models.OpConfirm,
			ExternalKey:   uuid.NewString(),
		}, orderID)
	})
	attempt(func() (int64, error) {
		return svc.CancelHold(ctx, &models.LedgerEntry{
			UserID:        userID,
			OperationType: models.OpCancel,
			ExternalKey:   uuid.NewString(),
		}, orderID)
	})

	wg.Wait()

	successes := 0
	failures := 0
	for _, e := range results {
		if e == nil {
			successes++
		} else {
			failures++
			assert.ErrorIs(t, e, storage.ErrHoldNotFound)
		}
	}
	assert.Equal(t, 1, successes, "ровно одна операция должна победить")
	assert.Equal(t, 1, failures)
}

func TestConcurrent_Accrual_And_Hold_RaceBalance(t *testing.T) {
	db := newTestDB(t)
	svc := service.NewBonusService(db)
	ctx := context.Background()
	userID := uuid.New()

	const seed = int64(1000)
	_, err := svc.Accrue(ctx, &models.LedgerEntry{
		UserID:        userID,
		OperationType: models.OpAccrual,
		Amount:        seed,
		ExternalKey:   uuid.NewString(),
	}, 30)
	require.NoError(t, err)

	const accruals = 5
	const holdAmount = int64(50)

	var wg sync.WaitGroup

	for i := 0; i < accruals; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = svc.Accrue(ctx, &models.LedgerEntry{
				UserID:        userID,
				OperationType: models.OpAccrual,
				Amount:        100,
				ExternalKey:   uuid.NewString(),
			}, 30)
		}()
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = svc.Hold(ctx, &models.LedgerEntry{
				UserID:        userID,
				OperationType: models.OpHold,
				Amount:        holdAmount,
				ExternalKey:   uuid.NewString(),
			}, uuid.New(), 24)
		}()
	}

	wg.Wait()

	balRepo := storage.NewBalanceRepo(db)
	bal, err := balRepo.GetBalance(ctx, userID)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, bal.Available, int64(0))
	assert.GreaterOrEqual(t, bal.Held, int64(0))
}
