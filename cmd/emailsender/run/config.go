package run

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/october93/engine/worker"
)

type Config struct {
	BugsnagAPIKey string
	LogLevel      string

	Worker worker.Config
}

func NewConfig() *Config {
	return &Config{}
}

// Load loads the config from a TOML file.
func (c *Config) Load(fpath string) error {
	bs, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}
	return c.Parse(string(bs))
}

// Parse parses the config from TOML.
func (c *Config) Parse(input string) error {
	_, err := toml.Decode(input, c)
	return err
}

// NewDefaultConfig returns the config that runs when no config is specified.
func NewDefaultConfig() (*Config, error) {
	c := NewConfig()
	return c, nil
}

// Validate returns an error if the config is invalid.
func (c *Config) Validate() error {
	return nil
}
