package metrics

import (
	"errors"
	"strings"
)

// Config provides the configuration required by the metric package.
type Config struct {
	Enabled         bool
	DatadogAgentURL string
}

// NewConfig returns a new instance of Config.
func NewConfig() Config {
	return Config{}
}

// Validate returns an error if the config is invalid.
func (c Config) Validate() error {
	if c.Enabled {
		if strings.Count(c.DatadogAgentURL, ":") != 1 {
			return errors.New(`DatadogAgentURL has wrong format, needs to be host:port`)
		}
	}
	return nil
}
