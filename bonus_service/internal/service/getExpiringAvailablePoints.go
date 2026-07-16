package service

import (
	"context"

	"github.com/google/uuid"
)

/*
GetExpiringAvailablePoints takes context; user id; number of days. Returns
number of all currently available points which will be expired in given days (0
if transaction is failed) and error (nil if transaction is successful).
*/
func (bs *BonusService) GetExpiringAvailablePoints(ctx context.Context, user uuid.UUID, days int) (int64, error) {
	sum, err := bs.BatchesDB.GetExpiringSum(ctx, user, days)
	if err != nil {
		return 0, err
	}
	return sum, err
}
