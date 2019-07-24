package datastore

import (
	"time"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

// SaveInvite saves an invite to the store
func (store *Store) SaveUserTip(m *model.UserTip) error {
	return saveUserTip(store, m)
}

// SaveInvite saves an invite to the store
func (tx *Tx) SaveUserTip(m *model.UserTip) error {
	return saveUserTip(tx, m)
}

func saveUserTip(e sqlx.Ext, m *model.UserTip) error {
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
		`INSERT INTO user_tips
		(
			id,
			user_id,
			card_id,
			alias_id,
			anonymous,
			amount,
			created_at,
			updated_at
		)
		VALUES
		(
			:id,
			:user_id,
			:card_id,
			:alias_id,
			:anonymous,
			:amount,
			:created_at,
			:updated_at
		)
		ON CONFLICT(id) DO UPDATE
		SET
			user_id     = EXCLUDED.user_id,
			card_id     = EXCLUDED.card_id,
			alias_id    = EXCLUDED.alias_id,
			anonymous   = EXCLUDED.anonymous,
			amount      = EXCLUDED.amount,
			created_at  = EXCLUDED.created_at,
			updated_at  = EXCLUDED.updated_at
		`, m)
	return errors.Wrap(err, "SaveUserTip failed")
}

func (store *Store) AssignAliasForUserTipsInThread(userID, threadRootID, aliasID globalid.ID) error {
	_, err := store.Exec(`
		UPDATE user_tips
		SET alias_id=$1
		WHERE
			user_id = $2
			AND
			card_id IN (SELECT id FROM cards WHERE thread_root_id = $3 UNION VALUES($3::uuid))
			AND
			anonymous = true
		`, aliasID, userID, threadRootID)
	if err != nil {
		return errors.Wrap(err, "AssignAliasForUserTipsInThread failed")
	}
	return nil
}
