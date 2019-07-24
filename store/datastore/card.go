package datastore

import (
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

// SaveCard saves a card
func (store *Store) SaveCard(m *model.Card) error {
	return saveCard(store, m)
}

// SaveCard saves a card
func (tx *Tx) SaveCard(m *model.Card) error {
	return saveCard(tx, m)
}

func saveCard(e sqlx.Ext, m *model.Card) error {
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
		`INSERT INTO cards
	(
		id,
		owner_id,
		thread_reply_id,
		thread_root_id,
		channel_id,
		title,
		content,
		url,
		background_image_path,
		background_color,
		anonymous,
		coins_earned,
		deleted_at,
		author_to_alias,
		alias_id,
		is_intro_card,
		shadowbanned_at,
		created_at,
		updated_at
	)
	VALUES
	(
		:id,
		:owner_id,
		:thread_reply_id,
		:thread_root_id,
		:channel_id,
		:title,
		:content,
		:url,
		:background_image_path,
		:background_color,
		:anonymous,
		:coins_earned,
		:deleted_at,
		:author_to_alias,
		:alias_id,
		:is_intro_card,
		:shadowbanned_at,
		:created_at,
		:updated_at
	)
	ON CONFLICT(id) DO UPDATE
	SET
    owner_id              = :owner_id,
    thread_reply_id       = :thread_reply_id,
    thread_root_id        = :thread_root_id,
		channel_id            = :channel_id,
    title                 = :title,
    content               = :content,
    url                   = :url,
	  background_image_path = :background_image_path,
	  background_color      = :background_color,
	  anonymous             = :anonymous,
		coins_earned          = :coins_earned,
	  deleted_at            = :deleted_at,
	  author_to_alias       = :author_to_alias,
  	alias_id              = :alias_id,
		is_intro_card         = :is_intro_card,
  	shadowbanned_at       = :shadowbanned_at,
    created_at            = :created_at,
    updated_at            = :updated_at
	WHERE cards.id = :id `, m)
	if err != nil {
		return errors.Wrap(err, "SaveCard failed")
	}
	return nil
}

// GetCard reads the given card from the database.
func (store *Store) GetCard(id globalid.ID) (*model.Card, error) {
	card := model.Card{}
	err := store.Get(&card, "SELECT * FROM cards where id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		return nil, errors.Wrap(err, "GetCard failed")
	}

	return &card, nil
}

// GetCards reads all cards from the database.
func (store *Store) GetCards() ([]*model.Card, error) {
	cards := []*model.Card{}
	err := store.Select(&cards, "SELECT * FROM cards WHERE deleted_at IS NULL ORDER BY updated_at DESC")
	if err != nil {
		return nil, errors.Wrap(err, "GetCards failed")
	}

	return cards, nil
}

// GetThread returns all cards with who are replying to the same card. This is
// used to display a thread on the application and excludes the thread root.
func (store *Store) GetThread(id, forUser globalid.ID) ([]*model.Card, error) {
	cards := []*model.Card{}
	err := store.Select(&cards,
		`WITH RECURSIVE
		 tree_search (thread_level) AS (
		     SELECT
		         0,
		         *
		     FROM cards
		     WHERE thread_root_id is NULL
		     UNION ALL
		     SELECT
		         ts.thread_level + 1,
		         cards.*
		     FROM cards, tree_search ts
		     WHERE cards.thread_reply_id = ts.id
		 )
		 SELECT * FROM tree_search
		 WHERE thread_level > 0
		 	AND thread_root_id = $1
			AND deleted_at IS NULL
			AND (shadowbanned_at IS NULL OR owner_id = $2)
			AND (owner_id NOT IN (SELECT blocked_user FROM user_blocks WHERE user_id = $2 AND blocked_user IS NOT NULL) OR alias_id IS NOT NULL)
			AND (alias_id NOT IN (SELECT blocked_alias FROM user_blocks WHERE user_id = $2 AND for_thread = $1) OR alias_id IS NULL)
		ORDER BY created_at`, id, forUser)

	if err != nil {
		return nil, errors.Wrap(err, "GetThread failed")
	}

	return cards, nil
}

