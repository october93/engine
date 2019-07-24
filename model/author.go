package model

import "github.com/october93/engine/kit/globalid"

// Author is what is used on the client and application side to display who
// wrote a card. This can be a user or an anonymous alias.
type Author struct {
	ID               globalid.ID `db:"id" json:"nodeId,omitempty"`
	DisplayName      string      `db:"display_name" json:"displayname"`
	Username         string      `db:"username" json:"username"`
	ProfileImagePath string      `db:"profile_image_path" json:"profileimg_path"`
	IsAnonymous      bool        `db:"is_anonymous" json:"isAnonymous"`
}
