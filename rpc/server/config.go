package server

import (
	"errors"

	"github.com/october93/engine/rpc"
	"github.com/october93/engine/rpc/protocol"
)

type Config struct {
	Host       string
	Port       int
	PublicPath string

	RPC      rpc.Config
	Protocol protocol.Config
}

func NewConfig() Config {
	return Config{
		RPC:      rpc.NewConfig(),
		Protocol: protocol.NewConfig(),
	}
}

// Validate returns an error if the config is invalid.
//TODO: (corylanou) add validation
func (c Config) Validate() error {
	if c.PublicPath == "" {
		return errors.New("PublicPath cannot be blank")
	}
	if err := c.RPC.Validate(); err != nil {
		return err
	}
	if err := c.Protocol.Validate(); err != nil {
		return err
	}
	return nil
}
