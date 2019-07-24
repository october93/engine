package datastore

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

// SaveCard saves a card
func (store *Store) SaveChannel(m *model.Channel) error {
	return saveChannel(store, m)
}

// SaveCard saves a card
func (tx *Tx) SaveChannel(m *model.Channel) error {
	return saveChannel(tx, m)
}

func saveChannel(e sqlx.Ext, m *model.Channel) error {
	if m == nil {
		return errors.New("provided model can not be nil")
	}
	if m.ID == globalid.Nil {
		m.ID = globalid.Next()
	}
	tn := time.Now().UTC()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = tn
	}
	m.UpdatedAt = tn
	_, err := sqlx.NamedExec(e,
		`INSERT INTO channels
	(
		id,
		name,
		handle,
		owner_id,
		description,
		is_default,
		private,
		created_at,
		updated_at
	)
	VALUES
	(
		:id,
		:name,
		:handle,
		:owner_id,
		:description,
		:is_default,
		:private,
		:created_at,
		:updated_at
	)
	ON CONFLICT(id) DO UPDATE
	SET

		name          = :name,
		handle        = :handle,
		owner_id      = :owner_id,
		description   = :description,
		private       = :private,
		is_default    = :is_default,
    created_at    = :created_at,
    updated_at    = :updated_at
	WHERE channels.id = :id `, m)
	if err != nil {
		return errors.Wrap(err, "SaveChannel failed")
	}
	return nil
}

// GetCard reads the given card from the database.
func (store *Store) GetChannel(id globalid.ID) (*model.Channel, error) {
	chann := model.Channel{}
	err := store.Get(&chann, "SELECT * FROM channels where id = $1", id)
	if err != nil {
		return nil, errors.Wrap(err, "GetChannel failed")
	}

	return &chann, nil
}

// GetCard reads the given card from the database.
func (store *Store) GetDefaultChannelIDs() ([]globalid.ID, error) {
	channs := []globalid.ID{}
	err := store.Select(&channs, "SELECT id FROM channels where is_default = true")
	if err != nil {
		return nil, errors.Wrap(err, "GetChannel failed")
	}

	return channs, nil
}

// GetCard reads the given card from the database.
func (store *Store) GetIsSubscribed(userID, channelID globalid.ID) (bool, error) {
	var ismember bool
	err := store.Get(&ismember, "SELECT true FROM channel_memberships where channel_id = $1 AND user_id = $2", channelID, userID)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, errors.Wrap(err, "GetIsSubscribed failed")
	}

	return true, nil
}

func (store *Store) IsSubscribedToChannels(userID globalid.ID, channelIDs []globalid.ID) (map[globalid.ID]bool, error) {
	var conditions []*model.Condition
	err := store.Select(&conditions, `SELECT   channel_id "id", true "condition"
	                                  FROM     channel_memberships
									  WHERE    channel_id = ANY($1::uuid[])
									  AND      user_id = $2
									  GROUP BY channel_id`, pq.Array(channelIDs), userID)
	if err != nil {
		return nil, errors.Wrap(err, "IsSubscribedToChannels failed")
	}
	result := make(map[globalid.ID]bool, len(conditions))
	for _, c := range conditions {
		result[c.ID] = c.Condition
	}
	return result, nil
}

// GetCard reads the given card from the database.
func (store *Store) GetSubscriberCount(channelID globalid.ID) (int, error) {
	var ismember int
	err := store.Get(&ismember, "SELECT count(*) FROM channel_memberships where channel_id = $1", channelID)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, errors.Wrap(err, "GetChannel failed")
	}

	return ismember, nil
}

// GetCard reads the given card from the database.
func (store *Store) GetChannelByHandle(handle string) (*model.Channel, error) {
	chann := model.Channel{}
	err := store.Get(&chann, "SELECT * FROM channels where handle = $1", handle)
	if err != nil {
		return nil, errors.Wrap(err, "GetChannel failed")
	}

	return &chann, nil
}

// GetCard reads the given card from the database.
func (store *Store) GetChannels() ([]*model.Channel, error) {
	channs := []*model.Channel{}
	err := store.Select(&channs, "SELECT * FROM channels")
	if err != nil {
		return nil, errors.Wrap(err, "GetChannel failed")
	}

	return channs, nil
}

