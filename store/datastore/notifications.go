package datastore

import (
	"time"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

// SaveNotification saves an notification to the store
func (store *Store) SaveNotification(m *model.Notification) error {
	return saveNotification(store, m)
}

// SaveNotification saves an notification to the store
func (tx *Tx) SaveNotification(m *model.Notification) error {
	return saveNotification(tx, m)
}

func saveNotification(e sqlx.Ext, m *model.Notification) error {
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
		`INSERT INTO notifications
		(
			id,
			user_id,
			target_id,
			target_alias_id,
			type,
			seen_at,
			opened_at,
			created_at,
			updated_at,
			deleted_at
		)
		VALUES
		(
			:id,
			:user_id,
			:target_id,
			:target_alias_id,
			:type,
			:seen_at,
			:opened_at,
			:created_at,
			:updated_at,
			:deleted_at
		)
		ON CONFLICT(id) DO UPDATE
		SET
			id              = :id,
			user_id         = :user_id,
			target_id       = :target_id,
			target_alias_id = :target_alias_id,
			type            = :type,
			seen_at         = :seen_at,
			opened_at       = :opened_at,
			created_at      = :created_at,
			updated_at      = :updated_at,
			deleted_at      = :deleted_at
		WHERE notifications.id = :id`, m)
	return errors.Wrap(err, "SaveNotification failed")
}

// GetNotifications reads all notifications for a given user from the database.
func (store *Store) GetNotifications(userID globalid.ID, page, skip int) ([]*model.Notification, error) {
	notifications := []*model.Notification{}
	err := store.Select(&notifications, `SELECT * FROM notifications WHERE user_id=$1 AND deleted_at IS NULL ORDER BY updated_at DESC OFFSET $2 LIMIT $3`, userID, page*skip, page)
	if err != nil {
		return nil, err
	}

	return notifications, errors.Wrap(err, "GetNotifications failed")
}

// GetNotifications reads all notifications for a given user from the database.
func (store *Store) GetNotification(notifID globalid.ID) (*model.Notification, error) {
	var notification model.Notification
	err := store.Get(&notification, `SELECT * FROM notifications WHERE id = $1 AND deleted_at IS NULL`, notifID)
	if err != nil {
		return nil, errors.Wrap(err, "GetNotification failed")
	}

	return &notification, nil
}

// GetNotifications reads all notifications for a given user from the database.
func (store *Store) UnseenNotificationsCount(userID globalid.ID) (int, error) {
	count := 0
	err := store.Get(&count, `SELECT COUNT(*) FROM notifications WHERE user_id=$1 AND deleted_at IS NULL AND seen_at IS NULL`, userID)
	if err != nil {
		return 0, errors.Wrap(err, "UnseenNotificationsCount failed")
	}

	return count, nil
}

// GetNotifications reads all notifications for a given user from the database.
func (store *Store) LatestForType(userID, targetID globalid.ID, typ string, unopenedOnly bool) (*model.Notification, error) {
	notif := model.Notification{}
	var err error

	if targetID != globalid.Nil {
		if unopenedOnly {
			err = store.Get(&notif, `SELECT * FROM notifications WHERE user_id=$1 AND type = $2 AND target_id = $3 AND opened_at IS NULL AND deleted_at IS NULL`, userID, typ, targetID)
		} else {
			err = store.Get(&notif, `SELECT * FROM notifications WHERE user_id=$1 AND type = $2 AND target_id = $3 AND deleted_at IS NULL`, userID, typ, targetID)
		}
	} else {
		if unopenedOnly {
			err = store.Get(&notif, `SELECT * FROM notifications WHERE user_id=$1 AND type = $2 AND opened_at IS NULL AND deleted_at IS NULL`, userID, typ)
		} else {
			err = store.Get(&notif, `SELECT * FROM notifications WHERE user_id=$1 AND type = $2 AND deleted_at IS NULL`, userID, typ)
		}
	}

	if err != nil {
		return nil, errors.Wrap(err, "LatestForType failed")
	}

	return &notif, nil
}

// UpdateNotifcationsSeen marks notifications as seen by setting the last seen at column to now.
func (store *Store) UpdateNotificationsSeen(ids []globalid.ID) error {
	if len(ids) == 0 {
		return nil
	}
	query, args, err := sqlx.In(`UPDATE notifications SET seen_at=now() WHERE id IN (?)`, ids)
	if err != nil {
		return errors.Wrap(err, "UpdateNotificationsSeen failed")
	}
	query = store.Rebind(query)
	var nothing []interface{}
	err = store.Select(&nothing, query, args...)
	return errors.Wrap(err, "UpdateNotificationsSeen failed")
}

// UpdateNotifcationsSeen marks notifications as seen by setting the last seen at column to now.
func (store *Store) UpdateAllNotificationsSeen(userID globalid.ID) error {
	_, err := store.Exec(`UPDATE notifications SET seen_at=now() WHERE user_id = $1`, userID)
	return errors.Wrap(err, "UpdateAllNotificationsSeen failed")
}

