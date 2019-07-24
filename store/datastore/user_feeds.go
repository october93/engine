package datastore

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

func (store *Store) AddCardsToTopOfFeed(userID globalid.ID, cardRankedIDs []globalid.ID) error {
	// reverse the array so it gets appended in the right order
	for i := len(cardRankedIDs)/2 - 1; i >= 0; i-- {
		opp := len(cardRankedIDs) - 1 - i
		cardRankedIDs[i], cardRankedIDs[opp] = cardRankedIDs[opp], cardRankedIDs[i]
	}
	var currentTop int
	err := store.Get(&currentTop, `SELECT COALESCE(MAX(position), 0) FROM user_feeds WHERE user_id = $1`, userID)
	if err != nil {
		return errors.Wrap(err, "AddCardsToTopOfFeed failed")
	}

	_, err = store.Exec(`
		INSERT INTO user_feeds (user_id, card_id, position)
			SELECT
				$1,
				unnest as card_id,
				(ROW_NUMBER() OVER (ORDER BY ordinality)) + $4
			FROM unnest($2::uuid[])
			WITH ORDINALITY
			LEFT JOIN cards ON unnest = cards.id
			WHERE cards.deleted_at IS NULL
				AND ( cards.shadowbanned_at IS NULL OR cards.owner_id = $1)
				AND cards.owner_id NOT IN (SELECT blocked_user FROM user_blocks WHERE user_id = $1 AND blocked_user IS NOT NULL)
				AND (cards.alias_id NOT IN (SELECT blocked_alias FROM user_blocks WHERE user_id = $1 AND for_thread = unnest) OR alias_id IS NULL)
	  ON CONFLICT (user_id, card_id)
	  DO UPDATE SET
			position = EXCLUDED.position,
			rank_updated_at = $3,
			updated_at = $3
			`, userID, pq.Array(cardRankedIDs), time.Now().UTC(), currentTop)
	if err != nil {
		return errors.Wrap(err, "AddCardsToTopOfFeed failed")
	}

	return nil
}

func (store *Store) NewContentAvailableForUser(userID, cardID globalid.ID) (bool, error) {
	var newAvail bool
	err := store.Get(&newAvail, `
		SELECT CASE
	    WHEN rank_updated_at IS NULL THEN FALSE
	    WHEN last_visited_at IS NULL AND rank_updated_at IS NOT NULL THEN TRUE
	    ELSE rank_updated_at > last_visited_at END
		FROM user_feeds WHERE user_id = $1 AND card_id = $2 LIMIT 1
		`, userID, cardID)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, errors.Wrap(err, "NewContentAvailableForUser failed")
	}
	return newAvail, nil
}

func (store *Store) NewContentAvailableForUserByCards(userID globalid.ID, cardIDs []globalid.ID) (map[globalid.ID]bool, error) {
	var conditions []*model.Condition
	err := store.Select(&conditions, `SELECT card_id "id",
                                      CASE
									  WHEN rank_updated_at IS NULL
										   THEN FALSE
                                      WHEN last_visited_at IS NULL AND rank_updated_at IS NOT NULL
                                           THEN TRUE
                                      ELSE rank_updated_at > last_visited_at
									  END "condition"
                                      FROM user_feeds
									  WHERE user_id = $1
									  AND card_id = ANY($2::uuid[])
									  GROUP BY card_id, rank_updated_at, last_visited_at`, userID, pq.Array(cardIDs))
	if err != nil {
		return nil, errors.Wrap(err, "NewContentAvailableForUserByCards failed")
	}
	result := make(map[globalid.ID]bool, len(cardIDs))
	for _, c := range conditions {
		result[c.ID] = c.Condition
	}
	return result, nil
}

func (store *Store) SetCardVisited(userID, cardID globalid.ID) error {
	_, err := store.Exec(`UPDATE user_feeds SET last_visited_at = $3 WHERE user_id = $1 AND card_id = $2`, userID, cardID, time.Now().UTC())
	if err != nil {
		return errors.Wrap(err, "SetCardVisited failed")
	}

	return nil
}