// GetThread returns all cards with who are replying to the same card. This is
// used to display a thread on the application and excludes the thread root.
func (store *Store) GetFlatReplies(id, forUser globalid.ID, latestFirst bool, limitTo int) ([]*model.Card, error) {
	cards := []*model.Card{}
	queryTail := ""

	if latestFirst {
		queryTail += " DESC"
	}

	if limitTo > 0 {
		queryTail += fmt.Sprintf(" LIMIT %v", limitTo)
	}

	query := fmt.Sprintf(`
		SELECT * FROM cards
		WHERE id IN (
			WITH RECURSIVE sub_replies(id, thread_reply_id) AS (
		  	SELECT id, thread_reply_id FROM cards WHERE thread_reply_id = $1
				UNION ALL
				SELECT c.id, c.thread_reply_id
				FROM sub_replies repl, cards c
				WHERE c.thread_reply_id = repl.id
		  ) SELECT id FROM sub_replies
		)
		AND deleted_at IS NULL
		AND (shadowbanned_at IS NULL OR owner_id = $2)
		AND (owner_id NOT IN (SELECT blocked_user FROM user_blocks WHERE user_id = $2 AND blocked_user IS NOT NULL) OR alias_id IS NOT NULL)
		AND (alias_id NOT IN (SELECT blocked_alias FROM user_blocks WHERE user_id = $2 AND for_thread = $1) OR alias_id IS NULL)
		ORDER BY created_at%v`, queryTail)

	err := store.Select(&cards, query, id, forUser)

	if err != nil {
		return nil, errors.Wrap(err, "GetThread failed")
	}

	return cards, nil
}

// GetThread returns all cards with who are replying to the same card. This is
// used to display a thread on the application and excludes the thread root.
func (store *Store) GetRankedImmediateReplies(cardID, forUser globalid.ID) ([]*model.Card, error) {
	cards := []*model.Card{}
	err := store.Select(&cards,
		`SELECT * FROM cards
		 WHERE thread_reply_id = $1
			AND deleted_at IS NULL
			AND (shadowbanned_at IS NULL OR owner_id = $2)
			AND (owner_id NOT IN (SELECT blocked_user FROM user_blocks WHERE user_id = $2 AND blocked_user IS NOT NULL) OR alias_id IS NOT NULL)
			AND (alias_id NOT IN (SELECT blocked_alias FROM user_blocks WHERE user_id = $2 AND for_thread = $1 AND blocked_alias IS NOT NULL) OR alias_id IS NULL)
		ORDER BY created_at`, cardID, forUser)

	if err != nil {
		return nil, errors.Wrap(err, "GetThread failed")
	}

	return cards, nil
}

func (store *Store) GetThreadCount(id globalid.ID) (int, error) {
	var count int
	err := store.Get(&count, `SELECT COUNT(*) FROM cards WHERE thread_root_id = $1 AND deleted_at IS NULL`, id)
	return count, errors.Wrap(err, "GetThreadCount failed")
}

func (store *Store) GetThreadCounts(ids []globalid.ID) (map[globalid.ID]int, error) {
	var counts []*model.Count
	err := store.Select(&counts, `SELECT   thread_root_id "user_id", COUNT(*)
	                              FROM     cards
								  WHERE    thread_root_id = ANY($1::uuid[])
								  AND      deleted_at IS NULL
								  GROUP BY thread_root_id`, pq.Array(ids))
	if err != nil {
		return nil, errors.Wrap(err, "GetThreadCounts failed")
	}
	result := make(map[globalid.ID]int, len(counts))
	for _, count := range counts {
		result[count.UserID] = count.Count
	}
	return result, nil
}

func (store *Store) GetCardsByIDs(ids []globalid.ID) ([]*model.Card, error) {
	cards := []*model.Card{}

	query, args, err := sqlx.In(`SELECT cards.* FROM unnest(ARRAY[?]::uuid[]) WITH ORDINALITY AS r(id, rn) JOIN cards USING (id) ORDER BY r.rn`, ids)
	if err != nil {
		return nil, err
	}
	query = store.Rebind(query)
	err = store.Select(&cards, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "GetCardsByIDs failed")
	}

	return cards, nil
}

func (store *Store) GetCardsWithDeletedByID(ids []globalid.ID) ([]*model.Card, error) {
	cards := []*model.Card{}
	if len(ids) == 0 {
		return cards, nil
	}
	query, args, err := sqlx.In(`SELECT * FROM cards WHERE id IN (?) ORDER BY array_positions(ARRAY[?], CAST(id AS text))`, ids, ids)
	if err != nil {
		return nil, err
	}
	query = store.Rebind(query)
	err = store.Select(&cards, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "GetCardsWithDeletedByID failed")
	}

	return cards, nil
}

