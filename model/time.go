package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

type dbTime pq.NullTime

// MarshalJSON marshals the underlying value to a
// proper JSON representation.
func (dt dbTime) MarshalJSON() ([]byte, error) {
	if dt.Valid {
		return json.Marshal(dt.Time)
	}
	return json.Marshal(nil)
}

func (dt *dbTime) UnmarshalJSON(b []byte) error {
	return dt.Time.UnmarshalJSON(b)
}

// Scan implements the Scanner interface.
func (dt *dbTime) Scan(value interface{}) error {
	dt.Time, dt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (dt dbTime) Value() (driver.Value, error) {
	if !dt.Valid {
		return nil, nil
	}
	return dt.Time, nil
}

func (dt dbTime) IsNil() bool {
	return !dt.Valid
}

func (dt dbTime) SameTimeAs(other dbTime) bool {
	return dt.Time == other.Time
}

func NewDBTime(t time.Time) dbTime {
	return dbTime{
		Time:  t,
		Valid: true,
	}
}

func NilDBTime() dbTime {
	return dbTime{
		Time:  time.Now().UTC(),
		Valid: false,
	}
}

//THIS IS SO FUCKING STUPID WHAT THE FUCK
type DBTime struct {
	WEIRDNAME time.Time //TODO: rename this to Time after upgrading graphql
	Valid     bool
}

// MarshalJSON marshals the underlying value to a
// proper JSON representation.
func (dt DBTime) MarshalJSON() ([]byte, error) {
	if dt.Valid {
		return json.Marshal(dt.WEIRDNAME)
	}
	return json.Marshal(nil)
}

func (dt *DBTime) UnmarshalJSON(b []byte) error {
	return dt.WEIRDNAME.UnmarshalJSON(b)
}

// Scan implements the Scanner interface.
func (dt *DBTime) Scan(value interface{}) error {
	dt.WEIRDNAME, dt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (dt DBTime) Value() (driver.Value, error) {
	if !dt.Valid {
		return nil, nil
	}
	return dt.WEIRDNAME, nil
}
