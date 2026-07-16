package service

import (
	"bonus-service/internal/models"
	"bonus-service/internal/storage"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

/*
validate takes context; user id; number of required available points;
idempotency key; arbitrary number of variables which have to be positive. It
begins new transaction; locks users balance; checks idempotency key uniqueness;
checks all required variables on positiveness; checks if user have enough
points. Returns reference to opened transaction, reference to user's balance,
and error (nil if validation is successful). May be called only within service
methods.
*/
func (bs *BonusService) validate(ctx context.Context, user uuid.UUID, requiredPoints int64, key string, positiveIntegers ...int64) (*sql.Tx, *models.Balance, error) {
	tx, err := bs.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	balance, err := bs.BalancesDB.GetBalanceForUpdate(ctx, tx, user)
	if err != nil {
		tx.Rollback()
		return tx, nil, err
	}

	temp, err := bs.LedgersDB.GetByExternalKey(ctx, key)
	if err != nil {
		tx.Rollback()
		return tx, balance, err
	}

	if temp != nil {
		tx.Rollback()
		return tx, balance, storage.ErrLedgerDuplicate
	}

	for _, i := range positiveIntegers {
		if i < 1 {
			tx.Rollback()
			return tx, balance, storage.ErrIncorrectInput
		}
	}

	if balance.Available < requiredPoints {
		tx.Rollback()
		return tx, balance, storage.ErrInsufficientBalance
	}

	return tx, balance, nil
}
