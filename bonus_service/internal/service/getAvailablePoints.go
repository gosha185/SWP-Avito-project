package service

import (
	"context"

	"github.com/google/uuid"
)

/*
GetAvailablePoints takes context; user id. Returns number of all available
points (0 if transaction is failed) and error (nil if transaction is
successful).
*/
func (bs *BonusService) GetAvailablePoints(ctx context.Context, user uuid.UUID) (int64, error) {
	balance, err := bs.BalancesDB.GetBalance(ctx, user)
	if err != nil {
		return 0, err
	}
	return balance.Available, nil
}
