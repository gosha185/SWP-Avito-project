package service

import (
	"context"

	"github.com/google/uuid"
)

/*
GetHeld takes context; user id. Returns number of all held points of user (0 if
transaction is failed) and error (nil if transaction is successful).
*/
func (bs *BonusService) GetHeld(ctx context.Context, user uuid.UUID) (int64, error) {
	balance, err := bs.BalancesDB.GetBalance(ctx, user)
	if err != nil {
		return 0, err
	}
	return balance.Held, nil
}
