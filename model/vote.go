package model

import (
	"time"

	"github.com/october93/engine/kit/globalid"
)

type VoteType string

const (
	Up   VoteType = "up"
	Down VoteType = "down"
)

type Vote struct {
	CardID      globalid.ID `db:"card_id"        json:"cardID"`
	UserID      globalid.ID `db:"user_id"        json:"userID"`
	Strength    float64     `db:"strength"       json:"strength"`
	SentToGraph bool        `db:"sent_to_graph"  json:"sentToGraph"`
	Type        VoteType    `db:"type"           json:"type"`
	CreatedAt   time.Time   `db:"created_at"     json:"createdAt"`
}
