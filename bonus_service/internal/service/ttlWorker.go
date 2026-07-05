package service

import "context"

// TTLWorker cancels expired holds.
func (bs *BonusService) TTLWorker(ctx context.Context) error {
    return bs.CancelAllExpiredHolds(ctx)
}