func (store *Store) ResetUserFeedTop(userID globalid.ID) error {
	_, err := store.Exec(`UPDATE user_feeds SET current_top = false WHERE current_top = true AND user_id = $1`, userID)
	if err != nil {
		return errors.Wrap(err, "ResetUserFeedTop failed")
	}

	_, err = store.Exec(`UPDATE user_feeds SET current_top = true WHERE position = (SELECT MAX(position) FROM user_feeds WHERE user_id = $1)`, userID)
	if err != nil {
		return errors.Wrap(err, "ResetUserFeedTop failed")
	}
	return nil
}

// GetFeedCardsFromCurrentTop fetches all the cards which are in a users current feed.
func (store *Store) GetFeedCardsFromCurrentTop(userID globalid.ID, perPage, page int) ([]*model.Card, error) {
	var cards []*model.Card

	err := store.Select(&cards, `
		SELECT cards.*
		FROM user_feeds LEFT JOIN cards ON user_feeds.card_id = cards.id
		WHERE user_feeds.user_id = $1
			AND user_feeds.position <= COALESCE((SELECT position FROM user_feeds WHERE user_id = $1 AND current_top = true ORDER BY position DESC LIMIT 1), (SELECT MAX(position) FROM user_feeds WHERE user_id = $1))
			AND cards.deleted_at IS NULL
			AND (cards.shadowbanned_at IS NULL OR cards.owner_id = $1)
			AND cards.owner_id NOT IN (SELECT blocked_user FROM user_blocks WHERE user_id = $1 AND blocked_user IS NOT NULL)
			AND (cards.alias_id NOT IN (SELECT blocked_alias FROM user_blocks WHERE user_id = $1 AND for_thread = user_feeds.card_id) OR alias_id IS NULL)
		ORDER BY user_feeds.position DESC, created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, perPage, page*perPage)

	return cards, errors.Wrap(err, "GetFeedCardsFromCurrentTop failed")
}

// GetFeedCardsFromCurrentTop fetches all the cards which are in a users current feed.
func (store *Store) GetFeedCardsFromCurrentTopWithQuery(userID globalid.ID, perPage, page int, searchString string) ([]*model.Card, error) {
	var cards []*model.Card

	err := store.Select(&cards, `
		SELECT cards.*
		FROM user_feeds
			LEFT JOIN cards ON user_feeds.card_id = cards.id
			LEFT JOIN users ON cards.owner_id = users.id
			LEFT JOIN anonymous_aliases ON cards.alias_id = anonymous_aliases.id
			LEFT JOIN channels ON cards.channel_id = channels.id
		WHERE user_feeds.user_id = $1
			AND user_feeds.position <= COALESCE((SELECT position FROM user_feeds WHERE user_id = $1 AND current_top = true), (SELECT MAX(position) FROM user_feeds WHERE user_id = $1))
			AND cards.deleted_at IS NULL
			AND (cards.shadowbanned_at IS NULL OR cards.owner_id = $1)
			AND cards.owner_id NOT IN (SELECT blocked_user FROM user_blocks WHERE user_id = $1 AND blocked_user IS NOT NULL)
			AND (cards.alias_id NOT IN (SELECT blocked_alias FROM user_blocks WHERE user_id = $1 AND for_thread = user_feeds.card_id) OR alias_id IS NULL)
			AND (cards.content ILIKE $4 OR (cards.alias_id IS NULL AND (users.username ILIKE $4 OR users.display_name ILIKE $4)) OR anonymous_aliases.username ILIKE $4 OR channels.handle ILIKE $4)
		ORDER BY user_feeds.position DESC, created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, perPage, page*perPage, fmt.Sprintf(`%%%s%%`, searchString))

	return cards, errors.Wrap(err, "GetFeedCardsFromCurrentTop failed")
}

