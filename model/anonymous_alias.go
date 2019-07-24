package model

import (
	"time"

	"github.com/october93/engine/kit/globalid"
)

type AnonymousAlias struct {
	ID               globalid.ID `db:"id"`
	DisplayName      string      `db:"display_name"`
	Username         string      `db:"username"`
	ProfileImagePath string      `db:"profile_image_path"`
	CreatedAt        time.Time   `db:"created_at"`
	UpdatedAt        time.Time   `db:"updated_at"`
	Inactive         bool        `db:"inactive"`
}

func (al *AnonymousAlias) Author() *Author {
	return &Author{
		ID:               al.ID,
		DisplayName:      al.DisplayName,
		Username:         al.Username,
		ProfileImagePath: al.ProfileImagePath,
		IsAnonymous:      true,
	}
}

func (al *AnonymousAlias) TaggableUser() TaggableUser {
	return TaggableUser{
		Username:       al.Username,
		DisplayName:    al.DisplayName,
		ProfilePicture: al.ProfileImagePath,
		Anonymous:      true,
	}
}
