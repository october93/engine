package model

import (
	"encoding/json"
	"regexp"
	"time"

	"github.com/october93/engine/kit/globalid"
)

// Notification represents a user notification which is displayed to the user
// inside the application. A notification informs the user about new content or
// a recent interaction, i.e. someonethe user is following has posted a new
// card.

const (
	// Types
	CommentType           = "comment"
	MentionType           = "mention"
	BoostType             = "boost"
	InviteAcceptedType    = "inviteAccepted"
	AnnouncementType      = "announcement"
	IntroductionType      = "introduction"
	NewInvitesType        = "newInvites"
	FollowType            = "follow"
	CoinsReceivedType     = "coinsReceived"
	PopularPostType       = "popularPost"
	LeaderboardRankType   = "leaderboardRank"
	FirstPostActivityType = "firstPostActivity"

	// ActionTypes
	OpenThreadAction = "openThread"
	// "threadRootID"
	OpenUserProfileAction = "openProfile"
	// "username"
	OpenComposerAction = "openComposer"
	// "startingPrompt"
	OpenInvitesAction = "openInvites"
	// no params
	OpenWalletAction = "openWallet"
	// no params

	NavigateToAction = "navigateTo"
	// mobileRouteName
	// webRouteName
)

type Notification struct {
	ID            globalid.ID `db:"id"`
	UserID        globalid.ID `db:"user_id"`
	TargetID      globalid.ID `db:"target_id"`
	TargetAliasID globalid.ID `db:"target_alias_id"`
	Type          string      `db:"type"`
	SeenAt        dbTime      `db:"seen_at"`
	OpenedAt      dbTime      `db:"opened_at"`
	CreatedAt     time.Time   `db:"created_at"`
	UpdatedAt     time.Time   `db:"updated_at"`
	DeletedAt     dbTime      `db:"deleted_at"`
}

type ExportedNotification struct {
	ID           globalid.ID `json:"id"`
	UserID       globalid.ID `json:"userID"`
	ImagePath    string      `json:"imagePath"`
	Message      string      `json:"message"`
	Timestamp    int64       `json:"timestamp"`
	Seen         bool        `json:"seen"`
	Opened       bool        `json:"opened"`
	ShowOnCardID globalid.ID `json:"showOnCardID,omitempty"`

	Type       string            `json:"type"`
	Action     string            `json:"action"`
	ActionData map[string]string `json:"actionData"`
}

func (n *Notification) MarshalJSON() ([]byte, error) {
	return json.Marshal(n)
}

func (eN *ExportedNotification) PlainMessage() string {
	// TODO (konrad): Regular expressions that do not contain any meta characters (things
	// like `\d`) are just regular strings. Using the `regexp` with such
	// expressions is unnecessarily complex and slow. Functions from the
	// `bytes` and `strings` packages should be used instead.
	removeBold := regexp.MustCompile(`\*\*`)        // nolint: megacheck
	removeBangEscapes := regexp.MustCompile(`\\\!`) // nolint: megacheck
	matchEscapes := regexp.MustCompile(`\\.`)
	removeEscapes := regexp.MustCompile(`\\`) // nolint: megacheck

	plain := removeBold.ReplaceAllString(eN.Message, "")
	plain = removeBangEscapes.ReplaceAllString(plain, "!")
	plain = matchEscapes.ReplaceAllStringFunc(plain, func(m string) string {
		return removeEscapes.ReplaceAllString(m, "")
	})

	return plain
}