// GetCardsByInterval returns cards posted in a certain time period.
func (store *Store) GetCardsByInterval(from, to time.Time) ([]*model.Card, error) {
	cards := []*model.Card{}
	err := store.Select(&cards, "SELECT * FROM cards where created_at >= $1 and created_at <= $2 ORDER BY updated_at DESC", from, to)
	return cards, errors.Wrap(err, "GetCardsByInterval failed")
}

// GetCardsByNodeInRange returns cards posted by a user in a certain time period.
func (store *Store) GetCardsByNodeInRange(nodeID globalid.ID, from, to time.Time) ([]*model.Card, error) {
	cards := []*model.Card{}
	err := store.Select(&cards, "SELECT * FROM cards where owner_id = $1 and created_at >= $2 and created_at <= $3 ORDER BY updated_at DESC", nodeID, from, to)
	return cards, errors.Wrap(err, "GetCardsByNodeInRange failed")
}

func (store *Store) SetIntroCardStatus(cardID globalid.ID, status bool) error {
	_, err := store.Exec(`UPDATE cards SET is_intro_card=$2 WHERE id = $1`, cardID, status)
	if err != nil {
		return errors.Wrap(err, "SetIntroCardStatus failed")
	}
	return nil
}

// GetCards reads all cards from the database.
func (store *Store) GetIntroCards() ([]*model.Card, error) {
	cards := []*model.Card{}
	err := store.Select(&cards, "SELECT * FROM cards WHERE deleted_at IS NULL AND is_intro_card = true")
	if err != nil {
		return nil, errors.Wrap(err, "GetIntroCards failed")
	}

	return cards, nil
}

// GetCards reads all cards from the database.
func (store *Store) GetIntroCardIDs() ([]globalid.ID, error) {
	cards := []globalid.ID{}
	err := store.Select(&cards, "SELECT id FROM cards WHERE deleted_at IS NULL AND is_intro_card = true ORDER BY created_at")
	if err != nil {
		return nil, errors.Wrap(err, "GetIntroCards failed")
	}

	return cards, nil
}

// GetPostedCardsForNode returns the number of cards posted by a user.
func (store *Store) GetPostedCardsForNode(nodeID globalid.ID, count, skip int) ([]*model.Card, error) {
	stmt := `SELECT * FROM cards where owner_id = $1 AND thread_root_id IS NULL AND deleted_at IS NULL AND alias_id IS NULL ORDER BY created_at DESC`
	if count > 0 {
		stmt = fmt.Sprintf("%s\nLIMIT %d", stmt, count)
		if skip > 0 {
			stmt = fmt.Sprintf("%s\nOFFSET %d", stmt, count*skip)
		}
	}
	cards := []*model.Card{}
	err := store.Select(&cards, stmt, nodeID)
	if err != nil {
		return nil, errors.Wrap(err, "GetPostedCardsForNode failed")
	}

	return cards, nil
}

// GetPostedCardsForNodeIncludingAnon returns the number of cards posted by a user.
func (store *Store) GetPostedCardsForNodeIncludingAnon(nodeID globalid.ID, count, skip int) ([]*model.Card, error) {
	stmt := `SELECT * FROM cards where owner_id = $1 AND thread_root_id IS NULL AND deleted_at IS NULL ORDER BY created_at DESC`
	if count > 0 {
		stmt = fmt.Sprintf("%s\nLIMIT %d", stmt, count)
		if skip > 0 {
			stmt = fmt.Sprintf("%s\nOFFSET %d", stmt, count*skip)
		}
	}
	cards := []*model.Card{}
	err := store.Select(&cards, stmt, nodeID)
	if err != nil {
		return nil, errors.Wrap(err, "GetPostedCardsForNodeIncludingAnon failed")
	}

	return cards, nil
}

func (store *Store) DeleteCard(id globalid.ID) error {
	_, err := store.Exec(`UPDATE cards SET deleted_at=now() WHERE id = $1 OR thread_root_id = $1`, id)
	if err != nil {
		return errors.Wrap(err, "DeleteCard failed")
	}
	_, err = store.Exec(`DELETE FROM notifications_comments WHERE card_id = $1`, id)
	return errors.Wrap(err, "DeleteCard failed")
}

func (store *Store) ShadowbanCard(cardID globalid.ID) error {
	_, err := store.Exec(`UPDATE cards SET shadowbanned_at=$1 WHERE id = $2`, time.Now().UTC(), cardID)
	if err != nil {
		return errors.Wrap(err, "ShadowbanCard failed")
	}
	return nil
}

