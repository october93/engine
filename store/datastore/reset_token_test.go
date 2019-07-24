package datastore_test

import (
	"testing"
	"time"

	model "github.com/october93/engine/model"
	datastore "github.com/october93/engine/store/datastore"
	"github.com/october93/engine/test"
)

func Test_ResetTokens(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping - short test detected")
	}

	cfg := datastore.NewTestConfig()
	cfg.Database = "engine_datastore_reset_tokens_test"
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

	resetToken := model.ResetToken{
		TokenHash: "secret",
		Expires:   time.Date(2008, time.November, 10, 23, 0, 0, 0, time.UTC),
	}
	user := model.User{
		Email:        "john@smith.com",
		Username:     "john",
		DisplayName:  "John Smith",
		PasswordHash: "secret",
		PasswordSalt: "salty secret",
	}

	t.Run("create_reset_token", func(t *testing.T) {
		if err := store.SaveUser(&user); err != nil {
			t.Fatal(err)
		}
		resetToken.UserID = user.ID
		if err := store.SaveResetToken(&resetToken); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("get_reset_token", func(t *testing.T) {
		rt, err := store.GetResetToken(user.ID)
		if err != nil {
			t.Fatal(err)
		}
		compareResetTokens(t, &resetToken, rt)
	})
}

func compareResetTokens(t *testing.T, exp, got *model.ResetToken) {
	t.Helper()
	if exp.TokenHash != got.TokenHash {
		t.Fatalf("unexpected Token.  exp %v, got %v", exp.TokenHash, got.TokenHash)
	}
	if exp.UserID != got.UserID {
		t.Fatalf("unexpected UserID.  exp %v, got %v", exp.UserID, got.UserID)
	}
	if exp, got := exp.Expires.Truncate(time.Millisecond), got.Expires.Truncate(time.Millisecond); !exp.Equal(got) {
		t.Fatalf("unexpected Expires.  exp %v, got %v", exp, got)
	}
	if exp, got := exp.CreatedAt.Truncate(time.Millisecond), got.CreatedAt.Truncate(time.Millisecond); !exp.Equal(got) {
		t.Fatalf("unexpected CreatedAt.  exp %v, got %v", exp, got)
	}
	if !exp.UpdatedAt.Truncate(time.Millisecond).Before(got.UpdatedAt) && !exp.UpdatedAt.Truncate(time.Millisecond).Equal(got.UpdatedAt.Truncate(time.Millisecond)) {
		t.Fatalf("unexpected UpdatedAt.  exp %v to be before %v", exp.UpdatedAt, got.UpdatedAt)
	}
}
