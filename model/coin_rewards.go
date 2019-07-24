package model

import (
	"time"

	"github.com/october93/engine/kit/globalid"
)

type CoinReward struct {
	ID             globalid.ID `db:"id"`
	UserID         globalid.ID `db:"user_id"`
	CoinsReceived  int64       `db:"coins_received"`
	LastRewardedOn dbTime      `db:"last_rewarded_on"`
	UpdatedAt      time.Time   `db:"updated_at"`
	CreatedAt      time.Time   `db:"created_at"`
}
