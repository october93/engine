package model

import (
	"time"

	"github.com/october93/engine/kit/globalid"
)

const (
	// FacebookProvider is the constant stored in the database to identify
	// Facebook OAuth 2.0 accounts.
	FacebookProvider = "facebook"
)

// OAuthAccount tracks the user accounts which have been authorized via a OAuth
// 2.0 provider. This information is required in order to allow for subsequent
// authentications to identify the associated user again.
type OAuthAccount struct {
	ID        globalid.ID `db:"id"`
	Provider  string      `db:"provider"`
	Subject   string      `db:"subject"`
	UserID    globalid.ID `db:"user_id"`
	CreatedAt time.Time   `db:"created_at"`
	UpdatedAt time.Time   `db:"updated_at"`
}

// NewOAuthAccount returns a new instance of OAuthAccount.
func NewOAuthAccount(provider, subject string, userID globalid.ID) *OAuthAccount {
	now := time.Now().UTC()
	return &OAuthAccount{
		ID:        globalid.Next(),
		Provider:  provider,
		Subject:   subject,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
