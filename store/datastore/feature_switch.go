package datastore

import (
	"time"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

func (store *Store) SaveFeatureSwitch(m *model.FeatureSwitch) error {
	return saveFeatureSwitch(store, m)
}

func (tx *Tx) SaveFeatureSwitch(m *model.FeatureSwitch) error {
	return saveFeatureSwitch(tx, m)
}

func saveFeatureSwitch(e sqlx.Ext, m *model.FeatureSwitch) error {
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
		`INSERT INTO feature_switches
	(
		id,
		name,
		state,
		testing_users,
		created_at,
		updated_at
	)
	VALUES
	(
		:id,
		:name,
		:state,
		:testing_users,
		:created_at,
		:updated_at
	)
	ON CONFLICT(id) DO UPDATE
	SET
		id               = :id,
		name             = :name,
		state            = :state,
		testing_users    = :testing_users,
		created_at       = :created_at,
		updated_at       = :updated_at
	WHERE feature_switches.ID = :id`, m)
	return errors.Wrap(err, "SaveFeatureSwitch failed")
}

// GetSavedLink returns the saved link for the provided ID
func (store *Store) GetFeatureSwitches() ([]*model.FeatureSwitch, error) {
	featureSwitch := []*model.FeatureSwitch{}
	err := store.Select(&featureSwitch, "SELECT * FROM feature_switches")
	return featureSwitch, errors.Wrap(err, "GetFeatureSwitches failed")
}

// GetSavedLink returns the saved link for the provided ID
func (store *Store) GetSwitchByName(name string) (*model.FeatureSwitch, error) {
	feature := model.FeatureSwitch{}
	err := store.Get(&feature, "SELECT * FROM feature_switches WHERE name = $1", name)
	return &feature, errors.Wrap(err, "GetSwitchByName failed")
}

// GetSavedLink returns the saved link for the provided ID
func (store *Store) DeleteSwitch(id globalid.ID) error {
	_, err := store.Exec("DELETE FROM feature_switches WHERE id = $1", id)
	return errors.Wrap(err, "DeleteSwitch failed")
}

// GetSavedLink returns the saved link for the provided ID
func (store *Store) GetOnSwitches() ([]*model.FeatureSwitch, error) {
	featureSwitch := []*model.FeatureSwitch{}
	err := store.Select(&featureSwitch, "SELECT * FROM feature_switches where state != 'off'")
	return featureSwitch, errors.Wrap(err, "GetOnSwitches failed")
}

// GetSavedLink returns the saved link for the provided ID
func (store *Store) ToggleUserForFeature(featureID, userID globalid.ID) error {
	fS := model.FeatureSwitch{}
	err := store.Get(&fS, "SELECT * FROM feature_switches where id = $1", featureID)

	if err != nil {
		return errors.Wrap(err, "ToggleUserForFeature failed")
	}

	if fS.TestingUsers[userID] {
		delete(fS.TestingUsers, userID)
	} else {
		fS.TestingUsers[userID] = true
	}

	err = store.SaveFeatureSwitch(&fS)
	return err
}

// GetSavedLink returns the saved link for the provided ID
func (store *Store) ChangeFeatureSwitchState(featureID globalid.ID, state string) error {
	fS := model.FeatureSwitch{}
	err := store.Get(&fS, "SELECT * FROM feature_switches where id = $1", featureID)
	if err != nil {
		return errors.Wrap(err, "ChangeFeatureSwitchState failed")
	}

	fS.State = state

	err = store.SaveFeatureSwitch(&fS)
	return err
}
