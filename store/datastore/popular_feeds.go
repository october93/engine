package datastore

import (
	"math"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

func (store *Store) GetCardsForPopularRankSince(t time.Time) ([]*model.PopularRankEntry, error) {
	var cards []*model.PopularRankEntry

	err := store.Select(&cards, `
		SELECT popular_ranks.*
		FROM popular_ranks
			LEFT JOIN cards ON popular_ranks.card_id = cards.id
			LEFT JOIN channels ON cards.channel_id = channels.id
		WHERE cards.deleted_at IS NULL
			AND (cards.shadowbanned_at IS NULL)
			AND cards.created_at >= $1
			AND channels.private = false
	`, t)

	return cards, errors.Wrap(err, "GetCardsForPopularRankSince failed")
}

func (store *Store) UpdatePopularRanksWithList(userID globalid.ID, ids []globalid.ID) error {
	_, err := store.Exec(`DELETE FROM user_popular_feeds where user_id = $1`, userID)
	if err != nil {
		return errors.Wrap(err, "GetCardsForPopularRankSince failed")
	}

	_, err = store.Exec(`
		INSERT INTO user_popular_feeds (user_id, card_id, position)
		SELECT
			$1,
			unnest,
			ROW_NUMBER() OVER (ORDER BY ordinality)
		FROM unnest($2::uuid[])
		WITH ORDINALITY
	`, userID, pq.Array(ids))

	return errors.Wrap(err, "GetCardsForPopularRankSince failed")
}

func (store *Store) UpdatePopularRanksForUser(userID globalid.ID) error {
	votePower := math.Pow(model.PopularRankPower, 1/(1-model.PopularRankPower))
	votePowerOverPower := votePower / model.PopularRankPower

	_, err := store.Exec(`
		INSERT INTO user_popular_feeds (user_id, card_id, position)
			SELECT
				$1,
				card_id,
				ROW_NUMBER() OVER (
					ORDER BY (
						CASE
							WHEN (upvote_count + score_mod + (comment_count * 1.5) - downvote_count) < 0
							THEN ((created_at_timestamp - $5) / $6) - (POWER(-(upvote_count + score_mod + (comment_count * 1.5) - downvote_count) + $2, $4) - $3)
							ELSE ((created_at_timestamp - $5) / $6) + POWER((upvote_count + score_mod + (comment_count * 1.5) - downvote_count) + $2, $4) - $3
						END
					) DESC
				)
			FROM popular_ranks LEFT JOIN cards ON popular_ranks.card_id = cards.id
			WHERE cards.deleted_at IS NULL AND cards.shadowbanned_at IS NULL
		ON CONFLICT (user_id, card_id)
		DO UPDATE SET
			position = EXCLUDED.position
			`, userID, votePower, votePowerOverPower, model.PopularRankPower, model.OctoberUnixOffset, model.TimeScalingFactor)

	if err != nil {
		return errors.Wrap(err, "BuildPopularFeedForUser failed")
	}

	return nil
}

func (store *Store) BuildInitialFeed(userID globalid.ID) error {
	votePower := math.Pow(model.PopularRankPower, 1/(1-model.PopularRankPower))
	votePowerOverPower := votePower / model.PopularRankPower

	_, err := store.Exec(`
		INSERT INTO user_feeds (user_id, card_id, position)
			SELECT
				$1,
				card_id,
				ROW_NUMBER() OVER (
					ORDER BY (
						CASE
							WHEN (upvote_count + score_mod + (comment_count * 1.5) - downvote_count) < 0
							THEN ((created_at_timestamp - $5) / $6) - (POWER(-(upvote_count + score_mod + (comment_count * 1.5) - downvote_count) + $2, $4) - $3)
							ELSE ((created_at_timestamp - $5) / $6) + POWER((upvote_count + score_mod + (comment_count * 1.5) - downvote_count) + $2, $4) - $3
						END
					)
				)
			FROM popular_ranks LEFT JOIN cards ON popular_ranks.card_id = cards.id
			WHERE cards.deleted_at IS NULL AND cards.shadowbanned_at IS NULL
			AND cards.channel_id IN (SELECT channel_id FROM channel_memberships WHERE user_id = $1)
		ON CONFLICT (user_id, card_id)
		DO UPDATE SET
			position = EXCLUDED.position
			`, userID, votePower, votePowerOverPower, model.PopularRankPower, model.OctoberUnixOffset, model.TimeScalingFactor)

	if err != nil {
		return errors.Wrap(err, "BuildPopularFeedForUser failed")
	}

	return nil
}

func (store *Store) GetPopularRankCardsForUser(userID globalid.ID, perPage, page int) ([]*model.Card, error) {
	var cards []*model.Card

	err := store.Select(&cards, `
		SELECT cards.*
		FROM user_popular_feeds
			LEFT JOIN cards ON user_popular_feeds.card_id = cards.id
		WHERE user_popular_feeds.user_id = $1
			AND cards.deleted_at IS NULL
			AND (cards.shadowbanned_at IS NULL OR cards.owner_id = $1)
		ORDER BY user_popular_feeds.position
		LIMIT $2 OFFSET $3
	`, userID, perPage, page*perPage)

	return cards, errors.Wrap(err, "GetFeedCardsFromCurrentTop failed")
}

func (store *Store) GetPopularRankForCard(cardID globalid.ID) (*model.PopularRankEntry, error) {
	user := model.PopularRankEntry{}
	err := store.Get(&user, "SELECT * FROM popular_ranks where card_id = $1", cardID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (store *Store) UpdatePopularRankForCard(cardID globalid.ID, viewCountChange, upCountChange, downCountChange, commentCountChange int64, scoreModChange float64) error {
	_, err := store.Exec(`
			UPDATE popular_ranks
			SET
				views = views + $2,
				upvote_count = upvote_count + $3,
				downvote_count = downvote_count + $4,
				comment_count = comment_count + $5,
				score_mod = score_mod + $6
				WHERE card_id = $1`, cardID, viewCountChange, upCountChange, downCountChange, commentCountChange, scoreModChange)
	return errors.Wrap(err, "UpdatePopularRankForCard failed")
}

func (store *Store) UpdateViewsForCards(cardIDs []globalid.ID) error {
	query, args, err := sqlx.In(`
			UPDATE popular_ranks
			SET
				views = views + 1
			WHERE card_id IN (?)`, cardIDs)
	if err != nil {
		return errors.Wrap(err, "UpdateViewsForCards failed")
	}
	query = store.Rebind(query)

	_, err = store.Exec(query, args...)
	return errors.Wrap(err, "UpdateViewsForCards failed")
}

func (store *Store) UpdateUniqueCommentersForCard(cardID globalid.ID) error {
	_, err := store.Exec(`
		UPDATE popular_ranks
		SET comment_count = (
			SELECT COUNT(DISTINCT(COALESCE(alias_id, owner_id)))
			FROM cards
			WHERE thread_root_id = $1
			OR id = $1
		)
		WHERE id = $1`, cardID)
	return errors.Wrap(err, "UpdatePopularRankForCard failed")
}

func (store *Store) SavePopularRank(m *model.PopularRankEntry) error {
	if m == nil {
		return errors.New("provided model can not be nil")
	}

	tn := time.Now().UTC()
	if m.CreatedAtTimestamp == 0 {
		m.CreatedAtTimestamp = tn.Unix()
	}

	_, err := store.NamedExec(`
	INSERT INTO popular_ranks
	(
		card_id,
		views,
		upvote_count,
		downvote_count,
		comment_count,
		unique_commenters_count,
		score_mod,
		created_at_timestamp
	)
	VALUES
	(
		:card_id,
		:views,
		:upvote_count,
		:downvote_count,
		:comment_count,
		:unique_commenters_count,
		:score_mod,
		:created_at_timestamp
	)
	ON CONFLICT (card_id) DO UPDATE
	SET
	views      = EXCLUDED.views,
	upvote_count      = EXCLUDED.upvote_count,
	downvote_count      = EXCLUDED.downvote_count,
	comment_count      = EXCLUDED.comment_count,
	unique_commenters_count = EXCLUDED.unique_commenters_count,
	score_mod      = EXCLUDED.score_mod
	`, m)

	return errors.Wrap(err, "SaveUserReaction failed")
}
