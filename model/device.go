package model

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"
)

// Device is the representation of the mobile device on which the mobile
// application is run. The platform specifies the operating system, i.e.
// Android, iOS.
type Device struct {
	Token    string
	Platform string
}

type Devices map[string]Device

func NewDevice(token, platform string) *Device {
	return &Device{Token: token, Platform: platform}
}

func (d *Devices) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}
	return json.Unmarshal(b, d)
}

func (d Devices) Value() (driver.Value, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

func (d Devices) UnmarshalJSON(b []byte) error {
	if d == nil {
		d = make(Devices)
	}
	var stuff map[string]Device
	err := json.Unmarshal(b, &stuff)
	if err != nil {
		return err
	}
	for key, value := range stuff {
		d[key] = value
	}
	return nil
}

func (d Devices) UnmarshalText(text []byte) error {
	err := json.Unmarshal(text, &d)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