func (store *Store) GetRankableCardsForUser(userID globalid.ID) ([]*model.PopularRankEntry, error) {
	cards := []*model.PopularRankEntry{}

	err := store.Select(&cards, `
		WITH
			muted_users AS (SELECT muted_user_id FROM user_mutes WHERE user_id = $1),
			followed_users AS (SELECT followee_id FROM user_follows WHERE follower_id = $1),
			followed_channels AS (SELECT channel_id FROM channel_memberships WHERE user_id = $1)
		SELECT
			popular_ranks.*
		FROM
			(
				SELECT card_id as id
				FROM user_card_ranks
					LEFT JOIN cards ON user_card_ranks.card_id = cards.id
				WHERE user_card_ranks.user_id = $1
					AND cards.deleted_at IS NULL
					AND cards.shadowbanned_at IS NULL
				UNION
				SELECT id FROM cards
				WHERE cards.created_at >= COALESCE((SELECT feed_last_updated_at FROM users WHERE id = $1) , now())
				AND cards.thread_root_id IS NULL
				AND cards.owner_id NOT IN (SELECT * FROM muted_users)
				AND cards.deleted_at IS NULL
				AND cards.shadowbanned_at IS NULL
				AND (
					(cards.owner_id IN (SELECT * FROM followed_users) AND cards.alias_id IS NULL)
					OR
					cards.channel_id IN (SELECT * FROM followed_channels)
				)
			) as cards
			LEFT JOIN popular_ranks ON cards.id = popular_ranks.card_id
	`, userID)

	return cards, errors.Wrap(err, "GetRankableCardsForUser failed")
}

func (store *Store) UpdateCardRanksForUser(userID globalid.ID, cardIDs []globalid.ID) error {
	if len(cardIDs) == 0 {
		return nil
	}
	_, err := store.Exec(`DELETE FROM user_card_ranks WHERE user_id = $1`, userID)
	if err != nil {
		return errors.Wrap(err, "UpdateCardRanksForUser failed")
	}

	if len(cardIDs) == 1 {
		_, err = store.Exec(`
			INSERT INTO user_card_ranks (user_id, card_id)
			VALUES ($1, $2)
			ON CONFLICT (user_id, card_id) DO NOTHING
			`, userID, cardIDs[0])
		return errors.Wrap(err, "UpdateCardRanksForUser failed")
	}

	_, err = store.Exec(`
		INSERT INTO user_card_ranks (user_id, card_id)
		SELECT $1, unnest
		FROM unnest($2::uuid[])
		ON CONFLICT (user_id, card_id) DO NOTHING
		`, userID, pq.Array(cardIDs))
	return errors.Wrap(err, "UpdateCardRanksForUser failed")
}

func (store *Store) SetFeedLastUpdatedForUser(userID globalid.ID, t time.Time) error {
	_, err := store.Exec(`UPDATE users SET feed_last_updated_at = $2 WHERE id = $1`, userID, t)
	if err != nil {
		return errors.Wrap(err, "SetFeedLastUpdatedForUser failed")
	}
	return nil
}

func (store *Store) GetCardsInFeed(userID globalid.ID) ([]globalid.ID, error) {
	var cards []globalid.ID
	err := store.Select(&cards, `SELECT card_id FROM user_feeds LEFT JOIN cards ON user_feeds.card_id = cards.id WHERE user_feeds.user_id = $1 AND cards.deleted_at IS NULL`, userID)
	return cards, errors.Wrap(err, "GetCardsInFeed failed")
}

func (store *Store) CountCardsInFeed(userID globalid.ID) (int64, error) {
	var ct int64
	err := store.Get(&ct, `SELECT count(*) FROM user_feeds WHERE user_id = $1`, userID)
	return ct, errors.Wrap(err, "GetCardsInFeed failed")
}

func (store *Store) DeleteCardFromFeeds(cardID globalid.ID) error {
	_, err := store.Exec(`DELETE FROM user_feeds WHERE card_id = $1`, cardID)
	if err != nil {
		return errors.Wrap(err, "DeleteCardFromFeeds failed")
	}

	return nil
}

func (store *Store) GetRankableCardsForExampleFeed(channels []globalid.ID, postedAfter time.Time) ([]*model.PopularRankEntry, error) {
	cards := []*model.PopularRankEntry{}

	err := store.Select(&cards, `
		SELECT
			popular_ranks.*
		FROM
			popular_ranks
			LEFT JOIN cards ON popular_ranks.card_id = cards.id
		WHERE
			cards.created_at >= $1
			AND cards.channel_id IN ($2)
	`, postedAfter, pq.Array(channels))

	return cards, errors.Wrap(err, "GetRankableCardsForUser failed")
}