func (store *Store) ShadowbanAllCardsForUser(userID globalid.ID) error {
	_, err := store.Exec(`UPDATE cards SET shadowbanned_at=$1 WHERE owner_id = $2`, time.Now().UTC(), userID)
	if err != nil {
		return errors.Wrap(err, "ShadowbaAllCardsForUser failed")
	}
	return nil
}

func (store *Store) UnshadowbanCard(cardID globalid.ID) error {
	_, err := store.Exec(`UPDATE cards SET shadowbanned_at = NULL WHERE id = $2`, time.Now().UTC(), cardID)
	if err != nil {
		return errors.Wrap(err, "UnshadowbanCard failed")
	}
	return nil
}

func (store *Store) DeleteAllCardsForUser(userID globalid.ID) error {
	_, err := store.Exec(`UPDATE cards SET deleted_at=now() WHERE owner_id = $1`, userID)
	if err != nil {
		return errors.Wrap(err, "DeleteCard failed")
	}
	return nil
}

func (store *Store) SubscribersForCard(cardID globalid.ID, typ string) ([]globalid.ID, error) {
	var ids []globalid.ID
	err := store.Select(&ids, `
		SELECT user_id
		FROM subscriptions
		WHERE card_id = $1 AND type = $2
	`, cardID, typ)
	return ids, errors.Wrap(err, "SubscribersForCard failed")
}

func (store *Store) CountCardsWithNoReactions() (int, error) {
	var ct int

	query, args, err := sqlx.In(`
		SELECT COUNT(*)
		FROM (
			(
				SELECT id FROM cards WHERE thread_root_id IS null
			) EXCEPT (
				SELECT DISTINCT card_id FROM user_reactions WHERE type IN (?)
			)
		) AS t1
		`, model.ExplicitCardReactions)
	if err != nil {
		return 0, err
	}
	query = store.Rebind(query)
	err = store.Get(&ct, query, args...)
	if err != nil {
		return 0, errors.Wrap(err, "CountCardsWithNoReactions failed")
	}

	return ct, nil
}

type CardComment struct {
	ByName         string      `db:"by_name"`
	ImagePath      string      `db:"image_path"`
	IsAnonymous    bool        `db:"is_anonymous"`
	Timestamp      time.Time   `db:"created_at"`
	AuthorUsername string      `db:"author_username"`
	AuthorID       globalid.ID `db:"id"`
}

func (cc *CardComment) DisplayName() string {
	r := cc.ByName
	if cc.IsAnonymous {
		r = "!" + r
	}
	return r
}

type CommentNotificationExportData struct {
	Comments              []*CardComment `db:"-"`
	RootOwnerID           globalid.ID    `db:"owner_id"`
	RootIsAnonymous       bool           `db:"root_is_anon"`
	PosterName            string         `db:"display_name"`
	PosterImage           string         `db:"profile_image_path"`
	CardContent           string         `db:"content"`
	ThreadRootID          globalid.ID    `db:"thread_root_id"`
	LatestCommentID       globalid.ID    `db:"latest_comment_id"`
	LatestCommentParentID globalid.ID    `db:"latest_comment_parent_id"`
}

func (store *Store) GetCommentExportData(n *model.Notification) (*CommentNotificationExportData, error) {
	var comments []*CardComment

	var data CommentNotificationExportData

	err := store.Select(&comments, `
			SELECT
				users.id,
				COALESCE(anonymous_aliases.username, users.display_name) as "by_name",
				COALESCE(anonymous_aliases.profile_image_path, users.profile_image_path) as "image_path",
				cards.alias_id IS NOT NULL as "is_anonymous",
				MAX(cards.created_at) as "created_at",
				COALESCE(anonymous_aliases.username, users.username) as "author_username"
			FROM notifications_comments
				LEFT JOIN cards ON notifications_comments.card_id = cards.id
				LEFT JOIN users ON cards.owner_id = users.id
				LEFT JOIN anonymous_aliases ON cards.alias_id = anonymous_aliases.id
			WHERE notifications_comments.notification_id = $1
	  	GROUP BY by_name, image_path, is_anonymous, author_username, users.id
			ORDER BY created_at DESC`, n.ID)
	if err != nil {
		return nil, errors.Wrap(err, "GetCommentExportData failed")
	}

	err = store.Get(&data, `
		SELECT
			cards.owner_id,
			cards.alias_id IS NOT NULL as "root_is_anon",
			COALESCE(anonymous_aliases.username, users.display_name) as "display_name",
			COALESCE(anonymous_aliases.profile_image_path, users.profile_image_path) as "profile_image_path",
			cards.content,
			COALESCE(cards.thread_root_id, cards.id) as "thread_root_id",
			latestComment.card_id as "latest_comment_id",
			latestComment.thread_reply_id as "latest_comment_parent_id"
		FROM cards
			LEFT JOIN users ON cards.owner_id = users.id
			LEFT JOIN anonymous_aliases ON cards.alias_id = anonymous_aliases.id
			LEFT JOIN (
				SELECT $1::uuid as "join_id", notifications_comments.*, cards.thread_reply_id
				FROM notifications_comments LEFT JOIN cards on notifications_comments.card_id = cards.id
				WHERE notifications_comments.notification_id = $2
				ORDER BY notifications_comments.created_at DESC
				LIMIT 1
			) as latestComment ON cards.id = latestComment.join_id
		WHERE cards.id = $1`, n.TargetID, n.ID)

	if err != nil {
		return nil, errors.Wrap(err, "GetCommentExportData failed")
	}

	data.Comments = comments

	return &data, nil
}

