package service

import (
	"context"
)

func (bs *BonusService) CancelAllExpiredHolds(ctx context.Context) error {
	tx, err := bs.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := bs.HoldsDB.CancelAllExpiredHolds(ctx, tx); err != nil {
		return err
	}

	return tx.Commit()
}
