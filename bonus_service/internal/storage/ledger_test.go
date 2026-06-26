package storage_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"bonus-service/internal/models"
	"bonus-service/internal/storage"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeEntry(userID uuid.UUID, key string, amount int64) *models.LedgerEntry {
	return &models.LedgerEntry{
		UserID:        userID,
		OperationType: models.OpAccrual,
		Amount:        amount,
		ExternalKey:   key,
		Metadata:      json.RawMessage(`{"ttl_days":"30"}`),
	}
}

func TestLedgerRepo_Insert(t *testing.T) {
	db := newTestDB(t)
	repo := storage.NewLedgerRepo(db)
	ctx := context.Background()
	userID := uuid.New()

	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	entry := makeEntry(userID, uuid.NewString(), 500)
	require.NoError(t, repo.Insert(ctx, tx, entry))
	require.NoError(t, tx.Commit())

	assert.NotZero(t, entry.ID)
	assert.False(t, entry.CreatedAt.IsZero())
}

func TestLedgerRepo_Insert_DuplicateKey(t *testing.T) {
	db := newTestDB(t)
	repo := storage.NewLedgerRepo(db)
	ctx := context.Background()
	userID := uuid.New()
	key := uuid.NewString()

	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	require.NoError(t, repo.Insert(ctx, tx, makeEntry(userID, key, 100)))
	require.NoError(t, tx.Commit())

	tx, err = db.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()
	err = repo.Insert(ctx, tx, makeEntry(userID, key, 200))
	assert.ErrorIs(t, err, storage.ErrLedgerDuplicate)
}

func TestLedgerRepo_GetByExternalKey(t *testing.T) {
	db := newTestDB(t)
	repo := storage.NewLedgerRepo(db)
	ctx := context.Background()
	userID := uuid.New()
	key := uuid.NewString()

	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	entry := makeEntry(userID, key, 300)
	require.NoError(t, repo.Insert(ctx, tx, entry))
	require.NoError(t, tx.Commit())

	found, err := repo.GetByExternalKey(ctx, key)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, int64(300), found.Amount)
	assert.Equal(t, key, found.ExternalKey)
}

func TestLedgerRepo_GetByExternalKey_NotFound(t *testing.T) {
	db := newTestDB(t)
	repo := storage.NewLedgerRepo(db)
	ctx := context.Background()

	found, err := repo.GetByExternalKey(ctx, "nonexistent-key")
	require.NoError(t, err)
	assert.Nil(t, found)
}

func TestLedgerRepo_GetHistory_Pagination(t *testing.T) {
	db := newTestDB(t)
	repo := storage.NewLedgerRepo(db)
	ctx := context.Background()
	userID := uuid.New()

	// insert 5 entries
	for i := 0; i < 5; i++ {
		tx, err := db.BeginTx(ctx, nil)
		require.NoError(t, err)
		entry := &models.LedgerEntry{
			UserID:        userID,
			OperationType: models.OpAccrual,
			Amount:        int64(100 * (i + 1)),
			ExternalKey:   uuid.NewString(),
			BatchID:       sql.NullInt64{},
		}
		require.NoError(t, repo.Insert(ctx, tx, entry))
		require.NoError(t, tx.Commit())
	}

	page1, err := repo.GetHistory(ctx, userID, 3, 0)
	require.NoError(t, err)
	assert.Len(t, page1, 3)

	page2, err := repo.GetHistory(ctx, userID, 3, 3)
	require.NoError(t, err)
	assert.Len(t, page2, 2)

	count, err := repo.GetHistoryCount(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, 5, count)
}

func TestLedgerRepo_GetHistory_Empty(t *testing.T) {
	db := newTestDB(t)
	repo := storage.NewLedgerRepo(db)
	ctx := context.Background()

	entries, err := repo.GetHistory(ctx, uuid.New(), 10, 0)
	require.NoError(t, err)
	assert.Empty(t, entries)
}
