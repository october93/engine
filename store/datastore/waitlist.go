package datastore

import (
	"time"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/october93/engine/model"
)

func (store *Store) SaveWaitlistEntry(m *model.WaitlistEntry) error {
	return saveWaitlistEntry(store, m)
}

func (tx *Tx) SaveWaitlistEntry(m *model.WaitlistEntry) error {
	return saveWaitlistEntry(tx, m)
}

func saveWaitlistEntry(e sqlx.Ext, m *model.WaitlistEntry) error {
	if m == nil {
		return errors.New("provided model can not be nil")
	}
	tn := time.Now().UTC()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = tn
	}
	_, err := sqlx.NamedExec(e,
		`INSERT INTO waitlist
		(
			email,
			comment,
			name,
			created_at
		)
		VALUES
		(
			:email,
			:comment,
			:name,
			:created_at
		)
		ON CONFLICT(email) DO UPDATE
		SET
			email = :email,
			comment = :comment,
			name = :name,
			created_at = :created_at
		WHERE waitlist.email = :email`, m)
	return errors.Wrap(err, "SaveWaitlistEntry failed")
}

func (store *Store) GetWaitlist() ([]*model.WaitlistEntry, error) {
	var waitlist []*model.WaitlistEntry
	err := store.Select(&waitlist, "SELECT * FROM waitlist ORDER BY created_at DESC")
	return waitlist, errors.Wrap(err, "GetWaitlist failed")
}

func (store *Store) UpdateWaitlistEntry(email, comment string) error {
	_, err := store.Exec("UPDATE waitlist SET comment = $1 WHERE email = $2", comment, email)
	return errors.Wrap(err, "UpdateWaitlistEntry failed")
}

func (store *Store) DeleteWaitlistEntry(email string) error {
	_, err := store.Exec(`DELETE FROM waitlist WHERE email = $1`, email)
	return errors.Wrap(err, "DeleteWaitlistEntry failed")
}
