package service

import (
	"bonus-service/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (bs *BonusService) Hold(ctx context.Context, entry *models.LedgerEntry, order uuid.UUID, hours int64) (int64, error) {
	tx, _, err := bs.validate(ctx, entry.UserID, entry.Amount, entry.ExternalKey, entry.Amount, hours) //Validation can not be omitted
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	hold := &models.Hold{0, entry.UserID, order, entry.Amount, "active", time.Now().Add(time.Hour * time.Duration(hours)), time.Now()}
	batches, err := bs.BatchesDB.GetExpiringBatchesForUpdate(ctx, tx, entry.UserID)
	if err != nil {
		return 0, err
	}
	err = bs.HoldsDB.CreateHold(ctx, tx, hold)
	if err != nil {
		return 0, err
	}
	rest := entry.Amount
	for _, batch := range batches {
		d := min(batch.Remaining, rest)
		rest -= d
		err = bs.BatchesDB.DecreaseBatchRemaining(ctx, tx, batch.ID, d) //One of the main features of the method
		if err != nil {
			return 0, err
		}
		err = bs.HoldBatchesDB.CreateHoldBatch(ctx, tx, hold.ID, entry.UserID, batch.ID, d) //One of the main features of the method
		if err != nil {
			return 0, err
		}
		if rest == 0 {
			break
		}
	}
	err = bs.BalancesDB.UpdateBalance(ctx, tx, entry.UserID, -entry.Amount, entry.Amount) //One of the main features of the method
	if err != nil {
		return 0, err
	}
	entry.Metadata = json.RawMessage(fmt.Sprintf(`{"order_id": "%s", "expires_at": "%s"}`, order.String(), hold.ExpiresAt.Format(time.RFC3339)))
	err = bs.LedgersDB.Insert(ctx, tx, entry) //Maintaining ledger can not be omitted
	if err != nil {
		return 0, err
	}
	return hold.ID, tx.Commit()
}
