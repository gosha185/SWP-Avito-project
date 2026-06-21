package service

import (
	"context"

	"github.com/google/uuid"
	"simple-ledger.itmo.ru/internal/models"
)

func (bs *BonusService) GetHistory(ctx context.Context, user uuid.UUID, limit int, offset int) ([]models.LedgerEntry, error) {
	return bs.LedgersDB.GetHistory(ctx, user, limit, offset)
}
