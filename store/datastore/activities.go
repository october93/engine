package datastore

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

func (store *Store) SaveActivity(m *model.Activity) error {
	return saveActivity(store, m)
}

func (tx *Tx) SaveActivity(m *model.Activity) error {
	return saveActivity(tx, m)
}

func saveActivity(e sqlx.Ext, m *model.Activity) error {
	if m == nil {
		return errors.New("provided model can not be nil")
	}
	if m.ID == globalid.Nil {
		m.ID = globalid.Next()
	}
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now().UTC()
	}
	_, err := sqlx.NamedExec(e,
		`INSERT INTO activities
	(
		id,
		rpc,
		data,
		user_id,
		error,
		created_at
	)
	VALUES
	(
		:id,
		:rpc,
		:data,
		:user_id,
		:error,
		:created_at
	)`, m)
	return errors.Wrap(err, "SaveActivity failed")
}

func (store *Store) GetLastActiveAt(userID globalid.ID) (time.Time, error) {
	var result time.Time
	err := store.Get(&result, `SELECT max(created_at) FROM activities WHERE user_id = $1`, userID)
	return result, errors.Wrap(err, "GetLastActiveAt failed")
}

func (store *Store) BatchGetLastActiveAt(ids []globalid.ID) ([]*model.DBTime, error) {
	query, args, err := sqlx.In(`
		SELECT t.max FROM unnest(ARRAY[?]::uuid[]) WITH ORDINALITY AS r(id, rn) LEFT JOIN (
			SELECT user_id, max(created_at) FROM activities GROUP BY user_id
			) AS t(id, max) USING (id) ORDER BY r.rn`, ids)
	if err != nil {
		return nil, errors.Wrap(err, "BatchGetLastActiveAt failed")
	}
	query = store.Rebind(query)
	times := []*model.DBTime{}
	if err = store.Select(&times, query, args...); err != nil {
		return nil, errors.Wrap(err, "BatchGetLastActiveAt failed")
	}
	if len(times) != len(ids) {
		return nil, fmt.Errorf("BatchGetLastActiveAt failed: unexpected number of times returned, likely requesting an invalid id\n")
	}
	return times, nil
}
