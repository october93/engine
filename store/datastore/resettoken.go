package datastore

import (
	"time"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

// SaveResetToken writes a reset token to the database.
func (store *Store) SaveResetToken(m *model.ResetToken) error {
	return saveResetToken(store, m)
}

// SaveResetToken writes a reset token to the database.
func (tx *Tx) SaveResetToken(m *model.ResetToken) error {
	return saveResetToken(tx, m)
}

func saveResetToken(e sqlx.Ext, m *model.ResetToken) error {
	if m == nil {
		return errors.New("provided model can not be nil")
	}
	tn := time.Now().UTC()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = tn
	}
	m.UpdatedAt = tn
	_, err := sqlx.NamedExec(e,
		`INSERT INTO reset_tokens
		(
			token_hash,
			user_id,
			expires,
			created_at,
			updated_at
		)
		VALUES
		(
			:token_hash,
			:user_id,
			:expires,
			:created_at,
			:updated_at
		)
		ON CONFLICT(user_id) DO UPDATE
		SET
			token_hash = :token_hash,
			user_id    = :user_id,
			expires    = :expires,
			created_at = :created_at,
			updated_at = :updated_at
		WHERE reset_tokens.user_id = :user_id`, m)
	return errors.Wrap(err, "SaveResetToken failed")
}

// GetResetToken reads a reset token from the database.
func (store *Store) GetResetToken(userID globalid.ID) (*model.ResetToken, error) {
	resetToken := model.ResetToken{}
	err := store.Get(&resetToken, "SELECT * FROM reset_tokens where user_id = $1", userID)
	return &resetToken, errors.Wrap(err, "GetResetToken failed")
}
