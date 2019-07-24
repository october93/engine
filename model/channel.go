package model

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/october93/engine/kit/globalid"
)

var nameRegex = regexp.MustCompile(fmt.Sprintf("^[a-zA-Z0-9_]{%d,%d}$", 2, 24))

type Channel struct {
	ID          globalid.ID `db:"id" json:"id"`
	Name        string      `db:"name" json:"name"`
	Description string      `db:"description" json:"description"`
	Private     bool        `db:"private" json:"private"`
	OwnerID     globalid.ID `db:"owner_id" json:"-"`
	Handle      string      `db:"handle" json:"-"`
	IsDefault   bool        `db:"is_default" json:"-"`

	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}

type ChannelSubscription struct {
	UserID     globalid.ID `db:"user_id"`
	ChannelID  globalid.ID `db:"channel_id"`
	Subscribed bool        `db:"subscribed"`
	Muted      bool        `db:"muted"`
	CreatedAt  time.Time   `db:"created_at" json:"-"`
}

type ChannelUserInfo struct {
	ChannelID   globalid.ID `db:"channel_id" json:"channel"`
	MemberCount int         `db:"member_count" json:"memberCount"`
	Subscribed  bool        `db:"subscribed" json:"subscribed"`
}

type ChannelEngagement struct {
	ChannelID       globalid.ID `db:"channel_id" json:"channelID"`
	TotalPosts      int         `db:"total_posts" json:"totalPosts"`
	TotalLikes      int         `db:"total_likes" json:"totalLikes"`
	TotalDislikes   int         `db:"total_dislikes" json:"totalDislikes"`
	TotalComments   int         `db:"total_comments" json:"totalComments"`
	TotalCommenters int         `db:"total_commenters" json:"totalCommenters"`
}

func ValidateChannelName(name string) error {
	if !nameRegex.MatchString(name) {
		return errors.New("channel names may contain only letters, numbers and underscores, and must be 2-20 characters long")
	}
	return nil
}
