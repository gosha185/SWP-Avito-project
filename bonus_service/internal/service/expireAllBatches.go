package service

import (
	"bonus-service/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

func (bs *BonusService) ExpireAllBatches(ctx context.Context) error {
	batches, err := bs.BatchesDB.GetExpiredBatches(ctx)
	if err != nil {
		return err
	}
	var tx *sql.Tx
	for _, batch := range batches {
		tx, err = bs.DB.BeginTx(ctx, nil)
		if err != nil {
			break
		}
		entry := &models.LedgerEntry{OperationType: models.OpExpiry, Amount: batch.Remaining, BatchID: batch.ID, CreatedAt: time.Now(), Metadata: json.RawMessage(`"expired": "batches"`)}
		_, err := bs.BalancesDB.GetBalanceForUpdate(ctx, tx, batch.UserID)
		if err != nil {
			break
		}
		err = bs.BalancesDB.UpdateBalance(ctx, tx, batch.UserID, -batch.Remaining, 0)
		if err != nil {
			break
		}
		err = bs.BatchesDB.DecreaseBatchRemaining(ctx, tx, batch.ID, batch.Remaining)
		if err != nil {
			break
		}
		err = bs.LedgersDB.Insert(ctx, tx, entry)
		if err != nil {
			break
		}
		err = tx.Commit()
		if err != nil {
			break
		}
	}
	if tx != nil {
		defer tx.Rollback()
	}
	return err
}
