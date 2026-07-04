package service

import (
	"context"
	"log"
)

// DeleteExpiredBatches deletes fully spent and expired batches.
func (bs *BonusService) DeleteExpiredBatches(ctx context.Context) error {
	rowsDeleted, err := bs.BatchesDB.DeleteExpiredZeroBatches(ctx, 30)
	if err != nil {
		return err
	}

	log.Printf("[DeleteExpiredBatches] deleted %d old batches", rowsDeleted)
	return nil
}
