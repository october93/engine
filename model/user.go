package model

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/october93/engine/kit/globalid"
	"golang.org/x/crypto/bcrypt"
)

const (
	usernameMinLength = 2
	usernameMaxLength = 20
)

var (
	// ErrUsernameTaken is used when a username is validated but it is already in
	// use by another user.
	ErrInsufficientBalance = errors.New("not enough coins to pay cost")
)

var usernameRegex = regexp.MustCompile(fmt.Sprintf("^[a-z0-9_]{%d,%d}$", usernameMinLength, usernameMaxLength))
var ErrWrongPassword = errors.New("password did not match")

type User struct {
	ID                globalid.ID `db:"id" json:"nodeId,omitempty"`
	DisplayName       string      `db:"display_name" json:"displayname"`
	FirstName         string      `db:"first_name" json:"firstName,omitempty"`
	LastName          string      `db:"last_name" json:"lastName,omitempty"`
	ProfileImagePath  string      `db:"profile_image_path" json:"profileimg_path"`
	CoverImagePath    string      `db:"cover_image_path" json:"cover_image_path"`
	Bio               string      `db:"bio" json:"userBio,omitempty"`
	Email             string      `db:"email" json:"email,omitempty"`
	PasswordHash      string      `db:"password_hash"`
	PasswordSalt      string      `db:"password_salt"`
	Username          string      `db:"username" json:"username"`
	Devices           Devices     `db:"devices" json:"devices,omitempty"`
	Admin             bool        `db:"admin" json:"admin,omitempty"`
	AllowEmail        bool        `db:"allow_email" json:"allow_email"`
	SearchKey         string      `db:"search_key" json:"searchKey,omitempty"`
	FeedUpdatedAt     dbTime      `db:"feed_updated_at"`
	JoinedFromInvite  globalid.ID `db:"joined_from_invite" json:"joinedFromInvite"`
	BotchedSignup     bool        `db:"botched_signup" json:"botchedSignup"`
	PossibleUninstall bool        `db:"possible_uninstall" json:"possibleUninstall"`
	GotDelayedInvites bool        `db:"got_delayed_invites" json:"-"`
	DisableFeed       bool        `db:"disable_feed" json:"-"`
	IsVerified        bool        `db:"is_verified" json:"-"`
	IsInternal        bool        `db:"is_internal" json:"-"`

	CoinBalance             int64  `db:"coin_balance" json:"-"`
	TemporaryCoinBalance    int64  `db:"temporary_coin_balance" json:"-"`
	CoinRewardLastUpdatedAt dbTime `db:"coin_reward_last_updated_at" json:"-"`

	IsDefault         bool      `db:"is_default" json:"-"`
	SeenIntroCards    bool      `db:"seen_intro_cards" json:"-"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt         dbTime    `db:"deleted_at"`
	FeedLastUpdatedAt dbTime    `db:"feed_last_updated_at"`
	ShadowbannedAt    dbTime    `db:"shadowbanned_at"`
	BlockedAt         dbTime    `db:"blocked_at" json:"blockedAt"`
}

type TaggableUser struct {
	Username       string `json:"username"`
	DisplayName    string `json:"displayName"`
	ProfilePicture string `json:"profilePicture"`
	Anonymous      bool   `json:"anonymous"`
}

// NewUser returns a new instance of User.
func NewUser(nodeID globalid.ID, username, email, displayName string) *User {
	return &User{
		ID:            nodeID,
		Username:      username,
		Email:         email,
		DisplayName:   displayName,
		BotchedSignup: false,
	}
}

// SetPassword assigns a new password by hashing it.
func (u *User) SetPassword(password string) error {
	hash, err := HashPassword(password)
	if err != nil {
		return err
	}
	u.PasswordHash = hash
	return nil
}

// PasswordMatches validates whether the given password is the same as the
// password stored in the user.
func (u *User) PasswordMatches(password string) (bool, error) {
	hash, err := base64.StdEncoding.DecodeString(u.PasswordHash)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword || err == bcrypt.ErrHashTooShort {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.Export(u.ID))
}

// HashPassword hashes the password using bcrypt.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(hash), nil
}

// ValidateUsername checks if the username satisfies the constraints.
func ValidateUsername(username string, blacklist map[string]bool) error {
	if !usernameRegex.MatchString(username) {
		if len(username) < usernameMinLength {
			return errors.New("username is too short")
		}
		if len(username) > usernameMaxLength {
			return errors.New("username is too long")
		}
		return errors.New("username can only use letters, numbers and underscores")
	}
	if blacklist[username] {
		return errors.New("username is unavailable")
	}
	return nil
}

// Export exports the user by copying the existing fields to the target fields
// of ExportedUser. Based on the given user ID more or less fields are
// exported.
func (u *User) Export(id globalid.ID) *ExportedUser {
	if id == u.ID {
		return &ExportedUser{
			ID:                   u.ID,
			DisplayName:          u.DisplayName,
			Email:                u.Email,
			Username:             u.Username,
			FirstName:            u.FirstName,
			LastName:             u.LastName,
			ProfileImagePath:     u.ProfileImagePath,
			CoverImagePath:       u.CoverImagePath,
			AllowEmail:           u.AllowEmail,
			Bio:                  u.Bio,
			Devices:              u.Devices,
			SearchKey:            u.SearchKey,
			Admin:                u.Admin,
			BotchedSignup:        u.BotchedSignup,
			CoinBalance:          u.CoinBalance,
			TemporaryCoinBalance: u.TemporaryCoinBalance,
		}
	}
	return &ExportedUser{
		ID:               u.ID,
		DisplayName:      u.DisplayName,
		Username:         u.Username,
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		ProfileImagePath: u.ProfileImagePath,
		CoverImagePath:   u.CoverImagePath,
		Bio:              u.Bio,
	}
}

func (u *User) Author() *Author {
	return &Author{
		ID:               u.ID,
		DisplayName:      u.DisplayName,
		Username:         u.Username,
		ProfileImagePath: u.ProfileImagePath,
		IsAnonymous:      false,
	}
}

// ExportedUser is a subset of user. This type exists in order to hide
// sensitive information.
type ExportedUser struct {
	ID               globalid.ID       `json:"nodeId,omitempty"`
	DisplayName      string            `json:"displayname"`
	FirstName        string            `json:"firstName,omitempty"`
	LastName         string            `json:"lastName,omitempty"`
	ProfileImagePath string            `json:"profileimg_path"`
	CoverImagePath   string            `json:"cover_image_path"`
	Bio              string            `json:"userBio"`
	Email            string            `json:"email,omitempty"`
	AllowEmail       bool              `json:"allow_email,omitempty"`
	Username         string            `json:"username"`
	Devices          map[string]Device `json:"devices,omitempty"`
	Admin            bool              `json:"admin,omitempty"`
	SearchKey        string            `json:"searchKey,omitempty"`
	BotchedSignup    bool              `json:"botchedSignup"`

	CoinBalance          int64 `json:"coinBalance,omitempty"`
	TemporaryCoinBalance int64 `json:"temporaryCoinBalance,omitempty"`
}

func (eu *ExportedUser) Import() *User {
	return &User{
		DisplayName:      eu.DisplayName,
		Email:            eu.Email,
		ID:               eu.ID,
		Username:         eu.Username,
		FirstName:        eu.FirstName,
		LastName:         eu.LastName,
		ProfileImagePath: eu.ProfileImagePath,
		CoverImagePath:   eu.CoverImagePath,
		Bio:              eu.Bio,
		AllowEmail:       eu.AllowEmail,
		Devices:          eu.Devices,
		BotchedSignup:    eu.BotchedSignup,
	}
}

func (u *User) TaggableUser() TaggableUser {
	return TaggableUser{
		Username:       u.Username,
		DisplayName:    u.DisplayName,
		ProfilePicture: u.ProfileImagePath,
		Anonymous:      false,
	}
}

type CoinBalances struct {
	CoinBalance          int64 `db:"coin_balance" json:"coinBalance"`
	TemporaryCoinBalance int64 `db:"temporary_coin_balance" json:"temporaryCoinBalance"`
}

func (u *User) CanAffordCost(amount int64) bool {
	return u.CoinBalance+u.TemporaryCoinBalance >= amount
}
