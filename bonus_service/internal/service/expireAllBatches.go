package service

import (
	"context"
)

/*
ExpireAllBatches takes context, sets remaining of all expired batches to 0,
updates balances, inserts ledger entries, and returns error (nil if transaction
is successful).
*/
func (bs *BonusService) ExpireAllBatches(ctx context.Context) error {
	tx, err := bs.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := bs.BatchesDB.ExpireAllBatches(ctx, tx); err != nil {
		return err
	}

	return tx.Commit()
}
