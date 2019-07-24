package debug

import "errors"

type Config struct {
	Profile bool
	Host    string
	Port    int
}

func NewConfig() Config {
	return Config{}
}

func (c Config) Validate() error {
	if c.Host == "" {
		return errors.New("Host cannot be blank")
	}
	if c.Port == 0 {
		return errors.New("Port cannot be 0")
	}
	return nil
}
