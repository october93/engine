package push

import "errors"

type Config struct {
	FanOut       bool
	NatsEndpoint string
}

func NewConfig() Config {
	return Config{}
}

// Validate returns an error if the config is invalid.
func (c Config) Validate() error {
	if c.FanOut && c.NatsEndpoint == "" {
		return errors.New("NatsEndpoint cannot be blank")
	}
	return nil
}
