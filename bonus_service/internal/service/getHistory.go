package service

import (
	"bonus-service/internal/models"
	"context"

	"github.com/google/uuid"
)

/*
GetHistory takes context; user id; limit; offset. Returns ledger entries of
given limited amount with given offset.
*/
func (bs *BonusService) GetHistory(ctx context.Context, user uuid.UUID, limit int, offset int) ([]models.LedgerEntry, error) {
	return bs.LedgersDB.GetHistory(ctx, user, limit, offset)
}
