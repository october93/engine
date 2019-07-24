package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/october93/engine/kit/globalid"
	"github.com/pkg/errors"
)

// Notification represents a user notification which is displayed to the user
// inside the application. A notification informs the user about new content or
// a recent interaction, i.e. someonethe user is following has posted a new
// card.

type AnnouncementActionData map[string]string

type Announcement struct {
	ID      globalid.ID `db:"id" json:"id"`
	UserID  globalid.ID `db:"from_user" json:"userID"`
	CardID  globalid.ID `db:"card_id" json:"cardID"`
	Message string      `db:"message" json:"message"`
	//Action     string                 `db:"action" json:"action"`
	//ActionData AnnouncementActionData `db:"action_data" json:"actionData"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt dbTime    `db:"deleted_at" json:"deletedAt"`
}

func (d *AnnouncementActionData) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}
	return json.Unmarshal(b, d)
}

func (d AnnouncementActionData) Value() (driver.Value, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

func (d AnnouncementActionData) UnmarshalJSON(b []byte) error {
	if d == nil {
		d = make(AnnouncementActionData)
	}
	var stuff map[string]string
	err := json.Unmarshal(b, &stuff)
	if err != nil {
		return err
	}
	for key, value := range stuff {
		d[key] = value
	}
	return nil
}

func (d AnnouncementActionData) UnmarshalText(text []byte) error {
	err := json.Unmarshal(text, &d)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
