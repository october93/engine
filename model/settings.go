package model

import (
	"time"
)

type Settings struct {
	ID              int       `db:"id"               json:"-"`
	MaintenanceMode bool      `db:"maintenance_mode" json:"maintenanceMode"`
	SignupsFrozen   bool      `db:"signups_frozen"   json:"signupsFrozen"`
	CreatedAt       time.Time `db:"created_at"       json:"-"`
	UpdatedAt       time.Time `db:"updated_at"       json:"-"`
}
