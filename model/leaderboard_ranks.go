package model

import (
	"time"

	"github.com/october93/engine/kit/globalid"
)

type LeaderboardRank struct {
	UserID      globalid.ID `db:"user_id"`
	Rank        int64       `db:"rank"`
	CoinsEarned int64       `db:"coins_earned"`
	CreatedAt   time.Time   `db:"created_at"`
}