func (store *Store) GetChannelsByID(ids []globalid.ID) ([]*model.Channel, error) {
	var channels []*model.Channel
	err := store.Select(&channels, `SELECT * FROM channels WHERE id = ANY($1::uuid[])`, pq.Array(ids))
	return channels, errors.Wrap(err, "GetChannelsByID failed")
}

// GetCard reads the given card from the database.
func (store *Store) GetChannelsForUser(userID globalid.ID) ([]*model.Channel, error) {
	channs := []*model.Channel{}
	err := store.Select(&channs, `
		SELECT channels.*
		FROM channels
		    LEFT JOIN (
		        SELECT channel_id, TRUE AS subscribed
		        FROM channel_memberships
		        WHERE user_id = $1
		    ) AS subs
		    ON channels.id = subs.channel_id
		    LEFT JOIN (
		        SELECT channel_id, count(*) AS member_count
		        FROM cards
		        GROUP BY channel_id
		    ) AS counts
		    ON channels.id = counts.channel_id
		WHERE (private = FALSE OR subscribed = TRUE)
		ORDER BY subscribed, coalesce(member_count, 0) DESC
	`, userID)
	if err != nil {
		return nil, errors.Wrap(err, "GetChannel failed")
	}

	return channs, nil
}

// GetCard reads the given card from the database.
func (store *Store) GetSubscribedChannels(userID globalid.ID) ([]*model.Channel, error) {
	channs := []*model.Channel{}
	err := store.Select(&channs, "SELECT channels.* FROM channel_memberships LEFT JOIN channels ON channel_memberships.channel_id = channels.id WHERE channel_memberships.user_id = $1", userID)
	if err != nil {
		return nil, errors.Wrap(err, "GetChannel failed")
	}

	return channs, nil
}

// GetCard reads the given card from the database.
func (store *Store) GetChannelInfos(userID globalid.ID) ([]*model.ChannelUserInfo, error) {
	channs := []*model.ChannelUserInfo{}
	err := store.Select(&channs, `
	SELECT channels.id as channel_id, subs.subscribed IS NOT NULL as subscribed, coalesce(member_count, 0) as member_count
	FROM channels
			LEFT JOIN (
					SELECT channel_id, TRUE AS subscribed
					FROM channel_memberships
					WHERE user_id = $1
			) AS subs
			ON channels.id = subs.channel_id
			LEFT JOIN (
					SELECT channel_id, count(*) AS member_count
					FROM channel_memberships
					GROUP BY channel_id
			) AS counts
			ON channels.id = counts.channel_id
	ORDER BY subscribed, coalesce(member_count, 0) DESC`, userID)
	if err != nil {
		return nil, errors.Wrap(err, "GetChannel failed")
	}

	return channs, nil
}

func (store *Store) JoinChannel(userID, channelID globalid.ID) error {
	_, err := store.Exec(`
		INSERT INTO channel_memberships (user_id, channel_id)
		VALUES ($1, $2)
		ON CONFLICT(user_id, channel_id)
		DO NOTHING`, userID, channelID)
	if err != nil {
		return errors.Wrap(err, "JoinChannel failed")
	}

	_, err = store.Exec(`
		DELETE FROM channel_mutes WHERE user_id = $1 AND channel_id = $2`, userID, channelID)
	if err != nil {
		return errors.Wrap(err, "JoinChannel failed")
	}

	return nil
}

func (store *Store) LeaveChannel(userID, channelID globalid.ID) error {
	_, err := store.Exec(`DELETE FROM channel_memberships WHERE user_id = $1 AND channel_id = $2`, userID, channelID)
	if err != nil {
		return errors.Wrap(err, "LeaveChannel failed")
	}

	_, err = store.Exec(`
		DELETE FROM channel_mutes WHERE user_id = $1 AND channel_id = $2`, userID, channelID)
	if err != nil {
		return errors.Wrap(err, "LeaveChannel failed")
	}

	return nil
}

func (store *Store) LeaveAllChannels(userID globalid.ID) error {
	_, err := store.Exec(`DELETE FROM channel_memberships WHERE user_id = $1`, userID)
	if err != nil {
		return errors.Wrap(err, "LeaveAllChannels failed")
	}
	return nil
}

