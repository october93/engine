package datastore_test

import (
	"encoding/json"
	"testing"

	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
	"github.com/october93/engine/rpc"
	"github.com/october93/engine/store/datastore"
	"github.com/october93/engine/test"
)

func TestActivities(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping - short test detected")
	}

	cfg := datastore.NewTestConfig()
	cfg.Database = "engine_datastore_activities_test"
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

	u := model.User{
		ID:           globalid.Next(),
		Email:        "kafka@october.news",
		Username:     "kafka",
		DisplayName:  "Kafka",
		PasswordHash: "secret",
		PasswordSalt: "salty secret",
	}

	data := rpc.ReactToCardParams{
		CardID:   globalid.Next(),
		Strength: 1.0,
		Reaction: "boost",
	}
	d, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	a := model.Activity{
		RPC:    "reactToCard",
		Data:   d,
		UserID: u.ID,
		Error:  "unknown card",
	}

	t.Run("create_activity", func(t *testing.T) {
		if err := store.SaveUser(&u); err != nil {
			t.Fatal(err)
		}
		if err := store.SaveActivity(&a); err != nil {
			t.Fatal(err)
		}
	})
}
