package models

import (
	"time"

	"github.com/google/uuid"
)

type Balance struct {
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	Available int64     `db:"available" json:"available"`
	Held      int64     `db:"held" json:"held"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
