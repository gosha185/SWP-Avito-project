package service

import (
	"bonus-service/internal/models"
	"context"
	"encoding/json"
	"time"
)

func (bs *BonusService) CancelAllExpiredHolds(ctx context.Context) error {
	tx, err := bs.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	entry := &models.LedgerEntry{OperationType: models.OpExpiry, CreatedAt: time.Now(), Metadata: json.RawMessage(`"expired": "holds"`)}
	holds, err := bs.HoldsDB.GetExpiredHolds(ctx, tx)
	if err != nil {
		return err
	}
	for _, hold := range holds {
		_, err := bs.BalancesDB.GetBalanceForUpdate(ctx, tx, hold.UserID)
		if err != nil {
			return err
		}
		err = bs.BalancesDB.UpdateBalance(ctx, tx, hold.UserID, -hold.Amount, 0)
		if err != nil {
			return err
		}
		err = bs.HoldsDB.UpdateHoldStatus(ctx, tx, hold.ID, "cancelled")
		if err != nil {
			return err
		}
		entry.Amount += hold.Amount
	}
	err = bs.LedgersDB.Insert(ctx, tx, entry)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	_, err = bs.HoldsDB.DeleteOldHolds(ctx, 0)
	return err
}
