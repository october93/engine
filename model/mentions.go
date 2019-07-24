package model

import (
	"time"

	"github.com/october93/engine/kit/globalid"
)

// Mention is an @mention in a card's text.
type Mention struct {
	ID             globalid.ID `json:"id" db:"id"`
	InCard         globalid.ID `json:"inCard" db:"in_card"`
	MentionedUser  globalid.ID `json:"mentionedUser" db:"mentioned_user"`
	MentionedAlias globalid.ID `json:"mentionedAlias" db:"mentioned_alias"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
	DeletedAt dbTime    `json:"deletedAt" db:"deleted_at"`
}

// links notifications and mentions
type NotificationMention struct {
	NotificationID globalid.ID `json:"notificationID" db:"notification_id"`
	MentionID      globalid.ID `json:"mentionID" db:"mention_id"`
	CreatedAt      time.Time   `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time   `json:"updatedAt" db:"updated_at"`
}
