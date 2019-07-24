package datastore_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/kr/pretty"
	"github.com/october93/engine/model"
	datastore "github.com/october93/engine/store/datastore"
	"github.com/october93/engine/test"
)

func TestSettings(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping - short test detected")
	}

	cfg := datastore.NewTestConfig()
	cfg.Database = "engine_datastore_settings_test"
	db := test.DBInit(t, cfg)
	// drop the database after the test is finished
	defer test.DBCleanup(t, db)

	store, err := datastore.New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	// Close the store
	defer func() {
		if e := store.Close(); e != nil {
			t.Fatalf("store.Close() failed.: %s", e)
		}
	}()

	settings := model.Settings{
		MaintenanceMode: true,
		SignupsFrozen:   true,
	}

	t.Run("save_settings", func(t *testing.T) {
		err := store.SaveSettings(&settings)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("get_settings", func(t *testing.T) {
		got, err := store.GetSettings()
		if err != nil {
			t.Fatal(err)
		}
		exp := &settings
		compareSettings(t, exp, got)
	})

	settings.MaintenanceMode = false

	t.Run("update_settings", func(t *testing.T) {
		err := store.SaveSettings(&settings)
		if err != nil {
			t.Fatal(err)
		}
		got, err := store.GetSettings()
		if err != nil {
			t.Fatal(err)
		}
		exp := &settings
		compareSettings(t, exp, got)
	})

	settings2 := model.Settings{
		MaintenanceMode: true,
		SignupsFrozen:   true,
	}

	t.Run("save_second_settings", func(t *testing.T) {
		err := store.SaveSettings(&settings2)
		if err != nil {
			t.Fatal(err)
		}
		got, err := store.GetSettings()
		if err != nil {
			t.Fatal(err)
		}
		exp := &settings2
		compareSettings(t, exp, got)
	})
}

func compareSettings(t *testing.T, exp, got *model.Settings) {
	t.Helper()
	if exp.CreatedAt.Unix() != got.CreatedAt.Unix() {
		t.Fatalf("unexpected CreatedAt. exp %v, got %v", exp.CreatedAt, got.CreatedAt)
	}
	if exp.UpdatedAt.Unix() != got.UpdatedAt.Unix() {
		t.Fatalf("unexpected UpdatedAt. exp %v, got %v", exp.UpdatedAt, got.UpdatedAt)
	}
	exp.CreatedAt = time.Time{}
	exp.UpdatedAt = time.Time{}
	got.CreatedAt = time.Time{}
	got.UpdatedAt = time.Time{}
	if !reflect.DeepEqual(exp, got) {
		t.Fatalf("unexpected result, diff: %v", pretty.Diff(exp, got))
	}
}
