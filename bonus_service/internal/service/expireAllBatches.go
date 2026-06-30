package service

import (
	"bonus-service/internal/models"
	"context"
	"encoding/json"
	"time"
)

func (bs *BonusService) ExpireAllBatches(ctx context.Context) error {
	tx, err := bs.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	entry := &models.LedgerEntry{OperationType: models.OpExpiry, CreatedAt: time.Now(), Metadata: json.RawMessage(`"expired": "batches"`)}
	batches, err := bs.BatchesDB.GetExpiredBatches(ctx, tx)
	if err != nil {
		return err
	}
	for _, batch := range batches {
		_, err := bs.BalancesDB.GetBalanceForUpdate(ctx, tx, batch.UserID)
		if err != nil {
			return err
		}
		err = bs.BalancesDB.UpdateBalance(ctx, tx, batch.UserID, -batch.Amount, 0)
		if err != nil {
			return err
		}
		err = bs.BatchesDB.DecreaseBatchRemaining(ctx, tx, batch.ID, batch.Amount)
		if err != nil {
			return err
		}
		entry.Amount += batch.Amount
	}
	_, err = bs.BatchesDB.DeleteExpiredZeroBatches(ctx, 0)
	if err != nil {
		return err
	}
	err = bs.LedgersDB.Insert(ctx, tx, entry)
	if err != nil {
		return err
	}
	return tx.Commit()
}
