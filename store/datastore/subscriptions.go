package datastore

import (
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

func (store *Store) SubscribeToCard(userID, cardID globalid.ID, typ string) error {
	_, err := store.Exec(`
		INSERT INTO subscriptions
		(
			user_id,
			card_id,
			type,
			created_at,
			updated_at
		)
		VALUES
		(
			$1,
			$2,
			$3,
			$4,
			$4
		)
		ON CONFLICT (user_id, card_id, type) DO UPDATE
		SET
			updated_at = $4
		WHERE subscriptions.user_id = $1 AND subscriptions.card_id = $2 AND subscriptions.type = $3`, userID, cardID, typ, time.Now().UTC())

	return errors.Wrap(err, "SubscribeToCard failed")
}

func (store *Store) UnsubscribeFromCard(userID, cardID globalid.ID, typ string) error {
	_, err := store.Exec(`
		DELETE FROM subscriptions
		WHERE user_id = $1 AND card_id = $2 AND type = $3`, userID, cardID, typ)
	return errors.Wrap(err, "UnsubscribeFromCard failed")
}

func (store *Store) SubscribedToTypes(userID, cardID globalid.ID) ([]string, error) {
	var types []string

	err := store.Select(&types, "SELECT type FROM subscriptions WHERE user_id = $1 AND card_id = $2", userID, cardID)
	if err != nil {
		return nil, errors.Wrap(err, "SubscribedToTypes failed")
	}
	return types, nil
}

func (store *Store) SubscribedToCards(userID globalid.ID, cardIDs []globalid.ID) (map[globalid.ID]bool, error) {
	var conditions []*model.Condition
	err := store.Select(&conditions, `SELECT   card_id "id", COUNT(*) > 0 "condition"
	                                  FROM     subscriptions
							          WHERE    user_id = $1
							          AND      card_id = ANY($2::uuid[])
									  GROUP BY card_id`, userID, pq.Array(cardIDs))
	if err != nil {
		return nil, errors.Wrap(err, "SubscribedToCards failed")
	}
	result := make(map[globalid.ID]bool, len(conditions))
	for _, c := range conditions {
		result[c.ID] = c.Condition
	}
	return result, nil
}
