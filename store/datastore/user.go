package datastore

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

var ErrUserNotFound = errors.New("user not found")

// SaveUser will update an existing user or create a new user record in the
// datastore.
func (store *Store) SaveUser(m *model.User) error {
	return saveUser(store, m)
}

// SaveUser will update an existing user or create a new user record in the
// datastore.
func (tx *Tx) SaveUser(m *model.User) error {
	return saveUser(tx, m)
}

func saveUser(e sqlx.Ext, m *model.User) error {
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
		`INSERT INTO users
		(
			id,
			display_name,
			first_name,
			last_name,
			profile_image_path,
			cover_image_path,
			bio,
			email,
			allow_email,
			password_hash,
			password_salt,
			username,
			devices,
			is_default,
			admin,
			search_key,
			disable_feed,
			feed_updated_at,
			joined_from_invite,
			shadowbanned_at,
			seen_intro_cards,
			coin_balance,
			is_verified,
			temporary_coin_balance,
			coin_reward_last_updated_at,
			blocked_at,
			created_at,
			updated_at,
			botched_signup,
			got_delayed_invites
		)
		VALUES
		(
			:id,
			:display_name,
			:first_name,
			:last_name,
			:profile_image_path,
			:cover_image_path,
			:bio,
			:email,
			:allow_email,
			:password_hash,
			:password_salt,
			:username,
			:devices,
			:is_default,
			:admin,
			:search_key,
			:disable_feed,
			:feed_updated_at,
			:joined_from_invite,
			:shadowbanned_at,
			:seen_intro_cards,
			:coin_balance,
			:is_verified,
			:temporary_coin_balance,
			:coin_reward_last_updated_at,
			:blocked_at,
			:created_at,
			:updated_at,
			:botched_signup,
			:got_delayed_invites
		)
		ON CONFLICT(id) DO UPDATE
		SET
			display_name        = :display_name,
			first_name          = :first_name,
			last_name           = :last_name,
			profile_image_path  = :profile_image_path,
			cover_image_path    = :cover_image_path,
			bio                 = :bio,
			email               = :email,
			allow_email         = :allow_email,
			password_hash       = :password_hash,
			password_salt       = :password_salt,
			username            = :username,
			devices             = :devices,
			is_default          = :is_default,
			admin               = :admin,
			search_key          = :search_key,
			disable_feed        =	:disable_feed,
			feed_updated_at     = :feed_updated_at,
			joined_from_invite  = :joined_from_invite,
			shadowbanned_at     = :shadowbanned_at,
			seen_intro_cards    = :seen_intro_cards,
			coin_balance        = :coin_balance,
			is_verified         = :is_verified,
			temporary_coin_balance = :temporary_coin_balance,
			coin_reward_last_updated_at = :coin_reward_last_updated_at,
			blocked_at          = :blocked_at,
			created_at          = :created_at,
			updated_at          = :updated_at,
			botched_signup      = :botched_signup,
			got_delayed_invites = :got_delayed_invites
		WHERE users.ID = :id`, m)
	return errors.Wrap(err, "SaveUser failed")
}

func (store *Store) BlockUser(blockingUser, blockedUser globalid.ID) error {
	_, err := store.Exec("INSERT INTO user_blocks (user_id, blocked_user, created_at) VALUES ($1, $2, $3) ON CONFLICT (user_id, blocked_user) DO NOTHING", blockingUser, blockedUser, time.Now().UTC())
	return err
}

func (store *Store) BlockAnonUserInThread(blockingUser, blockedAlias, threadID globalid.ID) error {
	_, err := store.Exec("INSERT INTO user_blocks (user_id, blocked_alias, for_thread, created_at) VALUES ($1, $2, COALESCE((SELECT thread_root_id FROM cards WHERE id = $3), $3), $4) ON CONFLICT (user_id, blocked_user) DO NOTHING", blockingUser, blockedAlias, threadID, time.Now().UTC())
	return err
}

func (store *Store) GetUser(id globalid.ID) (*model.User, error) {
	user := model.User{}
	err := store.Get(&user, "SELECT * FROM users where id = $1 AND deleted_at IS NULL", id)
	return &user, errors.Wrap(err, "GetUser failed")
}

