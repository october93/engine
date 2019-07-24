package datastore

import (
	"time"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

// SaveAnnouncement saves an Announcement to the store
func (store *Store) SaveAnnouncement(m *model.Announcement) error {
	return saveAnnouncement(store, m)
}

// SaveAnnouncement saves an Announcement to the store
func (tx *Tx) SaveAnnouncement(m *model.Announcement) error {
	return saveAnnouncement(tx, m)
}

func saveAnnouncement(e sqlx.Ext, m *model.Announcement) error {
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
		`INSERT INTO announcements
		(
			id,
			from_user,
			card_id,
			message,
			created_at,
			updated_at,
			deleted_at
		)
		VALUES
		(
			:id,
			:from_user,
			:card_id,
			:message,
			:created_at,
			:updated_at,
			:deleted_at
		)
		ON CONFLICT(id) DO UPDATE
		SET
			id             = :id,
			from_user      = :from_user,
			card_id        = :card_id,
			message        = :message,
			created_at     = :created_at,
			updated_at     = :updated_at,
			deleted_at     = :deleted_at
		WHERE announcements.id = :id`, m)
	return errors.Wrap(err, "SaveAnnouncement failed")
}

func (store *Store) GetAnnouncementForNotification(n *model.Notification) (*model.Announcement, error) {
	var a model.Announcement
	err := store.Get(&a, `SELECT * FROM announcements WHERE id = $1 AND deleted_at IS NULL LIMIT 1`, n.TargetID)
	if err != nil {
		return nil, errors.Wrap(err, "getAnnouncementForNotification failed")
	}
	return &a, nil
}

func (store *Store) DeleteAnnouncement(id globalid.ID) error {
	_, err := store.Exec(`UPDATE announcements SET deleted_at=now() WHERE id = $1`, id)
	if err != nil {
		return errors.Wrap(err, "DeleteAnnouncement failed")
	}
	_, err = store.Exec(`UPDATE notifications SET deleted_at=now() WHERE target_id = $1`, id)
	return err
}

func (store *Store) GetAnnouncements() ([]*model.Announcement, error) {
	a := []*model.Announcement{}
	err := store.Select(&a, `SELECT * FROM announcements`)
	return a, errors.Wrap(err, "GetAnnouncements failed")
}
