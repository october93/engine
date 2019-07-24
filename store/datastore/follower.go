package datastore

import (
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

func (store *Store) SaveFollower(followerID, followeeID globalid.ID) error {
	return saveFollower(store, followerID, followeeID)
}

func (tx *Tx) SaveFollower(followerID, followeeID globalid.ID) error {
	return saveFollower(tx, followerID, followeeID)
}

func saveFollower(e sqlx.Ext, followerID, followeeID globalid.ID) error {
	if followerID == globalid.Nil || followeeID == globalid.Nil {
		return errors.New("provided model cannot be nil")
	}
	_, err := e.Exec(
		`INSERT INTO user_follows
		(
			follower_id,
			followee_id,
			created_at
		)
		VALUES
		(
			$1,
			$2,
			$3
		)
		ON CONFLICT (follower_id, followee_id) DO NOTHING`, followerID, followeeID, time.Now().UTC())
	return errors.Wrap(err, "SaveFollower failed")
}

func (store *Store) FollowDefaultUsers(userID globalid.ID) error {
	_, err := store.Exec(
		`INSERT INTO user_follows
		(
			follower_id,
			followee_id,
			created_at
		)
		SELECT $1, id, $2 FROM users WHERE is_default = true
		ON CONFLICT (follower_id, followee_id) DO NOTHING`, userID, time.Now().UTC())
	return errors.Wrap(err, "FollowDefaultUsers failed")
}

func (store *Store) DeleteFollower(followerID, followeeID globalid.ID) error {
	_, err := store.Exec("DELETE FROM user_follows WHERE follower_id = $1 AND followee_id = $2", followerID, followeeID)
	return errors.Wrap(err, "DeleteFollower failed")
}

func (store *Store) GetFollowing(userID globalid.ID) ([]*model.User, error) {
	var following []*model.User
	err := store.Select(&following, "SELECT * FROM users WHERE id IN (SELECT followee_id FROM user_follows WHERE follower_id = $1)", userID)
	return following, errors.Wrap(err, "GetFollowing failed")
}

// SaveInvite saves an invite to the store
func (store *Store) SaveNotificationFollow(notifID, followerID, followeeID globalid.ID) error {
	return saveNotificationFollow(store, notifID, followerID, followeeID)
}

// SaveInvite saves an invite to the store
func (tx *Tx) SaveNotificationFollow(notifID, followerID, followeeID globalid.ID) error {
	return saveNotificationFollow(tx, notifID, followerID, followeeID)
}

func saveNotificationFollow(e sqlx.Ext, notifID, followerID, followeeID globalid.ID) error {
	if notifID == globalid.Nil || followeeID == globalid.Nil {
		return errors.New("provided model cannot be nil")
	}
	_, err := e.Exec(
		`INSERT INTO notifications_follows
		(
			notification_id,
			follower_id,
			followee_id,
			created_at
		)
		VALUES
		(
			$1,
			$2,
			$3,
			$4
		)
		ON CONFLICT (follower_id, followee_id) DO UPDATE
		SET
			notification_id  = $1,
			created_at       = $4
		WHERE notifications_follows.follower_id = $2 AND notifications_follows.followee_id = $3`, notifID, followerID, followeeID, time.Now().UTC())
	return errors.Wrap(err, "SaveNotificationFollow failed")
}

type Follower struct {
	ID        globalid.ID `db:"id"`
	Username  string      `db:"username"`
	Name      string      `db:"display_name"`
	ImagePath string      `db:"profile_image_path"`
	Timestamp time.Time   `db:"created_at"`
}

func (f *Follower) DisplayName() string {
	return f.Name
}

type FollowNotificationExportData struct {
	Followers []*Follower `db:"-"`
}

func (store *Store) IsFollowing(followerID, followeeID globalid.ID) (bool, error) {
	var rowCount int

	err := store.Get(&rowCount, `
			SELECT
				COUNT(*)
			FROM user_follows
			WHERE follower_id = $1
				AND followee_id = $2
			`, followerID, followeeID)
	if err != nil {
		return false, errors.Wrap(err, "GetFollowExportData failed")
	}

	return rowCount > 0, nil
}

func (store *Store) IsFollowings(followerID globalid.ID, followeeID []globalid.ID) (map[globalid.ID]bool, error) {
	var conditions []*model.Condition
	err := store.Select(&conditions, `SELECT   follower_id "id", COUNT(*) > 0 "condition"
	                                  FROM     user_follows
									  WHERE    follower_id = $1
									  AND      followee_id = ANY($2::uuid[])
									  GROUP BY follower_id`, followerID, pq.Array(followeeID))
	if err != nil {
		return nil, errors.Wrap(err, "GetThreadCounts failed")
	}
	result := make(map[globalid.ID]bool, len(conditions))
	for _, c := range conditions {
		result[c.ID] = c.Condition
	}
	return result, nil
}

func (store *Store) GetFollowExportData(n *model.Notification) (*FollowNotificationExportData, error) {
	var followers []*Follower

	var data FollowNotificationExportData

	err := store.Select(&followers, `
			SELECT
				users.id,
				users.username,
				users.display_name,
				users.profile_image_path,
				notifications_follows.created_at
			FROM notifications_follows
				LEFT JOIN users ON notifications_follows.follower_id = users.id
			WHERE notifications_follows.notification_id = $1
			ORDER BY created_at DESC`, n.ID)
	if err != nil {
		return nil, errors.Wrap(err, "GetFollowExportData failed")
	}

	data.Followers = followers

	return &data, nil
}
