package model

import (
	"database/sql/driver"
	"encoding/json"

	"errors"
)

// Image is a representation of an image hosted on the static file server of
// the backend, i.e. card image, profile image.
type Image struct {
	Path string `json:"path"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
}

type Images []Image

func NewImage(path string, x, y int) *Image {
	return &Image{Path: path, X: x, Y: y}
}

func (i *Images) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}
	return json.Unmarshal(b, i)
}

func (i Images) Value() (driver.Value, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}
