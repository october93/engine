package model

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/october93/engine/kit/globalid"
	"golang.org/x/crypto/bcrypt"
)

const resetTokenExpiry = 4 * time.Hour

// ResetToken represents a password reset triggered by a user.
type ResetToken struct {
	Token     globalid.ID `db:"-"          json:"token"`
	TokenHash string      `db:"token_hash" json:"tokenHash"`
	UserID    globalid.ID `db:"user_id"    json:"userID"`
	Expires   time.Time   `db:"expires"    json:"expires"`
	CreatedAt time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt time.Time   `db:"updated_at" json:"updated_at"`
}

// NewResetToken returns a new instance of ResetToken.
func NewResetToken(userID globalid.ID) (*ResetToken, error) {
	token := globalid.Next()
	tokenHash, err := HashResetToken(token)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	return &ResetToken{
		Token:     token,
		TokenHash: tokenHash,
		UserID:    userID,
		Expires:   now.Add(resetTokenExpiry),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func HashResetToken(token globalid.ID) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(token), 10)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(hash), nil
}

// Valid checks whether the reset password is still valid. A token to reset the
// password should only be valid for a fixed timespan.
func (rt *ResetToken) Valid() error {
	if rt.Expires.Before(time.Now()) {
		return errors.New("invite token has expired")
	}
	hash, err := base64.StdEncoding.DecodeString(rt.TokenHash)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword(hash, []byte(rt.Token))
}
