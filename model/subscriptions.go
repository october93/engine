package model

import (
	"time"

	"github.com/october93/engine/kit/globalid"
)

// A link is a URI saved by a user on the client for later retrieval.
type Subscription struct {
	UserID    globalid.ID `json:"userID" db:"user_id"`
	CardID    globalid.ID `json:"cardID" db:"card_id"`
	Type      string      `json:"type" db:"type"`
	CreatedAt time.Time   `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time   `json:"updatedAt" db:"updated_at"`
}
