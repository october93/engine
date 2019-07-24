package datastore

import (
	"time"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

// SaveNotification saves an notification to the store
func (store *Store) SaveMention(m *model.Mention) error {
	return saveMention(store, m)
}

// SaveNotification saves an notification to the store
func (tx *Tx) SaveMention(m *model.Mention) error {
	return saveMention(tx, m)
}

func saveMention(e sqlx.Ext, m *model.Mention) error {
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
		`INSERT INTO mentions
		(
			id,
			in_card,
			mentioned_user,
			mentioned_alias,
			created_at,
			updated_at,
			deleted_at
		)
		VALUES
		(
			:id,
			:in_card,
			:mentioned_user,
      :mentioned_alias,
			:created_at,
			:updated_at,
			:deleted_at
		)
		ON CONFLICT(id) DO UPDATE
		SET
			id                = :id,
			in_card           = :in_card,
			mentioned_user    = :mentioned_user,
      mentioned_alias   = :mentioned_alias,
			created_at        = :created_at,
			updated_at        = :updated_at,
			deleted_at        = :deleted_at
		WHERE mentions.id = :id`, m)
	return errors.Wrap(err, "SaveMention failed")
}

// SaveNotification saves an notification to the store
func (store *Store) SaveNotificationMention(m *model.NotificationMention) error {
	return saveNotificationMention(store, m)
}

// SaveNotification saves an notification to the store
func (tx *Tx) SaveNotificationMention(m *model.NotificationMention) error {
	return saveNotificationMention(tx, m)
}

func saveNotificationMention(e sqlx.Ext, m *model.NotificationMention) error {
	if m == nil {
		return errors.New("provided model can not be nil")
	}
	tn := time.Now().UTC()

	if m.CreatedAt.IsZero() {
		m.CreatedAt = tn
	}
	m.UpdatedAt = tn

	_, err := sqlx.NamedExec(e,
		`INSERT INTO notifications_mentions
		(
			notification_id,
			mention_id,
			created_at,
			updated_at
		)
		VALUES
		(
			:notification_id,
			:mention_id,
			:created_at,
			:updated_at
		)
		ON CONFLICT(notification_id, mention_id) DO UPDATE
		SET
			notification_id   = :notification_id,
			mention_id        = :mention_id,
			created_at        = :created_at,
			updated_at        = :updated_at
		WHERE notifications_mentions.notification_id = :notification_id AND notifications_mentions.mention_id = :mention_id`, m)
	return errors.Wrap(err, "SaveNotificationMention failed")
}

type MentionNotificationExportData struct {
	Name                 string      `db:"name"`
	IsAnonymous          bool        `db:"is_anon"`
	ImagePath            string      `db:"image_path"`
	InComment            bool        `db:"in_comment"`
	ThreadRoot           globalid.ID `db:"thread_root_id"`
	ThreadReply          globalid.ID `db:"thread_reply_id"`
	InCard               globalid.ID `db:"in_card"`
	InCardAuthorUsername string      `db:"username"`
}

func (store *Store) DeleteMentionsForCard(cardID globalid.ID) error {
	_, err := store.Exec(`UPDATE mentions SET deleted_at=now() WHERE in_card = $1`, cardID)
	return errors.Wrap(err, "DeleteMentionsForCard failed")
}

func (store *Store) GetMentionExportData(n *model.Notification) (*MentionNotificationExportData, error) {
	var data MentionNotificationExportData
	err := store.Get(&data, `
		SELECT
			COALESCE(anonymous_aliases.username, users.display_name) as "name",
			cards.alias_id IS NOT NULL as "is_anon",
			COALESCE(anonymous_aliases.profile_image_path, users.profile_image_path) as "image_path",
			cards.thread_root_id IS NOT NULL as "in_comment",
			COALESCE(cards.thread_root_id, mentions.in_card) as "thread_root_id",
			cards.thread_reply_id as "thread_reply_id",
			mentions.in_card,
			COALESCE(anonymous_aliases.username, users.username) as "username"
		FROM notifications_mentions
			LEFT JOIN mentions ON notifications_mentions.mention_id = mentions.id
			LEFT JOIN cards ON mentions.in_card = cards.id
			LEFT JOIN anonymous_aliases ON cards.alias_id = anonymous_aliases.id
			LEFT JOIN users ON cards.owner_id = users.id
		WHERE notifications_mentions.notification_id = $1
			AND mentions.deleted_at IS NULL
		`, n.ID)
	if err != nil {
		return nil, errors.Wrap(err, "GetMentionExportData failed")
	}
	return &data, nil
}
