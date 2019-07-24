package rpc

import (
	"errors"
)

type Config struct {
	Environment    string
	AdminPanelHost string
	WebEndpoint    string
	APIKey         string
	NatsEndpoint   string

	EmbedlyToken string
	DiffbotToken string

	ImageHost              string
	ProfileImagesPath      string
	CoverImagesPath        string
	CardImagesPath         string
	OriginalCardImagesPath string
	CardContentImagesPath  string
	AssetsPath             string
	PublicImagesPath       string
	BackgroundImagesPath   string

	AccessKeyID     string
	AccessKeySecret string
	S3Upload        bool
	S3Region        string
	S3Bucket        string

	FacebookAppID     string
	FacebookAppSecret string

	PushNotificationDelay int64
	SystemIconPath        string

	ResurfaceActivityThreshold int
	DisableUnverifiedUsers     bool
	AuthyAPIKey                string
	AutoVerify                 bool

	// Tokens
	UnitsPerCoin int64
}

func NewConfig() Config {
	return Config{}
}

// Validate returns an error if the config is invalid.
func (c Config) Validate() error {
	if c.WebEndpoint == "" {
		return errors.New("WebEndpoint cannot be blank")
	}
	if c.EmbedlyToken == "" {
		return errors.New("EmbedlyToken cannot be blank")
	}
	if c.ImageHost == "" {
		return errors.New("ImageHost cannot be blank")
	}
	if c.ProfileImagesPath == "" {
		return errors.New("ProfileImagesPath cannot be blank")
	}
	if c.CardImagesPath == "" {
		return errors.New("CardImagesPath cannot be blank")
	}
	if c.OriginalCardImagesPath == "" {
		return errors.New("OriginalCardImagesPath cannot be blank")
	}
	if c.CardContentImagesPath == "" {
		return errors.New("CardContentImagesPath cannot be blank")
	}
	if c.BackgroundImagesPath == "" {
		return errors.New("BackgroundImagesPath cannot be blank")
	}
	if c.AssetsPath == "" {
		return errors.New("AssetsPath cannot be blank")
	}
	if c.FacebookAppID == "" {
		return errors.New("FacebookAppID cannot be blank")
	}
	if c.FacebookAppSecret == "" {
		return errors.New("FacebookAppSecret cannot be blank")
	}
	if c.S3Upload {
		if c.AccessKeyID == "" {
			return errors.New("AccessKeyID cannot be blank if S3 upload is enabled")
		}
		if c.AccessKeySecret == "" {
			return errors.New("AccessKeySecret cannot be blank if S3 upload is enabled")
		}
		if c.S3Region == "" {
			return errors.New("S3Region cannot be blank if S3 upload is enabled")
		}
		if c.S3Bucket == "" {
			return errors.New("S3Upload cannot be blank if S3 upload is enabled")
		}
	}
	return nil
}
