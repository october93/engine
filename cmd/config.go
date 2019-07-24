package cmd

import (
	"errors"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/october93/engine/coinmanager"
	"github.com/october93/engine/debug"
	"github.com/october93/engine/gql"
	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/metrics"
	"github.com/october93/engine/rpc/push"
	"github.com/october93/engine/rpc/server"
	"github.com/october93/engine/search"
	"github.com/october93/engine/store"
	datastore "github.com/october93/engine/store/datastore"
	"github.com/october93/engine/worker"
)

type Config struct {
	Environment   string
	LogLevel      string
	BugsnagAPIKey string
	UseGraph      bool

	GraphQL     gql.Config
	Store       store.Config
	Datastore   datastore.Config
	Worker      worker.Config
	Server      server.Config
	Push        push.Config
	Metrics     metrics.Config
	Search      search.Config
	Debug       debug.Config
	CoinManager coinmanager.Config
}

func NewConfig() *Config {
	return &Config{
		GraphQL:     gql.NewConfig(),
		Datastore:   datastore.NewConfig(),
		Store:       store.NewConfig(),
		Worker:      worker.NewConfig(),
		Server:      server.NewConfig(),
		Push:        push.NewConfig(),
		Metrics:     metrics.NewConfig(),
		CoinManager: coinmanager.NewConfig(),
	}
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

	c.LogLevel = log.Info
	c.Server.Host = "localhost"
	c.Server.Port = 9000

	return c, nil
}

// Validate returns an error if the config is invalid.
func (c *Config) Validate() error {
	if c.LogLevel == "" {
		return errors.New("LogLevel is required")
	}

	if err := c.Store.Validate(); err != nil {
		return err
	}

	if err := c.GraphQL.Validate(); err != nil {
		return err
	}

	if err := c.Worker.Validate(); err != nil {
		return err
	}

	if err := c.Server.Validate(); err != nil {
		return err
	}
	if err := c.Metrics.Validate(); err != nil {
		return err
	}
	if err := c.Search.Validate(); err != nil {
		return err
	}
	if err := c.Debug.Validate(); err != nil {
		return err
	}
	if err := c.CoinManager.Validate(); err != nil {
		return err
	}
	return nil
}
