package service

import (
	"bonus-service/internal/storage"
	"database/sql"
)

type BonusService struct {
	BalancesDB    *storage.BalanceRepo
	BatchesDB     *storage.BatchRepo
	HoldsDB       *storage.HoldRepo
	LedgersDB     *storage.LedgerRepo
	HoldBatchesDB *storage.HoldBatchRepo
	DB            *sql.DB
}

func NewBonusService(DB *sql.DB) *BonusService {
	return &BonusService{storage.NewBalanceRepo(DB), storage.NewBatchRepo(DB), storage.NewHoldRepo(DB), storage.NewLedgerRepo(DB), storage.NewHoldBatchRepo(DB), DB}
}