// SaveInvite saves an invite to the store
func (store *Store) SaveNotificationComment(m *model.NotificationComment) error {
	return saveNotificationComment(store, m)
}

// SaveInvite saves an invite to the store
func (tx *Tx) SaveNotificationComment(m *model.NotificationComment) error {
	return saveNotificationComment(tx, m)
}

func saveNotificationComment(e sqlx.Ext, m *model.NotificationComment) error {
	if m == nil {
		return errors.New("provided model can not be nil")
	}

	tn := time.Now().UTC()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = tn
	}
	m.UpdatedAt = tn

	_, err := sqlx.NamedExec(e,
		`INSERT INTO notifications_comments
		(
			notification_id,
			card_id,
			created_at,
			updated_at
		)
		VALUES
		(
			:notification_id,
			:card_id,
			:created_at,
			:updated_at
		)
		ON CONFLICT(notification_id,card_id) DO UPDATE
		SET
			notification_id  = :notification_id,
			card_id          = :card_id,
			created_at       = :created_at,
			updated_at       = :updated_at
		WHERE notifications_comments.notification_id = :notification_id AND notifications_comments.card_id = :card_id `, m)
	return errors.Wrap(err, "SaveNotificationComment failed")
}

func (store *Store) FeaturedCommentsForUser(userID, cardID globalid.ID) ([]*model.Card, error) {
	//featured comments are currently:
	// The latest comment not belonging to you
	var cards []*model.Card
	err := store.Select(&cards, `SELECT * FROM cards WHERE thread_reply_id = $1 AND owner_id != $2 AND deleted_at IS NULL AND shadowbanned_at IS NULL ORDER BY created_at DESC LIMIT 1`, cardID, userID)
	if err != nil {
		return nil, errors.Wrap(err, "FeaturedCommentsForUser failed")
	}
	return cards, nil
}

func (store *Store) FeaturedCommentsForUserByCardIDs(userID globalid.ID, cardIDs []globalid.ID) (map[globalid.ID]*model.Card, error) {
	var cards []*model.Card
	err := store.Select(&cards, `SELECT   *
	                             FROM     cards
								 WHERE    thread_reply_id = ANY($1::uuid[])
								 AND      owner_id != $2
								 AND      deleted_at IS NULL
								 AND      shadowbanned_at IS NULL
								 ORDER BY created_at DESC LIMIT 1`, pq.Array(cardIDs), userID)
	if err != nil {
		return nil, errors.Wrap(err, "FeaturedCommentsForUserByCardIDs failed")
	}
	result := make(map[globalid.ID]*model.Card)
	id := cardIDs[0]
	for _, card := range cards {
		if id != card.ThreadReplyID {
			id := card.ThreadReplyID
			result[id] = card
		}
	}
	return result, nil
}

func (store *Store) MuteThread(userID, threadRootID globalid.ID) error {
	_, err := store.Exec(`INSERT INTO thread_mutes (user_id, thread_root_id) VALUES ($1, $2) ON CONFLICT (user_id, thread_root_id) DO NOTHING`, userID, threadRootID)
	if err != nil {
		return errors.Wrap(err, "MuteThread failed")
	}

	return nil
}

