package datastore

import (
	"time"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

func (store *Store) SaveSession(m *model.Session) error {
	return saveSession(store, m)
}

func (tx *Tx) SaveSession(m *model.Session) error {
	return saveSession(tx, m)
}

func saveSession(e sqlx.Ext, m *model.Session) error {
	if m == nil {
		return errors.New("provided model can not be nil")
	}
	if m.ID == globalid.Nil {
		m.ID = globalid.Next()
	}
	tn := time.Now().UTC()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = tn
	}
	m.UpdatedAt = tn
	_, err := sqlx.NamedExec(e,
		`INSERT INTO sessions
	(
		id,
		user_id,
		created_at,
		updated_at
	)
	VALUES
	(
		:id,
		:user_id,
		:created_at,
		:updated_at
	)
	ON CONFLICT(id) DO UPDATE
	SET
		id         = :id,
		user_id    = :user_id,
		created_at = :created_at,
		updated_at = :updated_at
	WHERE sessions.ID = :id`, m)
	return errors.Wrap(err, "SaveSession failed")
}

func (store *Store) GetSessions() ([]*model.Session, error) {
	sessions := []*model.Session{}
	err := store.Select(&sessions, `
		SELECT sessions.*,
			   users.id "user.id",
			   users.username "user.username",
			   users.email "user.email",
			   users.display_name "user.display_name"
		FROM sessions
		JOIN users ON sessions.user_id = users.id
		ORDER BY updated_at DESC`)
	return sessions, errors.Wrap(err, "GetSessions failed")
}

func (store *Store) GetSession(id globalid.ID) (*model.Session, error) {
	session := model.Session{}
	err := store.Get(&session,
		`SELECT sessions.*,
			   users.id "user.id",
			   users.username "user.username",
			   users.email "user.email",
			   users.display_name "user.display_name",
			   users.admin "user.admin"
		FROM sessions
		JOIN users on sessions.user_id = users.id
		WHERE sessions.id = $1`, id)
	return &session, errors.Wrap(err, "GetSession failed")
}

func (store *Store) DeleteExpiredSessions() (int64, error) {
	deadline := time.Now().Add(-model.SessionTTL)
	result, err := store.Exec("DELETE FROM sessions WHERE updated_at <= $1", deadline)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (store *Store) DeleteSession(id globalid.ID) error {
	_, err := store.Exec("DELETE FROM sessions WHERE id = $1", id)
	return errors.Wrap(err, "DeleteSession failed")
}

func (store *Store) DeleteSessionsForUser(userID globalid.ID) error {
	_, err := store.Exec("DELETE FROM sessions WHERE user_id = $1", userID)
	return errors.Wrap(err, "DeleteSession failed")
}
