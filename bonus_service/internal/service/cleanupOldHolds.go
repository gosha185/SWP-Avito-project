package service

import (
	"context"
	"log"
)

// CleanupOldHolds deletes old confirmed/cancelled holds.
// Runs daily to keep the holds table clean.
// Called by hold cleanup worker.
func (bs *BonusService) CleanupOldHolds(ctx context.Context, daysOld int) error {
	tx, err := bs.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	rowsDeleted, err := bs.HoldsDB.DeleteOldHolds(ctx, daysOld)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Printf("[CleanupOldHolds] deleted %d old holds (older than %d days)",
		rowsDeleted, daysOld)

	return nil
}
