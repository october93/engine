package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/october93/engine/kit/globalid"
)

type testingusers map[globalid.ID]bool

func (d *testingusers) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}
	return json.Unmarshal(b, d)
}

func (d testingusers) Value() (driver.Value, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

// A link is a URI saved by a user on the client for later retrieval.
type FeatureSwitch struct {
	ID           globalid.ID  `db:"id"`
	Name         string       `db:"name"`
	State        string       `db:"state"` // off, testing, on
	TestingUsers testingusers `db:"testing_users"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at"`
}
