package datastore

import "errors"

// Config provides the configuration required by the store package.
type Config struct {
	Database           string   `toml:"database"`
	Host               string   `toml:"host"`
	Port               int      `toml:"port"`
	User               string   `toml:"user"`
	Password           string   `toml:"password"`
	MigrationPath      string   `toml:"migrationPath"`
	Environment        string   `toml:"environment"`
	MaxConnections     int      `toml:"maxConnections"`
	ImageHost          string   `toml:"imageHost"`
	AnonymousIconsPath string   `toml:"anonymousIconsPath"`
	Tags               []string `toml:"tags"`
}

// NewConfig returns a new instance of Config.
func NewConfig() Config {
	return Config{}
}

// NewTestConfig returns a config used for testing
func NewTestConfig() Config {
	return Config{
		Database: "engine_test",
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
	}
}

// Validate returns an error if the config is invalid.
func (c Config) Validate() error {
	if c.Database == "" {
		return errors.New("Database cannot be blank")
	}
	if c.Host == "" {
		return errors.New("Host cannot be blank")
	}
	if c.Port == 0 {
		return errors.New("Port cannot be blank")
	}
	if c.User == "" {
		return errors.New("User cannot be blank")
	}
	if c.MigrationPath == "" {
		return errors.New("SchemaPath cannot be blank")
	}
	if c.Environment == "" {
		return errors.New("Environment cannot be blank")
	}
	if c.MaxConnections < 0 {
		return errors.New("MaxConnections cannot be negative")
	}
	if c.MaxConnections == 0 {
		return errors.New("MaxConnections should not be set to 0 (unlimited)")
	}
	if c.ImageHost == "" {
		return errors.New("ImageHost cannot be blank")
	}
	if c.AnonymousIconsPath == "" {
		return errors.New("AnonymousIconsPath cannot be blank")
	}
	return nil
}
