package service

import "context"

// HoldCleanupWorker deletes old confirmed/cancelled holds.
func (bs *BonusService) HoldCleanupWorker(ctx context.Context) error {
	return bs.CleanupOldHolds(ctx, 30)
}
