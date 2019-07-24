package model

import (
	"sync"
	"time"

	"github.com/october93/engine/kit/globalid"
)

// SessionTTL defines the time to live after which a session is expired.
const SessionTTL = 30 * 24 * time.Hour

// Session identifies an authenticated user.
type Session struct {
	ID        globalid.ID `db:"id"         json:"id"`
	UserID    globalid.ID `db:"user_id"    json:"userID,omitempty"`
	CreatedAt time.Time   `db:"created_at" json:"-"`
	UpdatedAt time.Time   `db:"updated_at" json:"-"`

	sync.RWMutex
	User *User `db:"user" json:"user,omitempty"`
}

// NewSession returns a new instance of session with a unique identifier.
func NewSession(user *User) *Session {
	session := Session{ID: globalid.Next(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
	if user != nil {
		session.User = user
		session.UserID = user.ID
	}
	return &session
}

// User returns the user attached to the session. This method is thread-safe as
// it is accessed for every authenticated RPC action.
func (s *Session) GetUser() *User {
	s.RLock()
	defer s.RUnlock()
	return s.User
}

// Identify attaches a user to the session. This method is thread-safe since
// the user it is accessed for every authenticated RPC actio
func (s *Session) Identify(user *User) {
	s.Lock()
	defer s.Unlock()
	s.User = user
	s.UserID = user.ID
}

func (s *Session) SetUser(user *User) {
	s.Lock()
	defer s.Unlock()
	s.User = user
	if user != nil {
		s.UserID = user.ID
	} else {
		s.UserID = globalid.Nil
	}
}

func (s *Session) Authenticated() bool {
	s.Lock()
	defer s.Unlock()
	return s.UserID != globalid.Nil
}
