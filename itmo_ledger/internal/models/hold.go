package models

import (
	"time"

	"github.com/google/uuid"
)

type Hold struct {
	ID        int64     `db:"id" json:"id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	OrderID   uuid.UUID `db:"order_id" json:"order_id"`
	Amount    int64     `db:"amount" json:"amount"`
	Status    string    `db:"status" json:"status"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
