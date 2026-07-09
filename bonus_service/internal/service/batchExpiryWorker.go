package service

import "context"

// BatchExpiryWorker expires old batches.
func (bs *BonusService) BatchExpiryWorker(ctx context.Context) error {
	return bs.ExpireAllBatches(ctx)
}
