package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"simple-ledger.itmo.ru/internal/models"
)

func (bs *BonusService) Hold(ctx context.Context, entry *models.LedgerEntry, order uuid.UUID, hours int64) (int64, error) {
	tx, _, err := bs.validate(ctx, entry.UserID, entry.Amount, entry.ExternalKey, entry.Amount, hours)
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
	hold, err = bs.HoldsDB.GetHoldByOrderID(ctx, tx, entry.UserID, order)
	if err != nil {
		return 0, err
	}
	rest := entry.Amount
	for _, batch := range batches {
		d := min(batch.Remaining, rest)
		rest -= d
		err = bs.BatchesDB.DecreaseBatchRemaining(ctx, tx, batch.ID, d)
		if err != nil {
			return 0, err
		}
		err = bs.HoldBatchesDB.CreateHoldBatch(ctx, tx, hold.ID, entry.UserID, batch.ID, d)
		if err != nil {
			return 0, err
		}
		if rest == 0 {
			break
		}
	}
	err = bs.BalancesDB.UpdateBalance(ctx, tx, entry.UserID, -entry.Amount, entry.Amount)
	if err != nil {
		return 0, err
	}
	entry.Metadata = json.RawMessage(fmt.Sprintf(`{"order_id": "%s", "expires_at": "%s"}`, order.String(), hold.ExpiresAt.Format(time.RFC3339)))
	err = bs.LedgersDB.Insert(ctx, tx, entry) //Maintaining ledger
	if err != nil {
		return 0, err
	}
	return hold.ID, tx.Commit()
}
