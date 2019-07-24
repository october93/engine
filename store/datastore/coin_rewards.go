package datastore

import (
	"github.com/october93/engine/model"
	"github.com/pkg/errors"
)

func (store *Store) GetCoinsReceivedNotificationData(notif *model.Notification) (*model.CoinReward, error) {
	var result model.CoinReward
	err := store.Get(&result,
		`SELECT *
		 FROM coin_rewards
		 WHERE id = $1`, notif.TargetID)
	return &result, errors.Wrap(err, "GetCoinsReceivedNotificationData failed")
}
