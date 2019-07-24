package model

import (
	"context"
	"time"

	"github.com/october93/engine/kit/globalid"
)

type VoteJob struct {
	ID        globalid.ID `db:"id"`
	NodeID    globalid.ID `db:"node_id"`
	CardID    globalid.ID `db:"card_id"`
	Strength  float64     `db:"strength"`
	Type      VoteType    `db:"type"`
	CreatedAt time.Time   `db:"created_at"`

	context.Context
	Cancel func()
}

func NewVoteJob(c context.Context, nodeID, cardID globalid.ID, voteType VoteType, strength float64) *VoteJob {
	ctx, cancel := context.WithCancel(c)
	return &VoteJob{
		ID:        globalid.Next(),
		NodeID:    nodeID,
		CardID:    cardID,
		Type:      voteType,
		Strength:  strength,
		CreatedAt: time.Now().UTC(),
		Context:   ctx,
		Cancel:    cancel,
	}
}
