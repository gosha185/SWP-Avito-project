package storage

import "errors"

var (
	ErrBalanceNotFound     = errors.New("balance not found")
	ErrBalanceNotUpdated   = errors.New("balance not updated")
	ErrBatchNotFound       = errors.New("batch not found")
	ErrHoldNotFound        = errors.New("hold not found")
	ErrLedgerDuplicate     = errors.New("duplicate external key")
	ErrIncorrectInput      = errors.New("incorrect input")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrOrderAlreadyHeld    = errors.New("order already has a hold")
)
