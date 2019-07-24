package datastore_test

import (
	"testing"
	"time"

	"github.com/october93/engine/kit/globalid"
	model "github.com/october93/engine/model"
	datastore "github.com/october93/engine/store/datastore"
	"github.com/october93/engine/test"
)

func Test_Sessions(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping - short test detected")
	}

	cfg := datastore.NewTestConfig()
	cfg.Database = "engine_datastore_sessions_test"
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

	user := model.User{
		ID:           globalid.Next(),
		Email:        "kafka@october.news",
		Username:     "kafka",
		DisplayName:  "Kafka",
		PasswordHash: "secret",
		PasswordSalt: "salty secret",
	}
	session := model.Session{
		UserID: user.ID,
	}

	t.Run("create_session", func(t *testing.T) {
		if err := store.SaveUser(&user); err != nil {
			t.Fatal(err)
		}
		if err := store.SaveSession(&session); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("get_session", func(t *testing.T) {
		slg, err := store.GetSession(session.ID)
		if err != nil {
			t.Fatal(err)
		}
		compareSessions(t, &session, slg)
	})

	t.Run("get_sessions", func(t *testing.T) {
		s, err := store.GetSessions()
		if err != nil {
			t.Fatal(err)
		}
		if exp, got := 1, len(s); exp != got {
			t.Fatalf("unexpected length of sessions.  exp %d, got %d", exp, got)
		}
		compareSessions(t, &session, s[0])
		if s[0].GetUser() == nil {
			t.Fatal("expected User field to be set")
		}
		if s[0].User.Username != user.Username {
			t.Fatalf("expected username %s, actual: %s", user.Username, s[0].User.Username)
		}
	})

	t.Run("delete_expired_sessions", func(t *testing.T) {
		i, err := store.DeleteExpiredSessions()
		if err != nil {
			t.Fatal(err)
		}
		if i != 0 {
			t.Fatalf("expected zero sessions to be deleted, got %d", i)
		}
	})

	t.Run("delete_sessions", func(t *testing.T) {
		if err := store.DeleteSession(session.ID); err != nil {
			t.Fatal(err)
		}
		s, err := store.GetSessions()
		if err != nil {
			t.Fatal(err)
		}
		if exp, got := 0, len(s); exp != got {
			t.Fatalf("unexpected length of sessions.  exp %d, got %d", exp, got)
		}
	})
}

func compareSessions(t *testing.T, exp, got *model.Session) {
	t.Helper()
	if exp.ID != got.ID {
		t.Fatalf("unexpected ID.  exp %v, got %v", exp.ID, got.ID)
	}
	if exp.UserID != got.UserID {
		t.Fatalf("unexpected UserID.  exp %v, got %v", exp.UserID, got.UserID)
	}
	if exp, got := exp.CreatedAt.Truncate(time.Millisecond), got.CreatedAt.Truncate(time.Millisecond); !exp.Equal(got) {
		t.Fatalf("unexpected CreatedAt.  exp %v, got %v", exp, got)
	}
	if !exp.UpdatedAt.Truncate(time.Millisecond).Before(got.UpdatedAt) && !exp.UpdatedAt.Truncate(time.Millisecond).Equal(got.UpdatedAt.Truncate(time.Millisecond)) {
		t.Fatalf("unexpected UpdatedAt.  exp %v to be before %v", exp.UpdatedAt, got.UpdatedAt)
	}
}
