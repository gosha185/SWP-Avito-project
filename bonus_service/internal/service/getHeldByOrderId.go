package service

import (
	"context"

	"github.com/google/uuid"
)

func (bs *BonusService) GetHeldByOrderId(ctx context.Context, user uuid.UUID, order uuid.UUID) (int64, error) {
	tx, err := bs.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	hold, err := bs.HoldsDB.GetHoldByOrderID(ctx, tx, user, order)
	if err != nil {
		return 0, err
	}
	return hold.Amount, nil
}
