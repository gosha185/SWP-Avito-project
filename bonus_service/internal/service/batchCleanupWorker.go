package service

import "context"

// BatchCleanupWorker deletes fully spent and expired batches.
func (bs *BonusService) BatchCleanupWorker(ctx context.Context) error {
    return bs.DeleteExpiredBatches(ctx)
}