func (store *Store) GetUserByEmail(email string) (*model.User, error) {
	user := model.User{}
	err := store.Get(&user, "SELECT * FROM users where email = $1 AND deleted_at IS NULL", email)
	if err != nil {
		return nil, errors.Wrap(err, "GetUserByEmail failed")
	}
	return &user, nil
}

func (store *Store) GetUsers() ([]*model.User, error) {
	users := []*model.User{}
	err := store.Select(&users, "SELECT * FROM users WHERE deleted_at IS NULL ORDER BY updated_at DESC")
	return users, errors.Wrap(err, "GetUsers failed")
}

func (store *Store) GetUserIDs() ([]globalid.ID, error) {
	var users []globalid.ID
	err := store.Select(&users, "SELECT id FROM users WHERE deleted_at IS NULL ORDER BY updated_at DESC")
	return users, errors.Wrap(err, "GetUserIDs failed")
}

func (store *Store) GetUsersByID(ids []globalid.ID) ([]*model.User, error) {
	query, args, err := sqlx.In(`SELECT users.* FROM unnest(ARRAY[?]::uuid[]) WITH ORDINALITY AS r(id, rn) LEFT OUTER JOIN users USING (id) ORDER BY r.rn;`, ids)
	if err != nil {
		return nil, errors.Wrap(err, "GetUsersByID failed")
	}
	query = store.Rebind(query)
	users := []*model.User{}
	err = store.Select(&users, query, args...)
	if len(users) != len(ids) {
		return nil, fmt.Errorf("GetUsersByID failed: unexpected number of users returned, likely requesting an invalid id\n")
	}
	return users, errors.Wrap(err, "GetUsersByID failed")
}

// GetUserByUsername reads a username identified by its username from the database.
func (store *Store) GetUserByUsername(username string) (*model.User, error) {
	user := model.User{}
	err := store.Get(&user, "SELECT * FROM users WHERE username = $1 AND deleted_at IS NULL", username)
	if err != nil {
		return nil, errors.Wrap(err, "GetUserByUsername failed")
	}
	return &user, nil
}

func (store *Store) GetUsersByUsernames(usernames []string) ([]*model.User, error) {
	var users []*model.User
	query, args, err := sqlx.In("SELECT * FROM users WHERE username IN (?) AND deleted_at IS NULL", usernames)
	if err != nil {
		return nil, errors.Wrap(err, "GetUsersByUsernames failed")
	}
	query = store.Rebind(query)
	err = store.Select(&users, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "GetUsersByUsernames failed")
	}
	return users, nil
}

func (store *Store) GetUserCount() (int, error) {
	var count int
	err := store.Get(&count, `SELECT COUNT(*) FROM users`)
	return count, errors.Wrap(err, "GetUserCount failed")
}

func (store *Store) DeleteUser(id globalid.ID) error {
	_, err := store.Exec(`UPDATE users SET deleted_at=now() WHERE id = $1`, id)
	return errors.Wrap(err, "DeleteUser failed")
}

func (store *Store) GetUsersWithoutDelayedInvites() ([]*model.User, error) {
	var result []*model.User
	err := store.Select(&result, `SELECT * FROM users WHERE got_delayed_invites = false AND blocked_at IS NULL AND deleted_at IS NULL`)
	return result, errors.Wrap(err, "GetUsersWithoutDelayedInvites")
}

func (store *Store) UpdateUserGotDelayedInvites(id globalid.ID) error {
	_, err := store.Exec(`UPDATE users SET got_delayed_invites = true WHERE id = $1`, id)
	return err
}

func (store *Store) MuteUser(userID, mutedUserID globalid.ID) error {
	_, err := store.Exec(`INSERT INTO user_mutes (user_id, muted_user_id) VALUES ($1, $2) ON CONFLICT (user_id, muted_user_id) DO NOTHING`, userID, mutedUserID)
	if err != nil {
		return errors.Wrap(err, "MuteUser q1 failed")
	}
	_, err = store.Exec(`DELETE FROM user_follows WHERE follower_id = $1 AND followee_id = $2`, userID, mutedUserID)
	if err != nil {
		return errors.Wrap(err, "MuteUser q2 failed")
	}
	return nil
}

