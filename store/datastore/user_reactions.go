package datastore

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

func (store *Store) SaveUserReaction(m *model.UserReaction) error {
	if m == nil {
		return errors.New("provided model can not be nil")
	}

	tn := time.Now().UTC()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = tn
	}
	m.UpdatedAt = tn

	_, err := store.NamedExec(`
	INSERT INTO user_reactions
	(
		user_id,
		card_id,
		alias_id,
		type,
		created_at,
		updated_at
	)
	VALUES
	(
		:user_id,
		:card_id,
		:alias_id,
		:type,
		:created_at,
		:updated_at
	)
	ON CONFLICT (user_id, card_id) DO UPDATE
	SET
		alias_id      = EXCLUDED.alias_id,
		type          = EXCLUDED.type,
		updated_at    = EXCLUDED.updated_at
	`, m)

	return errors.Wrap(err, "SaveUserReaction failed")
}

func (store *Store) GetUserReaction(userID, cardID globalid.ID) (*model.UserReaction, error) {
	var result model.UserReaction
	err := store.Get(&result,
		`SELECT *
		 FROM user_reactions
		 WHERE card_id = $1
		 AND   user_id = $2`, cardID, userID)
	return &result, errors.Wrap(err, "GetUserReaction failed")
}

func (store *Store) GetUserReactions(userID globalid.ID, cardIDs []globalid.ID) (map[globalid.ID]*model.UserReaction, error) {
	result := make(map[globalid.ID]*model.UserReaction, len(cardIDs))
	var reactions []*model.UserReaction
	err := store.Select(&reactions, `SELECT * FROM user_reactions WHERE user_id = $2 AND card_id = ANY($1::uuid[])`, pq.Array(cardIDs), userID)
	if err != nil && errors.Cause(err) != sql.ErrNoRows {
		return nil, errors.Wrap(err, "GetUserReactions failed")
	}
	for _, reaction := range reactions {
		result[reaction.CardID] = reaction
	}
	return result, nil
}

func (store *Store) DeleteUserReaction(userID, cardID globalid.ID) error {
	_, err := store.Exec(`DELETE FROM user_reactions WHERE card_id = $1 AND user_id = $2`, cardID, userID)
	return errors.Wrap(err, "DeleteUserReaction failed")
}
func (store *Store) DeleteUserReactionForType(userID, cardID globalid.ID, typ model.UserReactionType) (int64, error) {
	result, err := store.Exec(`DELETE FROM user_reactions WHERE card_id = $1 AND user_id = $2 AND type = $3`, cardID, userID, typ)
	if err != nil {
		return 0, errors.Wrap(err, "DeleteUserReactionForType failed")
	}
	return result.RowsAffected()
}

type CardLike struct {
	ByName      string    `db:"by_name"`
	ImagePath   string    `db:"image_path"`
	IsAnonymous bool      `db:"is_anonymous"`
	Timestamp   time.Time `db:"created_at"`
}

func (cc *CardLike) DisplayName() string {
	r := cc.ByName
	if cc.IsAnonymous {
		r = "!" + r
	}
	return r
}

type LikeNotificationExportData struct {
	Boosts         []*CardLike `db:"-"`
	IsComment      bool        `db:"is_comment"`
	PosterImage    string      `db:"profile_image_path"`
	CardContent    string      `db:"content"`
	ThreadRootID   globalid.ID `db:"thread_root_id"`
	ThreadReplyID  globalid.ID `db:"thread_reply_id"`
	AuthorUsername string      `db:"username"`
}

func (store *Store) GetLikeNotificationExportData(n *model.Notification) (*LikeNotificationExportData, error) {
	var boosts []*CardLike

	var boostData LikeNotificationExportData

	// gets the latest boost action for any user that's ever boosted, total count is number of users that have ever boosted
	err := store.Select(&boosts, `
		SELECT
			COALESCE(anonymous_aliases.username, users.display_name) as "by_name",
			COALESCE(anonymous_aliases.profile_image_path, users.profile_image_path) as "image_path",
			r.alias_id IS NOT NULL as "is_anonymous",
			r.created_at
		FROM notifications_reactions
			LEFT JOIN user_reactions AS r ON notifications_reactions.user_id = r.user_id AND notifications_reactions.card_id = r.card_id
			LEFT JOIN users ON r.user_id = users.id
			LEFT JOIN anonymous_aliases ON r.alias_id = anonymous_aliases.id
		WHERE notifications_reactions.notification_id = $1
		ORDER BY r.created_at DESC`, n.ID)
	if err != nil {
		return nil, errors.Wrap(err, "GetLikeNotificationExportData failed")
	}

	err = store.Get(&boostData, `
		SELECT
			cards.thread_root_id IS NOT NULL as "is_comment",
			users.profile_image_path,
			cards.content,
			COALESCE(cards.thread_root_id, cards.id) as "thread_root_id",
			cards.thread_reply_id,
			COALESCE(anonymous_aliases.username, users.username) as "username"
		FROM cards
			LEFT JOIN users ON cards.owner_id = users.id
			LEFT JOIN anonymous_aliases ON cards.alias_id = anonymous_aliases.id
		WHERE cards.id = $1`, n.TargetID)

	if err != nil {
		return nil, errors.Wrap(err, "GetLikeNotificationExportData failed")
	}

	boostData.Boosts = boosts

	return &boostData, nil
}
