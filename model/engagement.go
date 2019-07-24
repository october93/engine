package model

import (
	"time"

	"github.com/october93/engine/kit/globalid"
)

type Engagement struct {
	EngagedUsers       []*EngagedUser     `json:"engagedUsers"`
	EngagedUsersByType EngagedUsersByType `json:"engagedUsersByType"`
	Count              int                `json:"count"`
	CommentCount       int                `json:"commentCount"`
}

type EngagedUsersByType struct {
	Like    []*EngagedUser `json:"like"`
	Love    []*EngagedUser `json:"love"`
	Haha    []*EngagedUser `json:"haha"`
	Wow     []*EngagedUser `json:"wow"`
	Sad     []*EngagedUser `json:"sad"`
	Angry   []*EngagedUser `json:"angry"`
	Comment []*EngagedUser `json:"comment"`
	Tip     []*TippingUser `json:"tip"`
}

type EngagedUser struct {
	CardID           globalid.ID  `db:"card_id"            json:"-"`
	UserID           globalid.ID  `db:"user_id"            json:"userID,omitempty"`
	AliasID          globalid.ID  `db:"alias_id"           json:"aliasID,omitempty"`
	ProfileImagePath string       `db:"profile_image_path" json:"profileImagePath"`
	Username         string       `db:"username"           json:"username"`
	DisplayName      string       `db:"display_name"       json:"displayName"`
	Type             ReactionType `db:"type"               json:"type"`
	CreatedAt        time.Time    `db:"created_at"         json:"-"`
}

type TippingUser struct {
	TipID            globalid.ID `db:"tip_id"             json:"tipID"`
	UserID           globalid.ID `db:"user_id"            json:"userID,omitempty"`
	AliasID          globalid.ID `db:"alias_id"           json:"aliasID,omitempty"`
	Anonymous        bool        `db:"anonymous"          json:"anonymous"`
	ProfileImagePath string      `db:"profile_image_path" json:"profileImagePath"`
	Username         string      `db:"username"           json:"username"`
	DisplayName      string      `db:"display_name"       json:"displayName"`
	Amount           int         `db:"amount"             json:"amount"`

	CreatedAt time.Time `db:"created_at"         json:"-"`
}
