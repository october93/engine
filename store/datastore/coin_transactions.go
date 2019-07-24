package datastore

import (
	"time"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

// SaveCard saves a card
func (store *Store) SaveCoinTransaction(m *model.CoinTransaction) error {
	return saveCoinTransaction(store, m)
}

// SaveCard saves a card
func (tx *Tx) SaveCoinTransaction(m *model.CoinTransaction) error {
	return saveCoinTransaction(tx, m)
}

func saveCoinTransaction(e sqlx.Ext, m *model.CoinTransaction) error {
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
		`INSERT INTO coin_transactions
	(
		id,
		source_user_id,
		recipient_user_id,
		card_id,
		amount,
		type,
		created_at,
		updated_at
	)
	VALUES
	(
		:id,
		:source_user_id,
		:recipient_user_id,
		:card_id,
		:amount,
		:type,
		:created_at,
		:updated_at
	)
	ON CONFLICT(id) DO UPDATE
	SET
    source_user_id      = :source_user_id,
    recipient_user_id   = :recipient_user_id,
    card_id             = :card_id,
		amount              = :amount,
    type                = :type,
    created_at          = :created_at,
    updated_at          = :updated_at
	WHERE coin_transactions.id = :id `, m)
	if err != nil {
		return errors.Wrap(err, "SaveCoinTransaction failed")
	}
	return nil
}
