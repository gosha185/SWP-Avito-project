package models

import "time"

type Leader struct {
	RoleName      string    `db:"role_name"`
	LeaderID  string    `db:"leader_id"`
	UpdatedAt time.Time `db:"updated_at"`
}
