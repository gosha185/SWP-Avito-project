package service

import (
	"context"

	"github.com/google/uuid"
)

func (bs *BonusService) GetExpiringAvailablePoints(ctx context.Context, user uuid.UUID, days int) (int64, error) {
	sum, err := bs.BatchesDB.GetExpiringSum(ctx, user, days)
	if err != nil {
		return 0, err
	}
	return sum, err
}
