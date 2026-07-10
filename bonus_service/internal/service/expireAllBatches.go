package service

import (
	"context"
)

func (bs *BonusService) ExpireAllBatches(ctx context.Context) error {
	tx, err := bs.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = bs.BalancesDB.UpdateBalancesForExpiredBatches(ctx, tx)
	if err != nil {
		return err
	}

	_, err = bs.BatchesDB.ExpireAllBatches(ctx, tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}
