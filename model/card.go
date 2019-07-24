package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/october93/engine/kit/globalid"
)

type IdentityMap map[globalid.ID]globalid.ID

type Card struct {
	ID                  globalid.ID `db:"id"`
	OwnerID             globalid.ID `db:"owner_id"`
	AliasID             globalid.ID `db:"alias_id"`
	ThreadReplyID       globalid.ID `db:"thread_reply_id"`
	ThreadRootID        globalid.ID `db:"thread_root_id"`
	ChannelID           globalid.ID `db:"channel_id"`
	ThreadLevel         int         `db:"thread_level"`
	CoinsEarned         int         `db:"coins_earned"`
	Title               string      `db:"title"`
	Content             string      `db:"content"`
	URL                 string      `db:"url"`
	Anonymous           bool        `db:"anonymous"`
	BackgroundImagePath string      `db:"background_image_path"`
	BackgroundColor     string      `db:"background_color"`
	AuthorToAlias       IdentityMap `db:"author_to_alias"`
	IsIntroCard         bool        `db:"is_intro_card"`

	DeletedAt      dbTime    `db:"deleted_at"`
	ShadowbannedAt dbTime    `db:"shadowbanned_at"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`

	Author *Author
}

func NewCard(id globalid.ID, title, content string) Card {
	n := time.Now().UTC()
	return Card{
		ID:        globalid.Next(),
		OwnerID:   id,
		Title:     title,
		Content:   content,
		CreatedAt: n,
		UpdatedAt: n,
	}
}

// Reply returns whether this card is a reply to another card.
func (c *Card) Reply() bool {
	return c.ThreadReplyID != globalid.Nil
}

// IsComment returns whether the card is participating in a thread.
func (c *Card) IsComment() bool {
	return c.ThreadRootID != globalid.Nil
}

// ReplyTo ensures all the thread related fields are set accordingly to the
// threading model.
func (c *Card) ReplyTo(card *Card) {
	c.ThreadRootID = card.ID
	c.ThreadReplyID = card.ID
	if card.IsComment() {
		c.ThreadRootID = card.ThreadRootID
	}
}

func (c *Card) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Export())
}

func (c *Card) Export() *CardView {
	return &CardView{
		ID:                  c.ID,
		ThreadReplyID:       c.ThreadReplyID,
		ThreadRootID:        c.ThreadRootID,
		ThreadLevel:         c.ThreadLevel,
		Title:               c.Title,
		Content:             c.Content,
		URL:                 c.URL,
		CoinsEarned:         c.CoinsEarned,
		Anonymous:           c.Anonymous,
		BackgroundColor:     c.BackgroundColor,
		BackgroundImagePath: c.BackgroundImagePath,
		CreatedAt:           c.CreatedAt.Unix(),
	}
}

type CardView struct {
	ID                  globalid.ID `json:"cardID"`
	ThreadReplyID       globalid.ID `json:"replyID,omitempty"`
	ThreadRootID        globalid.ID `json:"rootID,omitempty"`
	ThreadLevel         int         `json:"threadLevel"`
	CoinsEarned         int         `json:"coinsEarned"`
	Title               string      `json:"title"`
	Content             string      `json:"body"`
	URL                 string      `json:"url"`
	BackgroundColor     string      `json:"bgColor"`
	BackgroundImagePath string      `json:"background_image_path"`
	Anonymous           bool        `json:"anonymous"`
	CreatedAt           int64       `json:"post_timestamp"`
}

func (c *CardView) ToCard() *Card {
	return &Card{
		ID:                  c.ID,
		ThreadReplyID:       c.ThreadReplyID,
		Title:               c.Title,
		Content:             c.Content,
		URL:                 c.URL,
		CoinsEarned:         c.CoinsEarned,
		Anonymous:           c.Anonymous,
		BackgroundImagePath: c.BackgroundImagePath,
		CreatedAt:           time.Unix(c.CreatedAt, 0).UTC(),
	}
}

func (d *IdentityMap) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}
	return json.Unmarshal(b, d)
}

func (d IdentityMap) Value() (driver.Value, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

type Viewer struct {
	AnonymousAlias         *AnonymousAlias `json:"alias,omitempty"`
	AnonymousAliasLastUsed bool            `json:"aliasLastUsed"`
}

type CardResponse struct {
	Card             *CardView         `json:"card"`
	Author           *Author           `json:"author"`
	Channel          *Channel          `json:"channel,omitempty"`
	Viewer           *Viewer           `json:"viewer,omitempty"`
	Replies          int               `json:"replies"`
	FeaturedComments *FeaturedComments `json:"featuredComments,omitempty"`
	Reactions        *Reaction         `json:"reaction,omitempty"`
	ViewerReaction   *UserReaction     `json:"viewerReaction,omitempty"`
	Engagement       *Engagement       `json:"engagement,omitempty"`
	Score            int64             `json:"score"`
	Subscribed       bool              `json:"subscribed"`
	IsMine           bool              `json:"isMine,omitempty"`
	Vote             *VoteResponse     `json:"vote,omitempty"`
	LatestComment    *CardResponse     `json:"latestComment,omitempty"`
	RankingReason    string            `json:"rankingReason,omitempty"`
}

type FeaturedComment struct {
	Card   *CardView `json:"card"`
	Author *Author   `json:"author"`
	New    bool      `json:"new"`
}

type FeaturedComments struct {
	Comments []*FeaturedComment
}

type VoteResponse struct {
	Type string `json:"type"`
}

type CardsResponse struct {
	Cards    []*CardResponse `json:"cards"`
	NextPage bool            `json:"hasNextPage"`
}

type NotificationComment struct {
	NotificationID globalid.ID `db:"notification_id"`
	CardID         globalid.ID `db:"card_id"`
	CreatedAt      time.Time   `db:"created_at"`
	UpdatedAt      time.Time   `db:"updated_at"`
}