func (store *Store) UnmuteUser(userID, mutedUserID globalid.ID) error {
	_, err := store.Exec(`DELETE FROM channel_mutes WHERE user_id = $1 AND muted_user_id = $2`, userID, mutedUserID)
	if err != nil {
		return errors.Wrap(err, "UnmuteUser failed")
	}
	return nil
}

func (store *Store) GetSafeUsersByPage(forUser globalid.ID, pageSize, pageNumber int, searchString string) ([]*model.User, error) {
	var rankings []*model.User
	var err error

	if searchString != "" {
		err = store.Select(&rankings, `
			SELECT users.*
			FROM users
				LEFT JOIN (SELECT * FROM user_follows WHERE follower_id = $1) as follows
				ON users.id = follows.followee_id
			WHERE id != $1
				AND shadowbanned_at IS NULL
				AND blocked_at IS NULL
				AND (
					users.username ILIKE $4
					OR
					users.display_name ILIKE $4
					OR
					users.bio ILIKE $4
				)
			ORDER BY followee_id, bio != '', id
			LIMIT $2
			OFFSET $3
		`, forUser, pageSize, pageSize*pageNumber, fmt.Sprintf(`%%%[1]s%%`, searchString))
	} else {
		err = store.Select(&rankings, `
			SELECT users.*
			FROM users
				LEFT JOIN (SELECT * FROM user_follows WHERE follower_id = $1) as follows
				ON users.id = follows.followee_id
			WHERE id != $1
				AND shadowbanned_at IS NULL
				AND blocked_at IS NULL
			ORDER BY followee_id, bio != '', id
			LIMIT $2
			OFFSET $3
		`, forUser, pageSize, pageSize*pageNumber)
	}

	return rankings, errors.Wrap(err, "GetSafeUsersByPage failed")
}

func (store *Store) ResetTemporaryCoins(newBalance int64) error {
	_, err := store.Exec(`UPDATE users SET temporary_coin_balance = $1`, newBalance)
	return errors.Wrap(err, "ResetTemporaryCoins failed")

}

func (store *Store) AwardCoins(userID globalid.ID, amount int64) error {
	_, err := store.Exec(`UPDATE users SET coin_balance = coin_balance + $1 WHERE id = $2`, amount, userID)
	return errors.Wrap(err, "UpdateCoinBalanceForUser failed")
}

func (store *Store) AwardTemporaryCoins(userID globalid.ID, amount int64) error {
	_, err := store.Exec(`UPDATE users SET temporary_coin_balance = temporary_coin_balance + $1 WHERE id = $2`, amount, userID)
	return errors.Wrap(err, "UpdateCoinBalanceForUser failed")
}

func (store *Store) PayCoinAmount(userID globalid.ID, amount int64) error {
	_, err := store.Exec(`
		UPDATE users
		SET
			temporary_coin_balance = GREATEST(temporary_coin_balance - $1, 0),
			coin_balance = LEAST(coin_balance + temporary_coin_balance - $1, coin_balance)
		WHERE id = $2`, amount, userID)
	return errors.Wrap(err, "UpdateCoinBalanceForUser failed")
}

func (store *Store) GetCurrentBalance(userID globalid.ID) (*model.CoinBalances, error) {
	var cB model.CoinBalances
	err := store.Get(&cB, `SELECT coin_balance, temporary_coin_balance FROM users WHERE id = $1`, userID)
	return &cB, errors.Wrap(err, "GetCurrentBalance failed")
}

func (store *Store) AddCoinsToBalance(userID globalid.ID, amount int) error {
	_, err := store.Exec(`UPDATE users SET coin_balance = coin_balance + $1 WHERE id = $2`, amount, userID)
	return errors.Wrap(err, "UpdateCoinBalanceForUser failed")
}

func (store *Store) SubtractCoinsFromBalance(userID globalid.ID, amount int) error {
	_, err := store.Exec(`
		UPDATE users
		SET
			coin_balance = coin_balance - $1
		WHERE id = $2`, amount, userID)

	return errors.Wrap(err, "DebitCoinsFromBalance failed")
}
