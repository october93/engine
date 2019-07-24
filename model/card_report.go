package model

import (
	"time"

	"github.com/october93/engine/kit/globalid"
)

// Mention is an @mention in a card's text.
type CardReport struct {
	CardID globalid.ID `json:"cardID" db:"card_id"`
	UserID globalid.ID `json:"userID" db:"user_id"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
