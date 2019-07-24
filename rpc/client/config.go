package client

import (
	"fmt"
	"strings"
	"time"
)

// Config provides the configuration for the client package.
type Config struct {
	Address string
	APIKey  string
	Timeout time.Duration
}

// NewConfig returns a new instance of Config.
func NewConfig() Config {
	return Config{Timeout: 10 * time.Minute}
}

// Validate checks whether the configuration meets the requirements.
func (c Config) Validate() error {
	if !strings.HasPrefix(c.Address, "ws://") && !strings.HasPrefix(c.Address, "wss://") {
		return fmt.Errorf("invalid address %q, please include protocol. ex: ws://", c.Address)
	}
	return nil
}
