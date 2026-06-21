package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"simple-ledger.itmo.ru/internal/models"
)

func (bs *BonusService) Accrue(ctx context.Context, entry *models.LedgerEntry, days int64) (int64, error) {
	tx, _, err := bs.validate(ctx, entry.UserID, 0, entry.ExternalKey, entry.Amount, days)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	batch := models.BonusBatch{0, entry.UserID, entry.Amount, entry.Amount, time.Now().Add(time.Hour * 24 * time.Duration(days)), time.Now()}
	err = bs.BatchesDB.CreateBatch(ctx, tx, &batch)
	if err != nil {
		return 0, err
	}
	err = bs.BalancesDB.UpdateBalance(ctx, tx, entry.UserID, entry.Amount, 0)
	if err != nil {
		return 0, err
	}
	entry.Metadata = json.RawMessage(fmt.Sprintf(`{"ttl_days": "%d", "expires_at": "%s"}`, days, batch.ExpiresAt.Format(time.RFC3339)))
	err = bs.LedgersDB.Insert(ctx, tx, entry) //Maintaining ledger
	if err != nil {
		return 0, err
	}
	return batch.ID, tx.Commit()
}
