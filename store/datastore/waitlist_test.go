package datastore_test

import (
	"testing"

	"github.com/october93/engine/model"
	"github.com/october93/engine/store/datastore"
	"github.com/october93/engine/test"
)

func TestWaitlist(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping - short test detected")
	}

	cfg := datastore.NewTestConfig()
	cfg.Database = "engine_datastore_waitlist_test"
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

	w := model.WaitlistEntry{
		Email:   "chad@october.news",
		Comment: "October's intern",
	}

	t.Run("create waitlist", func(t *testing.T) {
		if err = store.SaveWaitlistEntry(&w); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("get waitlist", func(t *testing.T) {
		var waitlist []*model.WaitlistEntry
		waitlist, err = store.GetWaitlist()
		if err != nil {
			t.Fatal(err)
		}
		if len(waitlist) != 1 {
			t.Fatalf("expected waitlist to have one entry")
		}
		if waitlist[0].Email != w.Email {
			t.Fatalf("expected Email to be %s, actual: %s", w.Email, waitlist[0].Email)
		}
		if waitlist[0].Comment != w.Comment {
			t.Fatalf("expected Comment to be %s, actual: %s", w.Comment, waitlist[0].Comment)
		}
		if waitlist[0].CreatedAt.IsZero() {
			t.Fatal("expected CreatedAt to be set")
		}
	})

	t.Run("delete waitlist", func(t *testing.T) {
		err = store.DeleteWaitlistEntry(w.Email)
		if err != nil {
			t.Fatal(err)
		}
		waitlist, err := store.GetWaitlist()
		if err != nil {
			t.Fatal(err)
		}
		if len(waitlist) != 0 {
			t.Fatalf("expected waitlist to be empty")
		}
		err = store.DeleteWaitlistEntry(w.Email)
		if err != nil {
			t.Fatal(err)
		}
	})
}
