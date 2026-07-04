package service

import "context"

// BatchCleanupWorker CleanupWorker expires old batches.
func (bs *BonusService) BatchCleanupWorker(ctx context.Context) error {
	return bs.ExpireAllBatches(ctx)
}
