package model

import "time"

type WaitlistEntry struct {
	Email     string    `db:"email"      json:"email"`
	Name      string    `db:"name"       json:"name"`
	Comment   string    `db:"comment"    json:"comment"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}
