package datastore_test

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/kr/pretty"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
	datastore "github.com/october93/engine/store/datastore"
	"github.com/october93/engine/test"
	"github.com/pkg/errors"
)

func Test_Users(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping - short test detected")
	}

	cfg := datastore.NewTestConfig()
	cfg.Database = "engine_datastore_users_test"
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
		ID:               globalid.Next(),
		DisplayName:      "Rob Pike",
		FirstName:        "Rob",
		LastName:         "Pike",
		ProfileImagePath: "/image/user.gif",
		Bio:              "I created go",
		Email:            "rob@golang.org",
		PasswordHash:     "hash",
		PasswordSalt:     "salt",
		Username:         "commander",
		Devices:          model.Devices{"platform": model.Device{Token: "token", Platform: "platform"}},
		Admin:            true,
		SearchKey:        "searchKey",
	}

	t.Run("create_user", func(t *testing.T) {
		if err := store.SaveUser(&u); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("get_users", func(t *testing.T) {
		us, err := store.GetUsers()
		if err != nil {
			t.Fatal(err)
		}
		if exp, got := 1, len(us); exp != got {
			t.Fatalf("unexpected length of users.  exp %d, got %d", exp, got)
		}
		exp := &u
		got := us[0]
		compareUsers(t, exp, got)
	})

	t.Run("get_user", func(t *testing.T) {
		ug, err := store.GetUser(u.ID)
		if err != nil {
			t.Fatal(err)
		}
		compareUsers(t, &u, ug)
	})

	t.Run("get_user_by_email", func(t *testing.T) {
		ue, err := store.GetUserByEmail(u.Email)
		if err != nil {
			t.Fatal(err)
		}
		compareUsers(t, &u, ue)
	})

	t.Run("update_user", func(t *testing.T) {
		u.DisplayName = "Bob Smith"
		u.FirstName = "Bob"
		u.LastName = "Smith"
		u.ProfileImagePath = "/image/bob/jpg"
		u.Bio = "I'm Bob"
		u.Email = "bob.smith@yahoo.com"
		u.PasswordHash = "newHash"
		u.PasswordSalt = "newSalt"
		u.Username = "bobsmith"
		u.Devices = model.Devices{"platform": model.Device{Token: "newToken", Platform: "newPlatform"}}
		u.Admin = false
		u.SearchKey = "newSearchKey"

		if err := store.SaveUser(&u); err != nil {
			t.Fatal(err)
		}

		u1, err := store.GetUser(u.ID)
		if err != nil {
			t.Fatal(err)
		}
		compareUsers(t, &u, u1)
	})

	t.Run("get_by_username", func(t *testing.T) {
		u1, err := store.GetUserByUsername(u.Username)
		if err != nil {
			t.Fatal(err)
		}
		compareUsers(t, &u, u1)
	})

	t.Run("get_by_usernames", func(t *testing.T) {
		ux, err := store.GetUsersByUsernames([]string{u.Username})
		if err != nil {
			t.Fatal(err)
		}
		compareUsers(t, &u, ux[0])
	})

	t.Run("create_user_with_nil_id", func(t *testing.T) {
		u = model.User{
			Email:        "kafka@october.news",
			Username:     "kafka",
			DisplayName:  "Franz Kafka",
			PasswordHash: "secret",
			PasswordSalt: "salty secret",
		}
		if err := store.SaveUser(&u); err != nil {
			t.Fatal(err)
		}
		if u.ID == globalid.Nil {
			t.Fatal("expected ID to be generated, was nil")
		}
	})

	t.Run("delete_user", func(t *testing.T) {
		if err := store.DeleteUser(u.ID); err != nil {
			t.Fatal(err)
		}
		_, err := store.GetUser(u.ID)
		err = errors.Cause(err)
		if err != sql.ErrNoRows {
			t.Fatalf("expected error %v but is %v", sql.ErrNoRows, err)
		}
		_, err = store.GetUserByEmail(u.Email)
		err = errors.Cause(err)
		if err != sql.ErrNoRows {
			t.Fatal(err)
		}
		_, err = store.GetUserByUsername(u.Username)
		err = errors.Cause(err)
		if err != sql.ErrNoRows {
			t.Fatal(err)
		}
		users, err := store.GetUsersByUsernames([]string{u.Username})
		if err != nil {
			t.Fatal(err)
		}
		if len(users) != 0 {
			t.Fatal("expected GetUsersByUsernames to return an empty slice")
		}
		users, err = store.GetUsers()
		if err != nil {
			t.Fatal(err)
		}
		for _, user := range users {
			if user.ID == u.ID {
				t.Fatal("expected user to not be returned by GetUsers()")
			}
		}

	})
}

func compareUsers(t *testing.T, exp, got *model.User) {
	t.Helper()
	if exp.ID != got.ID {
		t.Fatalf("unexptected ID.  exp %v, got %v", exp.ID, got.ID)
	}
	if exp.DisplayName != got.DisplayName {
		t.Fatalf("unexptected DisplayName.  exp %v, got %v", exp.DisplayName, got.DisplayName)
	}
	if exp.FirstName != got.FirstName {
		t.Fatalf("unexptected FirstName.  exp %v, got %v", exp.FirstName, got.FirstName)
	}
	if exp.LastName != got.LastName {
		t.Fatalf("unexptected LastName.  exp %v, got %v", exp.LastName, got.LastName)
	}
	if exp.ProfileImagePath != got.ProfileImagePath {
		t.Fatalf("unexptected ProfileImagePath.  exp %v, got %v", exp.ProfileImagePath, got.ProfileImagePath)
	}
	if exp.Bio != got.Bio {
		t.Fatalf("unexptected Bio.  exp %v, got %v", exp.Bio, got.Bio)
	}
	if exp.Email != got.Email {
		t.Fatalf("unexptected Email.  exp %v, got %v", exp.Email, got.Email)
	}
	if exp.PasswordHash != got.PasswordHash {
		t.Fatalf("unexptected PasswordHash.  exp %v, got %v", exp.PasswordHash, got.PasswordHash)
	}
	if exp.PasswordSalt != got.PasswordSalt {
		t.Fatalf("unexptected PasswordSalt.  exp %v, got %v", exp.PasswordSalt, got.PasswordSalt)
	}
	if exp.Username != got.Username {
		t.Fatalf("unexptected Username.  exp %v, got %v", exp.Username, got.Username)
	}
	if exp.Admin != got.Admin {
		t.Fatalf("unexptected Admin.  exp %v, got %v", exp.Admin, got.Admin)
	}
	if exp.SearchKey != got.SearchKey {
		t.Fatalf("unexptected SearchKey.  exp %v, got %v", exp.SearchKey, got.SearchKey)
	}
	if !reflect.DeepEqual(exp.Devices, got.Devices) {
		t.Fatalf("unexpected devices.\nexp:\n%s\ngot:\n%s\n", pretty.Sprint(exp.Devices), pretty.Sprint(got.Devices))
	}
	if exp, got := exp.CreatedAt.Truncate(time.Millisecond), got.CreatedAt.Truncate(time.Millisecond); !exp.Equal(got) {
		t.Fatalf("unexpected CreatedAt.  exp %v, got %v", exp, got)
	}
	if !exp.UpdatedAt.Truncate(time.Millisecond).Before(got.UpdatedAt) && !exp.UpdatedAt.Truncate(time.Millisecond).Equal(got.UpdatedAt.Truncate(time.Millisecond)) {
		t.Fatalf("unexpected UpdatedAt.  exp %v to be before %v", exp.UpdatedAt, got.UpdatedAt)
	}
}
