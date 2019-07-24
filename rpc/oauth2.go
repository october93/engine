package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const (
	extendTokenURL     = "https://graph.facebook.com/oauth/access_token?grant_type=fb_exchange_token" // nolint
	meURL              = "https://graph.facebook.com/me?fields=first_name,last_name,email"
	facebookPictureURL = "https://graph.facebook.com/me/picture?redirect=false&height=1080"
)

// OAuth2 handles all the OAuth2 related actions.
type OAuth2 interface {
	ExtendToken(ctx context.Context, token string) (AccessToken, error)
}

type oauth2 struct {
	ClientID     string
	ClientSecret string
}

// NewOAuth2 returns a new instance of oauth2.
func NewOAuth2(clientID, clientSecret string) *oauth2 {
	return &oauth2{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}

// ExtendToken uses a short-lived token to generate a long-lived token. The
// long-lived token is returned to the client and stored in the database for
// subsequent authentications.
func (oa2 *oauth2) ExtendToken(ctx context.Context, token string) (AccessToken, error) {
	u, err := url.Parse(extendTokenURL)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("client_id", oa2.ClientID)
	q.Set("client_secret", oa2.ClientSecret)
	q.Set("fb_exchange_token", token)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	var t accessToken
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&t)
	if err != nil {
		return nil, err
	}
	if t.Error != nil {
		return nil, errors.New(t.Error.Message)
	}
	return &t, nil
}

// AccessToken encapsulates requests to Facebook which require an access token.
type AccessToken interface {
	Token() string
	ExpiresAt() int64
	FacebookUser() (*FacebookUser, error)
}

// accessToken is the response payload for requesting an access token.
type accessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	Error       *struct {
		Message         string `json:"message"`
		Type            string `json:"type"`
		Code            int    `json:"code"`
		ErrorSubcode    int    `json:"error_subcode"`
		FacebookTraceID string `json:"fbtrace_id"`
	} `json:"error,omitempty"`
}

// NewAccessToken returns a new instance of Token.
func NewAccessToken(token string) *accessToken {
	return &accessToken{
		AccessToken: token,
	}
}

func (t *accessToken) Token() string {
	return t.AccessToken
}

func (t *accessToken) ExpiresAt() int64 {
	return t.ExpiresIn
}

// FacebookUser returns profile information related to a user's Facebook
// profile. The user is identified via their access token.
func (t *accessToken) FacebookUser() (*FacebookUser, error) {
	u, err := url.Parse(meURL)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("access_token", t.AccessToken)

	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	var fbu FacebookUser
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&fbu)
	if err != nil {
		return nil, err
	}
	if t.Error != nil {
		return nil, errors.New(t.Error.Message)
	}

	u, err = url.Parse(facebookPictureURL)
	if err != nil {
		return nil, err
	}
	q = u.Query()
	q.Set("access_token", t.AccessToken)
	u.RawQuery = q.Encode()
	resp, err = http.Get(u.String())
	if err != nil {
		return nil, err
	}
	dec = json.NewDecoder(resp.Body)
	var picture pictureResponse
	err = dec.Decode(&picture)
	if err != nil {
		return nil, err
	}
	fbu.ProfileImagePath = picture.Data.URL
	return &fbu, nil
}

// FacebookUser contains profile information related to a user's Facebook
// profile.
type FacebookUser struct {
	ID               string `json:"id"`
	Email            string `json:"email"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	ProfileImagePath string
}

func (f *FacebookUser) Name() string {
	return fmt.Sprintf("%s %s", f.FirstName, f.LastName)
}

type pictureResponse struct {
	Data struct {
		URL string `json:"url"`
	} `json:"data"`
}
