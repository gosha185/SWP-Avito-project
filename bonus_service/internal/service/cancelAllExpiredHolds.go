package service

import (
	"bonus-service/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func (bs *BonusService) CancelAllExpiredHolds(ctx context.Context) error {
	log.Printf("[TTLWorker] Starting...")
	holds, err := bs.HoldsDB.GetExpiredHolds(ctx)
	if err != nil {
		log.Printf("[TTLWorker] GetExpiredHolds error: %v", err)
		return err
	}

	log.Printf("[TTLWorker] Found %d expired holds", len(holds))

	if len(holds) == 0 {
		log.Printf("[TTLWorker] No expired holds to process")
		return nil
	}

	var tx *sql.Tx
	for _, hold := range holds {
		log.Printf("[TTLWorker] Processing hold %d (order_id: %s, amount: %d)", hold.ID, hold.OrderID, hold.Amount)
		tx, err = bs.DB.BeginTx(ctx, nil)
		if err != nil {
			log.Printf("[TTLWorker] BeginTx error: %v", err)
			break
		}
		entry := &models.LedgerEntry{OperationType: models.OpExpiry, Amount: hold.Amount, CreatedAt: time.Now(), Metadata: json.RawMessage(fmt.Sprintf(`{"order_id": "%s", "expires_at": "%s", "expired": "holds"}`, hold.OrderID.String(), hold.ExpiresAt.Format(time.RFC3339))), ExternalKey: fmt.Sprintf("hold_expiry_%d_%d", hold.ID, time.Now().UnixNano())}
		_, err = bs.BalancesDB.GetBalanceForUpdate(ctx, tx, hold.UserID)
		if err != nil {
			log.Printf("[TTLWorker] GetBalanceForUpdate error: %v", err)
			break
		}
		err = bs.HoldsDB.UpdateHoldStatus(ctx, tx, hold.ID, "cancelled")
		if err != nil {
			log.Printf("[TTLWorker] UpdateHoldStatus error: %v", err)
			break
		}
		var holdBatches []models.HoldBatch
		holdBatches, err = bs.HoldBatchesDB.GetHoldBatchesByHoldID(ctx, tx, hold.ID)
		if err != nil {
			log.Printf("[TTLWorker] GetHoldBatchesByHoldID error: %v", err)
			break
		}
		var sum int64
		for _, holdBatch := range holdBatches {
			var validness bool
			validness, err = bs.BatchesDB.IncreaseBatchRemaining(ctx, tx, holdBatch.BatchID, holdBatch.Amount)
			if err != nil {
				log.Printf("[TTLWorker] IncreaseBatchRemaining error (batch_id: %d): %v", holdBatch.BatchID, err)
				break
			}
			if validness {
				sum += holdBatch.Amount
			}
		}
		err = bs.BalancesDB.UpdateBalance(ctx, tx, hold.UserID, sum, -hold.Amount)
		if err != nil {
			log.Printf("[TTLWorker] UpdateBalance error: %v", err)
			break
		}
		err = bs.LedgersDB.Insert(ctx, tx, entry)
		if err != nil {
			log.Printf("[TTLWorker] Insert ledger error: %v", err)
			break
		}
		err = tx.Commit()
		if err != nil {
			log.Printf("[TTLWorker] Commit error: %v", err)
			break
		}
	}
	if tx != nil {
		defer tx.Rollback()
	}
	if err != nil {
		log.Printf("[TTLWorker] Finished with error: %v", err)
	} else {
		log.Printf("[TTLWorker] Finished successfully")
	}
	return err
}
