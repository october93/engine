package datastore

import (
	"time"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

func (store *Store) SaveScoreModification(m *model.ScoreModification) error {
	return saveScoreModification(store, m)
}

func (tx *Tx) SaveScoreModification(m *model.ScoreModification) error {
	return saveScoreModification(tx, m)
}

func saveScoreModification(e sqlx.Ext, m *model.ScoreModification) error {
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
		`INSERT INTO score_modifications
	(
		id,
		card_id,
		user_id,
		strength,
		created_at,
		updated_at
	)
	VALUES
	(
		:id,
		:card_id,
		:user_id,
		:strength,
		:created_at,
		:updated_at
	)
	ON CONFLICT(id) DO UPDATE
	SET
		id            = :id,
		user_id       = :user_id,
		card_id       = :card_id,
		strength      = :strength,
		created_at    = :created_at,
		updated_at    = :updated_at
	WHERE score_modifications.ID = :id`, m)
	return errors.Wrap(err, "SaveScoreModification failed")
}
