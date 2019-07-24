package rpc

import (
	"errors"
	"fmt"

	"github.com/october93/engine/kit/globalid"
)

var (
	// ErrUsernameTaken is used when a username is validated but it is already in
	// use by another user.
	ErrUsernameTaken = errors.New("username has already been taken")
	// ErrUsernameTaken is used when a username is validated but it is already in
	// use by another user.
	ErrChannelExists = errors.New("a channel with that name already exists")
	// ErrUserAlready exists is used when an attempt is made to create a user
	ErrUserAlreadyExists = errors.New("user already exists")
	// ErrNoImageData is used when no base64 data for upload image was provided
	ErrNoImageData = errors.New("no image data received")
	// ErrWrongPassword is used when the password did not match. The detail of this
	// error should be hidden from the client.
	ErrWrongPassword = errors.New("the username or password did not match")
	// ErrUserBlocked is used when an blocked user tries to authenticate
	ErrUserBlocked = errors.New("this account is currently inactive, please contact support")
	// ErrUserBlocked is used when an blocked user tries to authenticate
	ErrUserAlreadyRedeemedCode = errors.New("you have already redeemed an invite code")
)

// Error wraps internal errors in order to hide implementation detail.
type Error struct {
	msg   string
	cause error
}

// Erorr implements the error interface for Error.
func (e *Error) Error() string {
	return e.msg
}

// ErrCouldNotRetrieveUser is used when a database lookup for the given user
// has failed.
func ErrCouldNotRetrieveUser(cause error) error {
	return &Error{msg: "could not retrieve user", cause: cause}
}

// ErrForbidden is used when a user is authenticated but not authorized to perform this action.
func ErrForbidden() error {
	return &Error{msg: "operation forbidden"}
}

// ErrInvalidInviteToken is used when the given invite token does not exist.
func ErrInvalidInviteToken() error {
	return &Error{msg: "invalid invite token"}
}

// ErrResetTokenExpired is used when the invite token has expired.
func ErrResetTokenExpired() error {
	return &Error{msg: "invite token has expired"}
}

// ErrNodeNotFound is used when a user ID has been passed for which no node
// with the same node ID exists.
func ErrNodeNotFound(id globalid.ID) error {
	return &Error{msg: fmt.Sprintf("node with ID %v not found", id)}
}

// ErrInvalidToken is used when an invalid reset password token has been passed.
func ErrInvalidToken() error {
	return &Error{msg: "invalid token"}
}

func mapErrors(err error) error {
	if err.Error() == `pq: duplicate key value violates unique constraint "users_email_idx"` {
		return ErrUserAlreadyExists
	}
	return err
}
