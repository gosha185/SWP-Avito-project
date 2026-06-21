package service

import (
	"database/sql"

	"simple-ledger.itmo.ru/internal/storage"
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
