package gql

import "errors"

type Config struct {
	GraphAddress                 string
	TestPushNotificationIconPath string
}

func NewConfig() Config {
	return Config{}
}

func (c *Config) Validate() error {
	if c.GraphAddress == "" {
		return errors.New("GraphAddress cannot be nil")
	}
	return nil
}
