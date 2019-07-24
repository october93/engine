package model

import (
	"time"

	"github.com/october93/engine/kit/globalid"
)

type ScoreModification struct {
	ID        globalid.ID `db:"id"`
	CardID    globalid.ID `db:"card_id"`
	UserID    globalid.ID `db:"user_id"`
	Strength  float64     `db:"strength"`
	CreatedAt time.Time   `db:"created_at"`
	UpdatedAt time.Time   `db:"updated_at"`
}
