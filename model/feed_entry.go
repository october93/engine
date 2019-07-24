package model

import (
	"time"

	"github.com/october93/engine/kit/globalid"
)

type FeedEntry struct {
	ID               globalid.ID `db:"id"`
	UserID           globalid.ID `db:"user_id"`
	CardID           globalid.ID `db:"card_id"`
	FeedSessionIndex int         `db:"feed_session_index"`
	FeedSessionRank  int         `db:"feed_session_rank"`
	LastRankedAt     dbTime      `db:"last_ranked_at"`
	RerankedAt       dbTime      `db:"reranked_at"`
	LastNewContentAt dbTime      `db:"last_new_content_at"`
	CreatedAt        time.Time   `db:"created_at"`
	UpdatedAt        time.Time   `db:"updated_at"`
}

type FeedEntriesByID map[globalid.ID]FeedEntry

func (entry FeedEntry) EarlierInFeed(than FeedEntry) bool {
	if entry.FeedSessionIndex == than.FeedSessionIndex {
		if entry.FeedSessionRank == than.FeedSessionRank {
			return !entry.CreatedAt.Before(than.CreatedAt)
		}
		return entry.FeedSessionRank < than.FeedSessionRank
	}
	return entry.FeedSessionIndex > than.FeedSessionRank
}
