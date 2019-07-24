package model

import (
	"time"

	"github.com/october93/engine/kit/globalid"
)

type UserTip struct {
	ID        globalid.ID `db:"id" json:"id"`
	UserID    globalid.ID `db:"user_id" json:"userID"`
	CardID    globalid.ID `db:"card_id" json:"cardID"`
	AliasID   globalid.ID `db:"alias_id" json:"aliasID"`
	Anonymous bool        `db:"anonymous" json:"anonymous"`
	Amount    int         `db:"amount"   json:"amount"`

	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}
