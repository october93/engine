package model

import (
	"time"

	"github.com/october93/engine/kit/globalid"
)

type UserReactionType string

const (
	UserReactionsTableName = "user_reactions"

	ReactionLike    UserReactionType = "like"
	ReactionDislike UserReactionType = "dislike"
)

type UserReaction struct {
	UserID    globalid.ID      `db:"user_id"    json:"userID"`
	CardID    globalid.ID      `db:"card_id"    json:"cardID"`
	AliasID   globalid.ID      `db:"alias_id"   json:"aliasID,omitempty"`
	Type      UserReactionType `db:"type"       json:"type"`
	CreatedAt time.Time        `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time        `db:"updated_at" json:"updatedAt"`
}

func (ur *UserReaction) ToVoteResponse() *VoteResponse {
	typ := string(Down)
	if ur.Type == ReactionLike {
		typ = string(Up)
	}
	return &VoteResponse{
		Type: typ,
	}
}

func (ur *UserReaction) ToCardReaction() *Reaction {
	if ur.Type != ReactionLike {
		return nil
	}
	return &Reaction{
		ID:        ur.UserID,
		NodeID:    ur.UserID,
		CardID:    ur.CardID,
		AliasID:   ur.AliasID,
		Reaction:  Boost,
		CreatedAt: ur.CreatedAt,
		UpdatedAt: ur.UpdatedAt,
	}
}
