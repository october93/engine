package datastore

import (
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

// GetEngagement derives engagement for a given card from the cards and
// reactions table.
//
// The query is at its core a UNION between reactions and cards table.
// Reactions are filtered by boosts minus reactions to a user's own card. If
// there is a more recent bury from the user, the boost is not included at all.
//
// Only cards that are replies to the given card are included for the
// engagement.
//
// The engagement type is determined by the originating row. It's a 'comment'
// if it's from the cards table and a 'boost' if it is from the reactions
// table. In case there is a comment and a boost, precedence to the comment is
// given (see CASE WHEN statement).
//
// In order to fetch associated user or anonymous alias information JOINS are
// applied. An INNER JOIN for user because the user data is always available
// and a LEFT JOIN on anonymous alias because the anonymous alias data is only
// sometimes available.
//
// In order to protect the identity of an anonymous engagement, a COALESCE
// statement is used to give precedence to the anonymous information.
//
func (store *Store) GetEngagement(cardID globalid.ID) (*model.Engagement, error) {
	var engagement model.Engagement
	var engagedUsers []*model.EngagedUser
	err := store.Select(&engagedUsers, `
	SELECT user_id, username, display_name, alias_id, profile_image_path, type, min(created_at) "created_at" FROM (
			SELECT
				CASE WHEN user_reactions.alias_id IS NOT NULL THEN NULL ELSE user_id END "user_id",
				user_reactions.alias_id,
				COALESCE(
					anonymous_aliases.username,
					users.username
				) "username",
				COALESCE(
					anonymous_aliases.display_name,
					users.display_name
				) "display_name",
				COALESCE(
					anonymous_aliases.profile_image_path,
					users.profile_image_path
				) "profile_image_path",
				user_reactions.created_at,
				user_reactions.type
			FROM       user_reactions
			JOIN       users
			ON         user_reactions.user_id = users.id
			LEFT JOIN  anonymous_aliases
			ON         user_reactions.alias_id = anonymous_aliases.id
			WHERE      card_id = $1
			AND        type = $2
		UNION (
			SELECT
				CASE WHEN alias_id IS NOT NULL THEN NULL ELSE owner_id END "user_id",
				alias_id,
				COALESCE(
					anonymous_aliases.username,
					users.username
				) "username",
				COALESCE(
					anonymous_aliases.display_name,
					users.display_name
				) "display_name",
				COALESCE(
					anonymous_aliases.profile_image_path,
					users.profile_image_path
				) "profile_image_path",
				cards.created_at,
				'comment' "type"
			FROM       cards
			JOIN       users
			ON         cards.owner_id = users.id
			LEFT JOIN  anonymous_aliases
			ON     cards.alias_id = anonymous_aliases.id
			WHERE  (thread_reply_id = $1 OR thread_root_id = $1)
			AND    cards.deleted_at IS NULL
		)
	) t
	GROUP BY user_id, username, display_name, alias_id, profile_image_path, type
	ORDER BY alias_id DESC, type ASC, created_at DESC`, cardID, model.ReactionLike)
	if err != nil {
		return nil, errors.Wrap(err, "GetEngagement failed")
	}

	for _, engagedUser := range engagedUsers {
		engagedUser.CreatedAt = engagedUser.CreatedAt.UTC()
		switch engagedUser.Type {
		case model.Comment:
			engagement.EngagedUsersByType.Comment = append(engagement.EngagedUsersByType.Comment, engagedUser)
		case model.Like:
			engagement.EngagedUsersByType.Like = append(engagement.EngagedUsersByType.Like, engagedUser)
		}
	}

	var tippingUsers []*model.TippingUser
	err = store.Select(&tippingUsers, `
		SELECT
			user_tips.id as tip_id,
			CASE WHEN anonymous = true THEN null ELSE user_id END as user_id,
			alias_id,
			anonymous,
			amount,
			COALESCE(
				CASE
					WHEN anonymous = true THEN anonymous_aliases.profile_image_path
					ELSE users.profile_image_path
				END, ''
			) AS profile_image_path,
			COALESCE(
				CASE
					WHEN anonymous = true THEN anonymous_aliases.display_name
					ELSE users.display_name
				END, ''
			) AS display_name,
			COALESCE(
				CASE
					WHEN anonymous = true THEN anonymous_aliases.username
					ELSE users.username
				END, ''
				) AS username
		FROM user_tips
			LEFT JOIN users ON user_tips.user_id = users.id
			LEFT JOIN anonymous_aliases ON user_tips.alias_id = anonymous_aliases.id
		WHERE card_id = $1
		ORDER BY user_tips.created_at DESC
		`, cardID)
	if err != nil {
		return nil, errors.Wrap(err, "GetEngagement failed")
	}

	if len(tippingUsers) > 0 {
		engagement.EngagedUsersByType.Tip = tippingUsers
	} else if len(engagedUsers) == 0 && len(tippingUsers) == 0 {
		return nil, nil
	}

	var commentCount int
	err = store.Get(&commentCount, `SELECT COUNT(*) FROM cards WHERE CASE WHEN (SELECT thread_root_id FROM cards WHERE id = $1) IS NULL THEN thread_root_id = $1 ELSE thread_reply_id = $1 END AND deleted_at IS NULL`, cardID)

	if err != nil {
		return nil, err
	}

	engagement.CommentCount = commentCount
	engagement.Count = len(engagedUsers)
	return &engagement, nil
}

func (store *Store) GetEngagements(cardIDs []globalid.ID) (map[globalid.ID]*model.Engagement, error) {
	var engagedUsers []*model.EngagedUser
	err := store.Select(&engagedUsers, `
	SELECT card_id, user_id, username, display_name, alias_id, profile_image_path, type, min(created_at) "created_at" FROM (
			SELECT
			        card_id,
				CASE WHEN user_reactions.alias_id IS NOT NULL THEN NULL ELSE user_id END "user_id",
				user_reactions.alias_id,
				COALESCE(
					anonymous_aliases.username,
					users.username
				) "username",
				COALESCE(
					anonymous_aliases.display_name,
					users.display_name
				) "display_name",
				COALESCE(
					anonymous_aliases.profile_image_path,
					users.profile_image_path
				) "profile_image_path",
				user_reactions.created_at,
				user_reactions.type
			FROM       user_reactions
			JOIN       users
			ON         user_reactions.user_id = users.id
			LEFT JOIN  anonymous_aliases
			ON         user_reactions.alias_id = anonymous_aliases.id
			WHERE      card_id = ANY($1::uuid[])
			AND        type = $2
		UNION (
			SELECT
				cards.thread_root_id,
				CASE WHEN alias_id IS NOT NULL THEN NULL ELSE owner_id END "user_id",
				alias_id,
				COALESCE(
					anonymous_aliases.username,
					users.username
				) "username",
				COALESCE(
					anonymous_aliases.display_name,
					users.display_name
				) "display_name",
				COALESCE(
					anonymous_aliases.profile_image_path,
					users.profile_image_path
				) "profile_image_path",
				cards.created_at,
				'comment' "type"
			FROM       cards
			JOIN       users
			ON         cards.owner_id = users.id
			LEFT JOIN  anonymous_aliases
			ON     cards.alias_id = anonymous_aliases.id
			WHERE  thread_root_id = ANY($1::uuid[])
			AND    cards.deleted_at IS NULL
		)
	) t
	GROUP BY card_id, user_id, username, display_name, alias_id, profile_image_path, type
	ORDER BY alias_id DESC, type ASC, created_at DESC`, pq.Array(cardIDs), model.ReactionLike)
	if err != nil {
		return nil, errors.Wrap(err, "GetEngagements failed")
	}

	var counts []*model.Count
	err = store.Select(&counts, `SELECT thread_root_id "user_id", COUNT(*)
	                              FROM     cards
								  WHERE    thread_root_id = ANY($1)
								  AND      deleted_at IS NULL
								  GROUP BY thread_root_id`, pq.Array(cardIDs))
	if err != nil {
		return nil, errors.Wrap(err, "GetEngagements failed")
	}

	type tippingUserByCard struct {
		CardID globalid.ID `db:"card_id"`
		*model.TippingUser
	}
	var tippingUsersByCard []*tippingUserByCard
	err = store.Select(&tippingUsersByCard, `
		SELECT
		    card_id,
			user_tips.id "tip_id",
			CASE WHEN anonymous = true THEN null ELSE user_id END as user_id,
			alias_id,
			anonymous,
			amount,
			COALESCE(
				CASE
					WHEN anonymous = true THEN anonymous_aliases.profile_image_path
					ELSE users.profile_image_path
				END, ''
			) AS profile_image_path,
			COALESCE(
				CASE
					WHEN anonymous = true THEN anonymous_aliases.display_name
					ELSE users.display_name
				END, ''
			) AS display_name,
			COALESCE(
				CASE
					WHEN anonymous = true THEN anonymous_aliases.username
					ELSE users.username
				END, ''
				) AS username
		FROM user_tips
			LEFT JOIN users ON user_tips.user_id = users.id
			LEFT JOIN anonymous_aliases ON user_tips.alias_id = anonymous_aliases.id
			WHERE card_id = ANY($1::uuid[])
		ORDER BY user_tips.created_at DESC
		`, pq.Array(cardIDs))
	if err != nil {
		return nil, errors.Wrap(err, "GetEngagements failed")
	}

	engagements := make(map[globalid.ID]*model.Engagement, len(cardIDs))
	for _, cardID := range cardIDs {
		engagements[cardID] = &model.Engagement{}
	}
	for _, engagedUser := range engagedUsers {
		cardID := engagedUser.CardID
		engagedUser.CreatedAt = engagedUser.CreatedAt.UTC()
		switch engagedUser.Type {
		case model.Comment:
			engagements[cardID].EngagedUsersByType.Comment = append(engagements[cardID].EngagedUsersByType.Comment, engagedUser)
		case model.Like:
			engagements[cardID].EngagedUsersByType.Like = append(engagements[cardID].EngagedUsersByType.Like, engagedUser)
		}
	}
	for _, count := range counts {
		engagements[count.UserID].Count = count.Count
	}

	tipsByCardID := make(map[globalid.ID][]*model.TippingUser, len(cardIDs))
	for _, tip := range tippingUsersByCard {
		if tipsByCardID[tip.CardID] == nil {
			tipsByCardID[tip.CardID] = make([]*model.TippingUser, 1)
		}
		tipsByCardID[tip.CardID] = append(tipsByCardID[tip.CardID], tip.TippingUser)
	}
	for cardID, tippingUsers := range tipsByCardID {
		engagements[cardID].EngagedUsersByType.Tip = tippingUsers
	}

	return engagements, nil
}
