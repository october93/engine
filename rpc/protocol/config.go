package protocol

type Config struct {
	EnableBugSnag bool
}

func NewConfig() Config {
	return Config{}
}

// Validate returns an error if the config is invalid.
//TODO: (corylanou) add validation
func (c Config) Validate() error {
	return nil
}
