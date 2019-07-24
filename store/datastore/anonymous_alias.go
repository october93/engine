package datastore

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

// ErrNoAnonymousAliasLeft is used when all anonymous aliases have been
// assigned to individual users.
var ErrNoAnonymousAliasLeft = errors.New("no anonymous aliases left")

// CreateAnonymousAliases will update an existing user or create a new user record in the
// datastore.
func (store *Store) CreateAnonymousAliases() error {
	emojiIcons := []string{
		"airplane",
		"bear",
		"carrot",
		"corn",
		"duck",
		"hammer",
		"lollipop",
		"pepper",
		"snake",
		"whale",
		"alien",
		"bee",
		"cat",
		"cow",
		"eagle",
		"horse",
		"monkey",
		"pineapple",
		"ruler",
		"strawberry",
		"zebra",
		"ant",
		"bread",
		"caterpillar",
		"crab",
		"egg",
		"icecream",
		"mouse",
		"popcorn",
		"sandwich",
		"taco",
		"apple",
		"broccoli",
		"cheese",
		"crocodile",
		"frog",
		"key",
		"mushroom",
		"porcupine",
		"satellite",
		"thumbtack",
		"avocado",
		"burrito",
		"cherry",
		"croissant",
		"giraffe",
		"kiwi",
		"octopus",
		"potato",
		"saxophone",
		"tiger",
		"balloon",
		"butterfly",
		"chipmunk",
		"deer",
		"glasses",
		"koala",
		"orange",
		"pretzel",
		"scarf",
		"tomato",
		"banana",
		"cake",
		"chocolate",
		"dog",
		"goat",
		"lemon",
		"owl",
		"rabbit",
		"scooter",
		"unicorn",
		"basketball",
		"canoe",
		"coconut",
		"dolphin",
		"gorilla",
		"lion",
		"panda",
		"rhino",
		"shark",
		"violin",
		"bat",
		"brain",
		"cookie",
		"donut",
		"guitar",
		"lizard",
		"pear",
		"rocket",
		"snail",
		"watermelon",
	}

	icons := []string{
		/*"anchor",
		"apple",
		"bell",
		"cake",
		"candle",
		"car",
		"cat",
		"crab",
		"duck",
		"fish",
		"flower",
		"key",
		"robot",
		"rocket",
		"star",
		"tie",
		"tree",*/
	}

	colors := []string{
		/*"red",
		"pink",
		"orange",
		"yellow",
		"green",
		"blue",
		"purple",*/
	}

	blacklist := map[string]bool{
		"redrocket":  true,
		"pinkrocket": true,
		"pinkcat":    true,
	}

	existingHandles := map[string]bool{}

	aliases := []*model.AnonymousAlias{}
	err := store.Select(&aliases, "SELECT * FROM anonymous_aliases ORDER BY updated_at DESC")

	if err != nil {
		return err
	}

	for _, alias := range aliases {
		existingHandles[alias.Username] = true
	}

	for _, icon := range icons {
		for _, color := range colors {
			handle := fmt.Sprintf("%v%v", color, icon)
			if !existingHandles[handle] && !blacklist[handle] {
				imagePath := fmt.Sprintf("%s/%s/%s/%s.png", store.config.ImageHost, store.config.AnonymousIconsPath, color, icon)
				newAlias := model.AnonymousAlias{
					DisplayName:      "Anonymous",
					Username:         handle,
					ProfileImagePath: imagePath,
				}
				err := saveHandle(store, &newAlias)
				if err != nil {
					return err
				}
			}
		}
	}

	for _, name := range emojiIcons {
		//removes underscore
		handle := strings.Replace(name, "_", "", -1)

		if !existingHandles[handle] && !blacklist[handle] {
			imagePath := fmt.Sprintf("https://s3-us-west-2.amazonaws.com/%s/%s.png", store.config.AnonymousIconsPath, name)
			newAlias := model.AnonymousAlias{
				DisplayName:      "Anonymous",
				Username:         handle,
				ProfileImagePath: imagePath,
			}
			err := saveHandle(store, &newAlias)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// GetAnonymousAlias ---
func (store *Store) GetAnonymousAlias(id globalid.ID) (*model.AnonymousAlias, error) {
	alias := model.AnonymousAlias{}
	err := store.Get(&alias, "SELECT * FROM anonymous_aliases where id = $1", id)
	return &alias, errors.Wrap(err, "GetAnonymousAlias failed")
}

func (store *Store) GetAnonymousAliasesByID(ids []globalid.ID) ([]*model.AnonymousAlias, error) {
	query, args, err := sqlx.In(`SELECT anonymous_aliases.* FROM unnest(ARRAY[?]::uuid[]) WITH ORDINALITY AS r(id, rn) JOIN anonymous_aliases USING (id) ORDER BY r.rn`, ids)
	if err != nil {
		return nil, err
	}
	query = store.Rebind(query)
	aliases := []*model.AnonymousAlias{}
	err = store.Select(&aliases, query, args...)
	return aliases, errors.Wrap(err, "GetAnonymousAliasByID failed")
}

func (store *Store) GetAnonymousAliasLastUsed(userID, threadRootID globalid.ID) (bool, error) {
	var aliasID globalid.ID
	err := store.Get(&aliasID, `
	    SELECT alias_id FROM (
			SELECT alias_id, created_at
			FROM cards
			WHERE owner_id = $1
			AND (thread_root_id = $2 OR id = $2)
			UNION
			(SELECT user_reactions.alias_id, user_reactions.created_at
			 FROM user_reactions
			 WHERE user_id = $1
			 AND   card_id = ANY(SELECT id FROM cards WHERE thread_root_id = $2 OR id = $2)
			 AND   type = ANY (ARRAY[$3])
			)
			ORDER BY created_at DESC
			LIMIT 1
		) t`, userID, threadRootID, model.ReactionLike)

	if err == nil || err == sql.ErrNoRows {
		return aliasID != globalid.Nil, nil
	}

	return false, errors.Wrap(err, "GetAnonymousAliasLastUsed failed")
}

func (store *Store) GetAnonymousAliasByUsername(username string) (*model.AnonymousAlias, error) {
	alias := model.AnonymousAlias{}
	err := store.Get(&alias, "SELECT * FROM anonymous_aliases where username = $1", username)
	return &alias, errors.Wrap(err, "GetAnonymousAliasByUsername failed")
}

func (store *Store) GetAnonymousAliases() ([]*model.AnonymousAlias, error) {
	aliases := []*model.AnonymousAlias{}
	err := store.Select(&aliases, `SELECT * FROM anonymous_aliases ORDER BY updated_at DESC`)
	return aliases, errors.Wrap(err, "GetAnonymousAliases failed")
}

func (store *Store) GetUnusedAlias(cardID globalid.ID) (*model.AnonymousAlias, error) {
	var alias model.AnonymousAlias
	err := store.Get(&alias, `
		SELECT anonymous_aliases.*
		FROM   anonymous_aliases, cards
		WHERE  NOT inactive
		AND (CASE WHEN (SELECT COUNT(*) FROM cards WHERE cards.id = $1) > 0
			 THEN anonymous_aliases.id NOT IN (
				 SELECT uuid(value) FROM jsonb_each_text((SELECT author_to_alias FROM cards WHERE id = $1)))
			 ELSE
				TRUE
			 END)
		ORDER BY random() LIMIT 1`, cardID)
	if err == sql.ErrNoRows {
		return nil, ErrNoAnonymousAliasLeft
	} else if err != nil {
		return nil, errors.Wrap(err, "GetUnusedAlias failed")
	}
	return &alias, nil
}

func (store *Store) SaveAnonymousAlias(m *model.AnonymousAlias) error {
	return saveHandle(store, m)
}

func saveHandle(e sqlx.Ext, m *model.AnonymousAlias) error {
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
		`INSERT INTO anonymous_aliases
		(
			id,
			display_name,
			profile_image_path,
			username,
			created_at,
			updated_at,
			inactive
		)
		VALUES
		(
			:id,
			:display_name,
			:profile_image_path,
			:username,
			:created_at,
			:updated_at,
			:inactive
		)
		ON CONFLICT(id) DO UPDATE
		SET
			display_name       = :display_name,
			profile_image_path = :profile_image_path,
			username           = :username,
			created_at         = :created_at,
			updated_at         = :updated_at,
			inactive           = :inactive
		WHERE anonymous_aliases.ID = :id`, m)
	return errors.Wrap(err, "SaveAnonymousAlias failed")
}
