package service

import (
	"bonus-service/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
	"log"
)

func (bs *BonusService) ExpireAllBatches(ctx context.Context) error {
	log.Printf("[ExpireAllBatches] Starting...")

	batches, err := bs.BatchesDB.GetExpiredBatches(ctx)
	if err != nil {
		log.Printf("[ExpireAllBatches] GetExpiredBatches error: %v", err)
		return err
	}

	log.Printf("[ExpireAllBatches] Found %d expired batches", len(batches))

	for _, batch := range batches {
		log.Printf("[ExpireAllBatches] Processing batch %d", batch.ID)

		tx, err := bs.DB.BeginTx(ctx, nil)
		if err != nil {
			log.Printf("[ExpireAllBatches] BeginTx error: %v", err)
			return err
		}

		entry := &models.LedgerEntry{
			OperationType: models.OpExpiry,
			Amount:        batch.Remaining,
			BatchID:       sql.NullInt64{batch.ID, true},
			CreatedAt:     time.Now(),
			Metadata:      json.RawMessage(fmt.Sprintf(`{"expired": "batches", "batch_id": %d}`, batch.ID)),
			ExternalKey:   fmt.Sprintf("batch_expiry_%d_%d", batch.ID, time.Now().UnixNano()),
		}

		_, err = bs.BalancesDB.GetBalanceForUpdate(ctx, tx, batch.UserID)
		if err != nil {
			log.Printf("[ExpireAllBatches] GetBalanceForUpdate error: %v", err)
			tx.Rollback()
			continue
		}

		err = bs.BalancesDB.UpdateBalance(ctx, tx, batch.UserID, -batch.Remaining, 0)
		if err != nil {
			log.Printf("[ExpireAllBatches] UpdateBalance error: %v", err)
			tx.Rollback()
			continue
		}

		err = bs.BatchesDB.DecreaseBatchRemaining(ctx, tx, batch.ID, batch.Remaining)
		if err != nil {
			log.Printf("[ExpireAllBatches] DecreaseBatchRemaining error: %v", err)
			tx.Rollback()
			continue
		}

		err = bs.LedgersDB.Insert(ctx, tx, entry)
		if err != nil {
			log.Printf("[ExpireAllBatches] Insert ledger error: %v", err)
			tx.Rollback()
			continue
		}

		if err = tx.Commit(); err != nil {
			log.Printf("[ExpireAllBatches] Commit error: %v", err)
			tx.Rollback()
			continue
		}

		log.Printf("[ExpireAllBatches] Successfully expired batch %d", batch.ID)
	}

	log.Printf("[ExpireAllBatches] Finished")
	return nil
}