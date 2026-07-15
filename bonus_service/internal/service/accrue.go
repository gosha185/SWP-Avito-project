package service

import (
	"bonus-service/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

/*
Accrue takes context; reference to ledger entry containing user id, idempotency
key, and amount; amount of days for batch's lifetime. It validates request;
fills in the rest of information in entry; updates user balance; creates bonus
batch; inserts ledger entry. It returns id of the new batch (0 if transaction is
failed) and error (nil if transaction is successful).
*/
func (bs *BonusService) Accrue(ctx context.Context, entry *models.LedgerEntry, days int64) (int64, error) {
	tx, _, err := bs.validate(ctx, entry.UserID, 0, entry.ExternalKey, entry.Amount, days) //Validation can not be omitted
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	batch := models.BonusBatch{0, entry.UserID, entry.Amount, entry.Amount, time.Now().Add(time.Hour * 24 * time.Duration(days)), time.Now()}

	err = bs.BatchesDB.CreateBatch(ctx, tx, &batch) //One of the main features of the method
	if err != nil {
		return 0, err
	}

	err = bs.BalancesDB.UpdateBalance(ctx, tx, entry.UserID, entry.Amount, 0) //One of the main features of the method
	if err != nil {
		return 0, err
	}

	entry.Metadata = json.RawMessage(fmt.Sprintf(`{"ttl_days": "%d", "expires_at": "%s"}`, days, batch.ExpiresAt.Format(time.RFC3339)))
	err = bs.LedgersDB.Insert(ctx, tx, entry) //Maintaining ledger can not be omitted
	if err != nil {
		return 0, err
	}

	return batch.ID, tx.Commit()
}
