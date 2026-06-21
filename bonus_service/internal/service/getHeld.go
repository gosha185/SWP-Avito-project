package service

import (
	"context"

	"github.com/google/uuid"
)

func (bs *BonusService) GetHeld(ctx context.Context, user uuid.UUID) (int64, error) {
	balance, err := bs.BalancesDB.GetBalance(ctx, user)
	if err != nil {
		return 0, err
	}
	return balance.Held, nil
}
