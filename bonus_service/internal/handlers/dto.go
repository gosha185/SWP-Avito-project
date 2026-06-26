package handlers

import (
	"bonus-service/internal/models"
	"time"

	"github.com/google/uuid"
)

type accrualRequest struct {
	UserID uuid.UUID `json:"user_id"`
	Amount int64     `json:"amount"`
	Days   int       `json:"days"`
}

type accrualResponse struct {
	BatchId int64 `json:"batch_id"`
}

type createHoldRequest struct {
	UserID  uuid.UUID `json:"user_id"`
	OrderID uuid.UUID `json:"order_id"`
	Amount  int64     `json:"amount"`
	Hours   int64     `json:"hours"`
}

type createHoldResponse struct {
	HoldID int64 `json:"hold_id"`
}

type confirmHoldResponse struct {
	HoldID int64 `json:"hold_id"`
}

type cancelHoldResponse struct {
	HoldID int64 `json:"hold_id"`
}

type getBalanceResponse struct {
	Available int64 `json:"available"`
	Expiring  int64 `json:"expiring"`
}

type transaction struct {
	OperationType models.OperationType `json:"operation_type"`
	Amount        int64                `json:"amount"`
	CreatedAt     time.Time            `json:"created_at"`
}

type GetHistoryResponse struct {
	Transactions []transaction `json:"transactions"`
}

type GetHoldResponse struct {
	Amount int64 `json:"amount"`
}
