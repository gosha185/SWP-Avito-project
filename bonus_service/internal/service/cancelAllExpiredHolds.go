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

	_, err = bs.BalancesDB.UpdateBalancesForExpiredHolds(ctx, tx)
	if err != nil {
		return err
	}

	_, err = bs.LedgersDB.InsertCancelEntries(ctx, tx)
	if err != nil {
		return err
	}

	_, err = bs.HoldsDB.CancelAllExpiredHolds(ctx, tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}
