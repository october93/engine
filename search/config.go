package search

import "errors"

// Config provides the configuration required by the search package.
type Config struct {
	ApplicationID    string
	AlgoliaAPIKey    string
	AlgoliaSearchKey string
	IndexName        string
}

// NewConfig returns a new instance of Config.
func NewConfig() Config {
	return Config{}
}

// Validate returns an error if the config is invalid.
func (c Config) Validate() error {
	if c.ApplicationID != "" && c.AlgoliaAPIKey == "" {
		return errors.New("API key must not be empty if indexer is enabled")
	}
	if c.ApplicationID != "" && c.AlgoliaSearchKey == "" {
		return errors.New("Search key must not be empty if indexer is enabled")
	}
	if c.ApplicationID != "" && c.IndexName == "" {
		return errors.New("Index name must not be empty if indexer is enabled")
	}
	return nil
}
