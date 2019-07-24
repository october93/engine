package datastore_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/kr/pretty"
	"github.com/october93/engine/model"
	"github.com/october93/engine/store/datastore"
	"github.com/october93/engine/test"
)

func TestOAuthAccounts(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping - short test detected")
	}

	cfg := datastore.NewTestConfig()
	cfg.Database = "engine_datastore_oauth_accounts_test"
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

	u1 := model.User{
		Username:     "peggy",
		DisplayName:  "Peggy",
		Email:        "peggy@october.news",
		PasswordHash: "hash",
		PasswordSalt: "salt",
	}

	u2 := model.User{
		Username:     "sue",
		DisplayName:  "Sue",
		Email:        "sue@october.news",
		PasswordHash: "hash",
		PasswordSalt: "salt",
	}

	oa1 := &model.OAuthAccount{
		Provider: model.FacebookProvider,
		Subject:  "123",
	}

	oa2 := &model.OAuthAccount{
		Provider: model.FacebookProvider,
		Subject:  "abc",
	}

	t.Run("create_oauth_account", func(t *testing.T) {
		if err := store.SaveUser(&u1); err != nil {
			t.Fatal(err)
		}
		oa1.UserID = u1.ID
		if err := store.SaveOAuthAccount(oa1); err != nil {
			t.Fatal(err)
		}
		if err := store.SaveUser(&u2); err != nil {
			t.Fatal(err)
		}
		oa2.UserID = u2.ID
		if err := store.SaveOAuthAccount(oa2); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("get_oauth_account_by_subject", func(t *testing.T) {
		got, err := store.GetOAuthAccountBySubject(oa1.Subject)
		if err != nil {
			t.Fatal(err)
		}
		exp := oa1
		compareOAuthAccounts(t, exp, got)
		got, err = store.GetOAuthAccountBySubject(oa2.Subject)
		if err != nil {
			t.Fatal(err)
		}
		exp = oa2
		compareOAuthAccounts(t, exp, got)
	})

}

func compareOAuthAccounts(t *testing.T, e, g *model.OAuthAccount) {
	t.Helper()
	exp := *e
	got := *g

	exp.CreatedAt = exp.CreatedAt.Truncate(time.Millisecond)
	got.CreatedAt = got.CreatedAt.Truncate(time.Millisecond)
	if !exp.CreatedAt.Equal(got.CreatedAt) {
		t.Fatalf("unexpected CreatedAt.  exp %v, got %v", exp.CreatedAt, got.CreatedAt)
	}

	exp.UpdatedAt = exp.UpdatedAt.Truncate(time.Millisecond)
	got.UpdatedAt = got.UpdatedAt.Truncate(time.Millisecond)
	if !exp.UpdatedAt.Equal(got.UpdatedAt) {
		t.Fatalf("unexpected UpdatedAt.  exp %v, got %v", exp.UpdatedAt, got.UpdatedAt)
	}

	exp.CreatedAt = time.Time{}
	got.CreatedAt = time.Time{}

	exp.UpdatedAt = time.Time{}
	got.UpdatedAt = time.Time{}

	if !reflect.DeepEqual(exp, got) {
		t.Fatalf("compareOAuthAccounts(): unexpected result, diff: %v", pretty.Diff(exp, got))
	}
}
