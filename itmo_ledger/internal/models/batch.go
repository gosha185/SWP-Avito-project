package models

import (
	"time"

	"github.com/google/uuid"
)

type BonusBatch struct {
	ID        int64     `db:"id" json:"id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	Amount    int64     `db:"amount" json:"amount"`
	Remaining int64     `db:"remaining" json:"remaining"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type ExpiryBreakdown struct {
	DaysLeft int   `json:"days_left"`
	Amount   int64 `json:"amount"`
}
