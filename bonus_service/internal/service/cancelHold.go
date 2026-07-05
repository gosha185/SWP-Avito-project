package service

import (
	"bonus-service/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (bs *BonusService) CancelHold(ctx context.Context, entry *models.LedgerEntry, order uuid.UUID) (int64, error) {
	tx, _, err := bs.validate(ctx, entry.UserID, 0, entry.ExternalKey)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	hold, err := bs.HoldsDB.GetHoldByOrderID(ctx, tx, entry.UserID, order)
	if err != nil {
		return 0, err
	}
	holdBatches, err := bs.HoldBatchesDB.GetHoldBatchesByHoldID(ctx, tx, hold.ID)
	if err != nil {
		return 0, err
	}
	var sum int64
	for _, holdBatch := range holdBatches {
		var validness bool
		validness, err = bs.BatchesDB.IncreaseBatchRemaining(ctx, tx, holdBatch.BatchID, holdBatch.Amount)
		if err != nil {
			return 0, err
		}
		if validness {
			sum += holdBatch.Amount
		}
	}
	err = bs.HoldsDB.UpdateHoldStatus(ctx, tx, hold.ID, "cancelled")
	if err != nil {
		return 0, err
	}
	err = bs.BalancesDB.UpdateBalance(ctx, tx, entry.UserID, sum, -hold.Amount)
	if err != nil {
		return 0, err
	}
	entry.Metadata = json.RawMessage(fmt.Sprintf(`{"order_id": "%s", "expires_at": "%s"}`, order.String(), hold.ExpiresAt.Format(time.RFC3339)))
	entry.Amount = hold.Amount
	err = bs.LedgersDB.Insert(ctx, tx, entry)
	if err != nil {
		return 0, err
	}
	return hold.ID, tx.Commit()
}
