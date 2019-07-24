package worker

import "errors"

// Config provides the configuration required by the worker package.
type Config struct {
	EnableBugSnag bool
	EmbedlyToken  string

	ThumborHost string
	ThumborPort int

	Port int

	FCMServerKey string
	AMPKey       string

	NSQDAddress       string
	NSQLookupdAddress string
	SendGridAPIKey    string

	SlackWebhook         string
	GraphMonitorInterval int

	CoinUpdater CoinUpdaterConfig
}

type CoinUpdaterConfig struct {
	LeaderboardUpdateFreqency int64
	LeaderboardPlacesLimit    int
}

// NewConfig returns a new instance of Config.
func NewConfig() Config {
	return Config{}
}

// Validate returns an error if the config is invalid.
func (c Config) Validate() error {
	if c.NSQDAddress == "" {
		return errors.New("NSQDAddress cannot be blank")
	}
	if c.NSQLookupdAddress == "" {
		return errors.New("NSQLookupdAddress cannot be blank")
	}
	if c.SendGridAPIKey == "" {
		return errors.New("SendGridAPIKey cannot be blank")
	}
	return nil
}
