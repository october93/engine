package model

import "github.com/october93/engine/kit/globalid"

type Count struct {
	UserID globalid.ID `db:"user_id" json:"id"`
	Count  int         `db:"count"   json:"count"`
}

type Condition struct {
	ID        globalid.ID `db:"id"        json:"id"`
	Condition bool        `db:"condition" json:"condition"`
}

type UserEngagement struct {
	UserID                 globalid.ID `json:"userID"`
	DaysActive             int         `json:"daysActive"`
	PostCount              int         `json:"postCount"`
	CommentCount           int         `json:"commentCount"`
	ReactedCount           int         `json:"reactedCount"`
	ReceivedReactionsCount int         `json:"receivedReactionsCount"`
	FollowedUsersCount     int         `json:"followedUsersCount"`
	FollowedCount          int         `json:"followedCount"`
	InvitedCount           int         `json:"invitedCount"`
	Score                  float64     `json:"score"`
}
