package service

import (
	"bonus-service/internal/models"
	"context"

	"github.com/google/uuid"
)

func (bs *BonusService) GetHistory(ctx context.Context, user uuid.UUID, limit int, offset int) ([]models.LedgerEntry, error) {
	return bs.LedgersDB.GetHistory(ctx, user, limit, offset)
}