func (store *Store) MuteChannel(userID, channelID globalid.ID) error {
	_, err := store.Exec(`
		INSERT INTO channel_mutes (user_id, channel_id)
		VALUES ($1, $2)
		ON CONFLICT(user_id, channel_id)
		DO NOTHING`, userID, channelID)
	if err != nil {
		return errors.Wrap(err, "MuteChannel failed")
	}

	_, err = store.Exec(`DELETE FROM channel_memberships WHERE user_id = $1 AND channel_id = $2`, userID, channelID)
	if err != nil {
		return errors.Wrap(err, "MuteChannel failed")
	}
	return nil
}

func (store *Store) UnmuteChannel(userID, channelID globalid.ID) error {
	_, err := store.Exec(`DELETE FROM channel_mutes WHERE user_id = $1 AND channel_id = $2`, userID, channelID)
	if err != nil {
		return errors.Wrap(err, "UnmuteChannel failed")
	}
	return nil
}

// AddUserToDefaultChannels stop
func (store *Store) AddUserToDefaultChannels(userID globalid.ID) error {
	_, err := store.Exec(`
		INSERT INTO channel_memberships (user_id, channel_id)
		SELECT $1, id FROM channels WHERE is_default = true
		ON CONFLICT(user_id, channel_id)
		DO NOTHING`, userID)
	if err != nil {
		return errors.Wrap(err, "AddUserToDefaultChannels failed")
	}
	return nil
}

// AddUserToDefaultChannels stop
func (store *Store) AddAllUsersToChannel(channelID globalid.ID) error {
	_, err := store.Exec(`
		INSERT INTO channel_memberships (user_id, channel_id)
		SELECT id, $1 FROM users
		ON CONFLICT(user_id, channel_id)
		DO NOTHING`, channelID)
	if err != nil {
		return errors.Wrap(err, "AddUserToDefaultChannels failed")
	}
	return nil
}

// GetCardsForChannel returns cards posted to a channel.
func (store *Store) GetCardsForChannel(channelID globalid.ID, count, skip int, forUser globalid.ID) ([]*model.Card, error) {
	cards := []*model.Card{}
	err := store.Select(&cards, `
		SELECT * FROM cards
		WHERE channel_id = $1
			AND deleted_at IS NULL
			AND (shadowbanned_at IS NULL OR owner_id = $4)
		ORDER BY created_at DESC
		LIMIT $2
		OFFSET $3`, channelID, count, count*skip, forUser)
	if err != nil {
		return nil, errors.Wrap(err, "GetCardsForChannel failed")
	}

	return cards, nil
}

// GetCardsForChannel returns cards posted to a channel.
func (store *Store) GetChannelEngagements() ([]*model.ChannelEngagement, error) {
	cards := []*model.ChannelEngagement{}
	err := store.Select(&cards, `
		SELECT
			channel_id,
			COUNT(*) as total_posts,
			COALESCE(SUM(commentcount), 0) AS total_comments,
			COALESCE(SUM(commenterscount), 0) AS total_commenters,
			COALESCE(SUM(likecount), 0) AS total_likes,
			COALESCE(SUM(dislikecount), 0) AS total_dislikes
		FROM
		cards
			LEFT JOIN
			(SELECT thread_root_id, count(*) AS commentcount FROM cards GROUP BY thread_root_id) AS comments ON cards.id = comments.thread_root_id
			LEFT JOIN
			(SELECT thread_root_id, COUNT(DISTINCT(COALESCE(alias_id, owner_id))) AS commenterscount
		  FROM cards WHERE thread_root_id IS NOT NULL
		  GROUP BY thread_root_id) AS commenters ON cards.id = commenters.thread_root_id
			LEFT JOIN
			(SELECT card_id, count(*) AS dislikecount FROM user_reactions WHERE type = 'dislike' GROUP BY card_id) AS dislikes ON cards.id = dislikes.card_id
			LEFT JOIN
			(SELECT card_id, count(*) AS likecount FROM user_reactions WHERE type = 'like' GROUP BY card_id) AS likes ON cards.id = likes.card_id
		WHERE channel_id IS NOT NULL AND cards.thread_root_id IS NULL
		GROUP BY channel_id
		`)
	if err != nil {
		return nil, errors.Wrap(err, "GetChannelEngagements failed")
	}

	return cards, nil
}