func (store *Store) UnmuteThread(userID, threadRootID globalid.ID) error {
	_, err := store.Exec(`DELETE FROM thread_mutes WHERE user_id = $1 AND threadRootID = $2`, userID, threadRootID)
	if err != nil {
		return errors.Wrap(err, "UnmuteUser failed")
	}
	return nil
}

func (store *Store) CountPostsByAliasInThread(aliasID, threadRootID globalid.ID) (int, error) {
	var count int
	err := store.Get(&count, `SELECT COUNT(*) FROM cards WHERE alias_id = $1 AND thread_root_id = $2`, aliasID, threadRootID)
	return count, errors.Wrap(err, "CountPostsByAliasInThread failed")
}

func (store *Store) AddToCoinsEarned(cardID globalid.ID, amount int) error {
	_, err := store.Exec(`UPDATE cards SET coins_earned = coins_earned + $1 WHERE id = $2`, amount, cardID)
	return errors.Wrap(err, "AddToCoinsEarned failed")
}

func (store *Store) GetUsersNeedingFirstPostNotification() ([]*model.User, error) {
	// users who have gotten at least one post or comment on any post and have no notif for first or popular posts
	var cards []*model.User
	err := store.Select(&cards, `
		SELECT users.* FROM (
			SELECT
			  DISTINCT cards.owner_id as user_id
			FROM
			  user_reactions
			  LEFT JOIN cards ON user_reactions.card_id = cards.id
			UNION
			SELECT
			  DISTINCT parents.owner_id as user_id
			FROM
			  cards
			  LEFT JOIN cards AS parents ON cards.thread_reply_id = parents.id
			WHERE
			  cards.thread_reply_id IS NOT null
			EXCEPT
			SELECT
			  user_id
			FROM notifications
			WHERE type = $1 OR type = $2
		) as root LEFT JOIN users ON root.user_id = users.id
		`, model.FirstPostActivityType, model.PopularPostType)
	if err != nil {
		return nil, errors.Wrap(err, "GetPostsNeedingFirstPostNotification failed")
	}
	return cards, nil
}

func (store *Store) UsersWithPopularPostsSinceTime(t time.Time) ([]*model.User, error) {
	// Get posts and comments that have at least one reply or comment and no first-post notificaiton
	var cards []*model.User
	err := store.Select(&cards, `
		SELECT users.* FROM (
			SELECT
				owner_id
			FROM
				cards
				LEFT JOIN (
					SELECT card_id, count(*)
					FROM user_reactions
					WHERE created_at >= $1
					GROUP BY card_id
				) AS rcount ON cards.id = rcount.card_id
				LEFT JOIN (
					SELECT thread_reply_id AS card_id, count(*)
					FROM cards
					WHERE created_at >= $1
					GROUP BY thread_reply_id
				) AS ccount ON cards.id = ccount.card_id
			GROUP BY owner_id
			HAVING SUM(COALESCE(rcount.count,0)) + SUM(COALESCE(ccount.count, 0)) >= 5
		) as root LEFT JOIN users ON root.owner_id = users.id
		`, t)
	if err != nil {
		return nil, errors.Wrap(err, "GetPopularPostSinceTime failed")
	}
	return cards, nil
}

type PopularPostNotificationExportData struct {
	ThreadRootID        globalid.ID `db:"thread_root_id"`
	CommentCardID       globalid.ID `db:"comment_card_id"`
	CommentCardUsername string      `db:"comment_card_username"`
	ParentCommentID     globalid.ID `db:"parent_comment_id"`
}

func (store *Store) GetPopularPostNotificationExportData(n *model.Notification) (*PopularPostNotificationExportData, error) {
	var data PopularPostNotificationExportData
	err := store.Get(&data, `
		SELECT
		    coalesce(thread_root_id, cards.id) AS thread_root_id,
				CASE WHEN thread_root_id IS NULL THEN NULL ELSE cards.id END AS comment_card_id,
				CASE WHEN thread_root_id IS NULL THEN NULL ELSE COALESCE(anonymous_aliases.username, users.username) END AS comment_card_username,
		    CASE WHEN thread_root_id IS NOT NULL AND thread_root_id != thread_reply_id THEN thread_reply_id ELSE NULL END AS parent_comment_id
		FROM
			cards
			LEFT JOIN users ON cards.owner_id = users.id
			LEFT JOIN anonymous_aliases ON cards.alias_id = anonymous_aliases.id
		WHERE cards.id = $1
	`, n.TargetID)
	return &data, errors.Wrap(err, "GetPopularPostNotificationExportData failed")
}
