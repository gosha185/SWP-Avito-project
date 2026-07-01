package service

import (
	"bonus-service/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

func (bs *BonusService) CancelAllExpiredHolds(ctx context.Context) error {
	holds, err := bs.HoldsDB.GetExpiredHolds(ctx)
	if err != nil {
		return err
	}
	var tx *sql.Tx
	for _, hold := range holds {
		tx, err = bs.DB.BeginTx(ctx, nil)
		if err != nil {
			break
		}
		entry := &models.LedgerEntry{OperationType: models.OpCancel, Amount: hold.Amount, CreatedAt: time.Now(), Metadata: json.RawMessage(fmt.Sprintf(`{"order_id": "%s", "expires_at": "%s", "expired": "holds"}`, hold.OrderID.String(), hold.ExpiresAt.Format(time.RFC3339)))}
		_, err = bs.BalancesDB.GetBalanceForUpdate(ctx, tx, hold.UserID)
		if err != nil {
			break
		}
		err = bs.BalancesDB.UpdateBalance(ctx, tx, hold.UserID, -hold.Amount, 0)
		if err != nil {
			break
		}
		err = bs.HoldsDB.UpdateHoldStatus(ctx, tx, hold.ID, "cancelled")
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
