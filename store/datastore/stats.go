package datastore

import (
	"time"

	"github.com/october93/engine/model"
	"github.com/pkg/errors"
)

func (store *Store) GetPostCount(startDate, endDate time.Time) ([]*model.Count, error) {
	var result []*model.Count
	err := store.Select(&result, `SELECT owner_id "user_id", COUNT(*) FROM cards WHERE thread_root_id IS NULL AND created_at BETWEEN $1 AND $2 GROUP BY owner_id`, startDate, endDate)
	return result, errors.Wrap(err, "GetPostCount failed")
}

func (store *Store) GetDaysActive(startDate, endDate time.Time) ([]*model.Count, error) {
	var result []*model.Count
	err := store.Select(&result, `SELECT user_id, COUNT(created_at)
								  FROM   (SELECT DISTINCT user_id, created_at::date FROM activities WHERE error = '' AND created_at BETWEEN $1 AND $2) t
				                  GROUP  BY user_id`, startDate, endDate)
	return result, errors.Wrap(err, "GetDaysActive failed")
}

func (store *Store) GetCommentCount(startDate, endDate time.Time) ([]*model.Count, error) {
	var result []*model.Count
	err := store.Select(&result, `SELECT owner_id "user_id", COUNT(*) FROM cards WHERE thread_root_id IS NOT NULL AND created_at BETWEEN $1 AND $2 GROUP BY owner_id`, startDate, endDate)
	return result, errors.Wrap(err, "GetCommentCount failed")
}

func (store *Store) GetReactedCount(startDate, endDate time.Time) ([]*model.Count, error) {
	var result []*model.Count
	err := store.Select(&result, `SELECT user_id, COUNT(*) FROM user_reactions WHERE created_at BETWEEN $1 AND $2 GROUP BY user_id`, startDate, endDate)
	return result, errors.Wrap(err, "GetReactedCount failed")
}

func (store *Store) GetReceivedReactionsCount(startDate, endDate time.Time) ([]*model.Count, error) {
	var result []*model.Count
	err := store.Select(&result, `SELECT cards.owner_id "user_id", COUNT(*)
	                              FROM   user_reactions
								  JOIN   cards
								  ON     user_reactions.card_id = cards.id
								  WHERE  cards.owner_id != user_reactions.user_id
								  AND    user_reactions.created_at BETWEEN $1 AND $2
								  GROUP BY cards.owner_id`, startDate, endDate)
	return result, errors.Wrap(err, "GetReceivedReactions failed")
}

func (store *Store) GetInvitedCount(startDate, endDate time.Time) ([]*model.Count, error) {
	var result []*model.Count
	err := store.Select(&result, `SELECT   node_id "user_id", COUNT(*) FROM invites
	                              WHERE    id IN (SELECT joined_from_invite FROM users WHERE created_at BETWEEN $1 AND $2)
								  GROUP BY node_id`, startDate, endDate)
	return result, errors.Wrap(err, "GetInviteCount failed")
}

func (store *Store) GetFollowedUsersCount(startDate, endDate time.Time) ([]*model.Count, error) {
	var result []*model.Count
	err := store.Select(&result, `SELECT follower_id "user_id", COUNT(*) FROM user_follows WHERE created_at BETWEEN $1 AND $2 GROUP BY follower_id`, startDate, endDate)
	return result, errors.Wrap(err, "GetFollowedUsersCount failed")
}

func (store *Store) GetFollowedCount(startDate, endDate time.Time) ([]*model.Count, error) {
	var result []*model.Count
	err := store.Select(&result, `SELECT followee_id "user_id", COUNT(*) FROM user_follows WHERE created_at BETWEEN $1 AND $2 GROUP BY followee_id`, startDate, endDate)
	return result, errors.Wrap(err, "GetFollowedCount failed")
}

func (store *Store) GetTotalReplyCount(startDate, endDate time.Time) ([]*model.Count, error) {
	var result []*model.Count
	// TODO (konrad): rename user_id to id
	err := store.Select(&result, `SELECT thread_root_id "user_id", COUNT(*)
                                  FROM     cards
                                  WHERE    created_at BETWEEN $1 AND $2
																	AND thread_root_id IS NOT NULL
                                  GROUP BY thread_root_id
                                  ORDER BY count DESC`, startDate, endDate)
	return result, errors.Wrap(err, "GetTotalReplyCount failed")
}

func (store *Store) GetTotalLikeCount(startDate, endDate time.Time) ([]*model.Count, error) {
	var result []*model.Count
	// TODO (konrad): rename user_id to id
	err := store.Select(&result, `SELECT card_id "user_id", COUNT(*) FROM user_reactions WHERE type = 'like' AND created_at BETWEEN $1 AND $2 GROUP BY card_id ORDER BY count DESC`, startDate, endDate)
	return result, errors.Wrap(err, "GetTotalLikeCount")
}

func (store *Store) GetTotalDislikeCount(startDate, endDate time.Time) ([]*model.Count, error) {
	var result []*model.Count
	// TODO (konrad): rename user_id to id
	err := store.Select(&result, `SELECT card_id "user_id", COUNT(*) FROM user_reactions WHERE type = 'dislike' AND created_at BETWEEN $1 AND $2 GROUP BY card_id ORDER BY count DESC`, startDate, endDate)
	return result, errors.Wrap(err, "GetTotalLikeCount")
}

func (store *Store) GetUniqueUserCommentCount(startDate, endDate time.Time) ([]*model.Count, error) {
	var result []*model.Count
	// TODO (konrad): rename user_id to id
	err := store.Select(&result, `
		SELECT   thread_root_id "user_id", count(DISTINCT (thread_root_id, owner_id))
		FROM     cards
		WHERE    thread_root_id IS NOT NULL
		AND      owner_id NOT IN
				 (
						SELECT id
						FROM   users
						WHERE  username IN ('paul',
											'eugene',
											'chris',
											'konrad',
											'kai',
											'kingsley',
											'juan',
											'tomas',
											'richagoyal',
											'root',
											'october'))
		AND      created_at BETWEEN $1 AND $2
		GROUP BY thread_root_id
		HAVING   count(DISTINCT (thread_root_id, owner_id)) > 0
		ORDER BY count DESC`, startDate, endDate)
	return result, errors.Wrap(err, "GetUniqueUserCommentCount failed")
}
