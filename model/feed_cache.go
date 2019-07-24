package model

import (
	"time"

	"github.com/october93/engine/kit/globalid"
)

type FeedCacheItem struct {
	UserID      globalid.ID `db:"user_id"`
	CardID      globalid.ID `db:"card_id"`
	FeedSession int         `db:"feed_session"`
	Position    int         `db:"position"`
	CreatedAt   time.Time   `db:"created_at"`
}

func ToFeedEntries(cacheItems []*FeedCacheItem, newCardIDs []globalid.ID) []*FeedEntry {
	newCards := make(map[globalid.ID]bool)

	for _, card := range newCardIDs {
		newCards[card] = true
	}

	result := make([]*FeedEntry, len(cacheItems))
	for i, cacheItem := range cacheItems {
		result[i] = &FeedEntry{
			CardID:           cacheItem.CardID,
			UserID:           cacheItem.UserID,
			FeedSessionIndex: cacheItem.FeedSession,
			FeedSessionRank:  cacheItem.Position,
		}
		if newCards[cacheItem.CardID] {
			result[i].LastNewContentAt = NewDBTime(time.Now().UTC())
		}
	}
	return result
}
