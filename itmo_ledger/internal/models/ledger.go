package models

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type OperationType string

const (
	OpAccrual OperationType = "accrual"
	OpHold    OperationType = "hold"
	OpConfirm OperationType = "confirm"
	OpCancel  OperationType = "cancel"
	OpExpiry  OperationType = "expiry"
)

type LedgerEntry struct {
	ID            int64           `db:"id" json:"id"`
	UserID        uuid.UUID       `db:"user_id" json:"user_id"`
	OperationType OperationType   `db:"operation_type" json:"operation_type"`
	Amount        int64           `db:"amount" json:"amount"`
	BatchID       sql.NullInt64   `db:"batch_id" json:"batch_id,omitempty"`
	ExternalKey   string          `db:"external_key" json:"external_key"`
	CreatedAt     time.Time       `db:"created_at" json:"created_at"`
	Metadata      json.RawMessage `db:"metadata" json:"metadata,omitempty"`
}
