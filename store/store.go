package store

import (
	"database/sql"

	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/model"
	"github.com/october93/engine/store/datastore"
	"github.com/pkg/errors"
)

var (
	RootUser = &model.User{
		Username:    "root",
		DisplayName: "root",
		Email:       "tech@october.news",
		Admin:       true,
	}
)

// Store is the persistence access layer to the graph and the domain model.
// Store acts as a wrapper in order to make access to both parts more
// convenient.
type Store struct {
	*datastore.Store
	config *Config
	log    log.Logger
}

// NewStore returns a new instance of Store.
func NewStore(c *Config, l log.Logger) (*Store, error) {
	store, err := datastore.New(c.Datastore)
	if err != nil {
		return nil, err
	}
	return &Store{
		Store:  store,
		config: c,
		log:    l,
	}, nil
}

// EnsureRootUser makes sure there is always one user available to interact
// with Engine. This is used in particular to make the administration panel
// accessible.
func (s *Store) EnsureRootUser() error {
	user, err := s.GetUserByUsername(RootUser.Username)
	if err == nil {
		err = user.SetPassword(s.config.RootUserPassword)
		if err != nil {
			return err
		}
		return s.SaveUser(user)
	}
	err = RootUser.SetPassword(s.config.RootUserPassword)
	if err != nil {
		return err
	}
	return s.SaveUser(RootUser)
}

func (s *Store) EnsureSettings() (*model.Settings, error) {
	settings, err := s.GetSettings()
	if errors.Cause(err) == sql.ErrNoRows {
		settings = &model.Settings{}
		return settings, s.SaveSettings(settings)
	}
	return settings, err
}

// Populate inserts necessary data like anonymous aliases or tags.
func (s *Store) Populate() error {
	return s.Store.CreateAnonymousAliases()
}

func (s *Store) Close() error {
	// TODO: close snapshotter?
	return s.Store.Close()
}
