package models

import "github.com/google/uuid"

type HoldBatch struct {
	HoldID      int64     `db:"hold_id" json:"hold_id"`
	BatchUserID uuid.UUID `db:"batch_user_id" json:"batch_user_id"`
	BatchID     int64     `db:"batch_id" json:"batch_id"`
	Amount      int64     `db:"amount" json:"amount"`
}