// UpdateNotifcationsOpened marks notifications as opened by setting the last opened at column to now.
func (store *Store) UpdateNotificationsOpened(ids []globalid.ID) error {
	if len(ids) == 0 {
		return nil
	}
	query, args, err := sqlx.In(`UPDATE notifications SET opened_at=(?) WHERE id IN (?)`, time.Now(), ids)
	if err != nil {
		return errors.Wrap(err, "UpdateNotificationsOpened failed")
	}
	query = store.Rebind(query)
	var nothing []interface{}
	err = store.Select(&nothing, query, args...)
	return errors.Wrap(err, "UpdateNotificationsOpened failed")
}

func (store *Store) DeleteNotification(id globalid.ID) error {
	_, err := store.Exec(`UPDATE notifications SET deleted_at=now() WHERE id = $1`, id)
	return errors.Wrap(err, "DeleteNotification failed")
}

func (store *Store) ClearEmptyNotifications() error {
	// delete ANY empty notifications
	_, err := store.Exec(`
		UPDATE notifications SET deleted_at=now()
		WHERE type IN ('boost', 'comment', 'mention', 'follow')
			AND id NOT IN (
			SELECT DISTINCT notification_id FROM notifications_follows
				UNION
			SELECT DISTINCT notification_id FROM notifications_mentions
				UNION
			SELECT DISTINCT notification_id FROM notifications_comments
				UNION
			SELECT DISTINCT notification_id FROM notifications_reactions
		)`)
	return errors.Wrap(err, "ClearEmptyNotifications failed")
}

func (store *Store) ClearEmptyNotificationsForUser(userID globalid.ID) error {
	// delete ANY empty notifications
	_, err := store.Exec(`
		UPDATE notifications SET deleted_at=now()
		WHERE type IN ('boost', 'comment', 'mention')
			AND user_id = $1
			AND id NOT IN (
			SELECT DISTINCT notification_id FROM notifications_mentions
				UNION
			SELECT DISTINCT notification_id FROM notifications_comments
				UNION
			SELECT DISTINCT notification_id FROM notifications_reactions
		)`, userID)
	return errors.Wrap(err, "ClearEmptyNotifications failed")
}

func (store *Store) DeleteNotificationsForCard(cardID globalid.ID) error {
	// delete mention actions
	_, err := store.Exec(`DELETE FROM notifications_mentions WHERE mention_id IN (SELECT id FROM mentions WHERE in_card = $1)`, cardID)
	if err != nil {
		return err
	}

	_, err = store.Exec(`DELETE FROM notifications_comments WHERE card_id = $1`, cardID)
	if err != nil {
		return err
	}

	return store.ClearEmptyNotifications()
}

type LeaderboardNotificationExportData struct {
	NotificationID globalid.ID `db:"notification_id"`
	Rank           int64       `db:"rank"`
}

func (store *Store) GetLeaderboardNotificationExportData(notifID globalid.ID) (*LeaderboardNotificationExportData, error) {
	ret := LeaderboardNotificationExportData{}
	err := store.Get(&ret, `SELECT notification_id, rank FROM notifications_leaderboard_data WHERE notification_id = $1`, notifID)
	if err != nil {
		return nil, errors.Wrap(err, "GetLeaderboardNotificationExportData failed")
	}

	return &ret, nil
}

func (store *Store) SaveLeaderboardNotificationData(notifID globalid.ID, rank int) error {
	_, err := store.Exec(
		`INSERT INTO notifications_leaderboard_data
		(
			notification_id,
			rank,
			created_at
		)
		VALUES
		(
			$1,
			$2,
			$3
		)
		ON CONFLICT(notification_id) DO UPDATE
		SET
			rank       = EXCLUDED.rank
		`, notifID, rank, time.Now().UTC())
	return errors.Wrap(err, "SaveReactionForNotification failed")
}

func (store *Store) SaveReactionForNotification(notifID, userID, cardID globalid.ID) error {
	_, err := store.Exec(
		`INSERT INTO notifications_reactions
		(
			notification_id,
			user_id,
			card_id,
			created_at,
			updated_at
		)
		VALUES
		(
			$1,
			$2,
			$3,
			$4,
			$4
		)
		ON CONFLICT(notification_id,user_id,card_id) DO UPDATE
		SET
			created_at       = EXCLUDED.created_at,
			updated_at       = EXCLUDED.updated_at
		`, notifID, userID, cardID, time.Now().UTC())
	return errors.Wrap(err, "SaveReactionForNotification failed")
}

func (store *Store) DeleteReactionForNotification(notifID, userID, cardID globalid.ID) error {
	_, err := store.Exec(
		`DELETE FROM notifications_reactions WHERE notification_id = $1 AND user_id = $2 AND card_id = $3
		`, notifID, userID, cardID)
	return errors.Wrap(err, "SaveReactionForNotification failed")
}
