package store

import (
	"errors"

	datastore "github.com/october93/engine/store/datastore"
)

// Config provides the configuration required by the store package.
type Config struct {
	Datastore        datastore.Config
	RootUserPassword string
}

// NewConfig returns a new instance of Config.
func NewConfig() Config {
	return Config{}
}

// Validate returns an error if the config is invalid.
func (c Config) Validate() error {
	if c.RootUserPassword == "" {
		return errors.New("RootUserPassword cannot be blank")
	}
	return c.Datastore.Validate()
}
