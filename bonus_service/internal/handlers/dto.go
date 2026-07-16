package handlers

import (
	"bonus-service/internal/models"
	"time"

	"github.com/google/uuid"
)

// accrualRequest represents the payload for creating a bonus accrual.
type accrualRequest struct {
	UserID uuid.UUID `json:"user_id"`
	Amount int64     `json:"amount"`
	Days   int       `json:"days"`
}

// accrualResponse contains the ID of the created accrual batch.
type accrualResponse struct {
	BatchId int64 `json:"batch_id"`
}

// createHoldRequest represents the payload for reserving bonus points for an order.
type createHoldRequest struct {
	UserID  uuid.UUID `json:"user_id"`
	OrderID uuid.UUID `json:"order_id"`
	Amount  int64     `json:"amount"`
	Hours   int64     `json:"hours"`
}

// createHoldResponse contains the ID of the created hold.
type createHoldResponse struct {
	HoldID int64 `json:"hold_id"`
}

// confirmHoldResponse contains the ID of the confirmed hold.
type confirmHoldResponse struct {
	HoldID int64 `json:"hold_id"`
}

// cancelHoldResponse contains the ID of the cancelled hold.
type cancelHoldResponse struct {
	HoldID int64 `json:"hold_id"`
}

// getBalanceResponse contains the user's available bonus balance.
type getBalanceResponse struct {
	Available int64 `json:"available"`
}

// getExpirationsResponse contains the amount of bonus points that will expire.
type getExpirationsResponse struct {
	Expiring  int64 `json:"expiring"`
}

// transaction represents a single bonus ledger transaction returned in the history response.
type transaction struct {
	OperationType models.OperationType `json:"operation_type"`
	Amount        int64                `json:"amount"`
	CreatedAt     time.Time            `json:"created_at"`
}

// GetHistoryResponse contains a paginated list of user transactions.
type GetHistoryResponse struct {
	Transactions []transaction `json:"transactions"`
}

// GetHoldResponse contains the total amount of bonus points currently on hold.
type GetHoldResponse struct {
	Amount int64 `json:"amount"`
}
