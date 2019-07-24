package model

import (
	"github.com/october93/engine/kit/globalid"
)

type AttentionNode struct {
	ID                globalid.ID `json:"id"`
	CardRankTableSize int         `json:"cardRankTableSize"`
	FollowerCount     int         `json:"followerCount"`
	FollowingCount    int         `json:"followingCount"`
}
