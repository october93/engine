package datastore

import (
	"time"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/october93/engine/model"
)

func (store *Store) SaveSettings(m *model.Settings) error {
	return saveSettings(store, m)
}

func (tx *Tx) SaveSettings(m *model.Settings) error {
	return saveSettings(tx, m)
}

func saveSettings(e sqlx.Ext, m *model.Settings) error {
	if m == nil {
		return errors.New("provided model can not be nil")
	}
	// enforce one row constraint
	m.ID = 1
	tn := time.Now().UTC()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = tn
	}
	m.UpdatedAt = tn
	_, err := sqlx.NamedExec(e,
		`INSERT INTO settings
		(
			id,
			maintenance_mode,
			signups_frozen,
			created_at,
			updated_at
		)
		VALUES
		(
			:id,
			:maintenance_mode,
			:signups_frozen,
			:created_at,
			:updated_at
		)
		ON CONFLICT(id) DO UPDATE
		SET
			maintenance_mode = :maintenance_mode,
			signups_frozen   = :signups_frozen,
			created_at       = :created_at,
			updated_at       = :updated_at
		WHERE settings.id = :id`, m)
	return errors.Wrap(err, "SaveSettings failed")
}

func (store *Store) GetSettings() (*model.Settings, error) {
	var settings model.Settings
	err := store.Get(&settings, "SELECT * FROM settings LIMIT 1")
	if err != nil {
		return nil, errors.Wrap(err, "GetSettings failed")
	}
	return &settings, nil
}
