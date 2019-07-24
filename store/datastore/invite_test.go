package datastore_test

import (
	"testing"
	"time"

	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
	datastore "github.com/october93/engine/store/datastore"
	"github.com/october93/engine/test"
)

func Test_Invites(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping - short test detected")
	}

	cfg := datastore.NewTestConfig()
	cfg.Database = "engine_datastore_invites_test"
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

	invite := *model.NewInviteWithParams(
		globalid.Next(),
		globalid.Next(),
		"token",
		time.Date(2010, time.November, 10, 23, 0, 0, 0, time.UTC),
		time.Date(2011, time.November, 10, 23, 0, 0, 0, time.UTC),
	)
	user := model.User{
		ID:           invite.NodeID,
		Email:        "kafka@october.news",
		Username:     "kafka",
		DisplayName:  "Kafka",
		PasswordHash: "secret",
		PasswordSalt: "salty secret",
	}

	t.Run("create_invite", func(t *testing.T) {
		if err := store.SaveUser(&user); err != nil {
			t.Fatal(err)
		}
		if err := store.SaveInvite(&invite); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("get_invites", func(t *testing.T) {
		invites, err := store.GetInvites()
		if err != nil {
			t.Fatal(err)
		}
		if exp, got := 1, len(invites); exp != got {
			t.Fatalf("unexpected length of invites.  exp %d, got %d", exp, got)
		}
		exp := &invite
		got := invites[0]
		compareInvites(t, exp, got)
	})

	t.Run("get_invite_by_token", func(t *testing.T) {
		it, err := store.GetInviteByToken(invite.Token)
		if err != nil {
			t.Fatal(err)
		}
		compareInvites(t, &invite, it)
	})

	t.Run("update_invite", func(t *testing.T) {
		invite.Token = "new_token"
		invite.RemainingUses = 1
		if err := store.SaveInvite(&invite); err != nil {
			t.Fatal(err)
		}
		it, err := store.GetInviteByToken(invite.Token)
		if err != nil {
			t.Fatal(err)
		}
		compareInvites(t, &invite, it)
	})

	t.Run("delete_invite", func(t *testing.T) {
		err := store.DeleteInvite(invite.ID)
		if err != nil {
			t.Fatal(err)
		}
		invites, err := store.GetInvites()
		if err != nil {
			t.Fatal(err)
		}
		if exp, got := 0, len(invites); exp != got {
			t.Fatalf("unexpected length of invites.  exp %d, got %d", exp, got)
		}
	})

	t.Run("create_empty_invite", func(t *testing.T) {
		invite := model.Invite{
			NodeID: globalid.Next(),
			Token:  "12345=",
		}
		user := model.User{
			ID:           invite.NodeID,
			Email:        "franz@october.news",
			Username:     "franz",
			DisplayName:  "Franz",
			PasswordHash: "secret",
			PasswordSalt: "salty secret",
		}
		if err := store.SaveUser(&user); err != nil {
			t.Fatal(err)
		}
		if err := store.SaveInvite(&invite); err != nil {
			t.Fatal(err)
		}
	})
}

func compareInvites(t *testing.T, exp, got *model.Invite) {
	t.Helper()
	if exp.ID != got.ID {
		t.Fatalf("unexpected ID.  exp %v, got %v", exp.ID, got.ID)
	}
	if exp.NodeID != got.NodeID {
		t.Fatalf("unexpected NodeID.  exp %v, got %v", exp.NodeID, got.NodeID)
	}
	if exp.Token != got.Token {
		t.Fatalf("unexpected Token.  exp %v, got %v", exp.Token, got.Token)
	}
	if exp.RemainingUses != got.RemainingUses {
		t.Fatalf("unexpected Redeemed.  exp %v, got %v", exp.RemainingUses, got.RemainingUses)
	}
	if exp, got := exp.CreatedAt.Truncate(time.Millisecond), got.CreatedAt.Truncate(time.Millisecond); !exp.Equal(got) {
		t.Fatalf("unexpected CreatedAt.  exp %v, got %v", exp, got)
	}
	if !exp.UpdatedAt.Truncate(time.Millisecond).Before(got.UpdatedAt) && !exp.UpdatedAt.Truncate(time.Millisecond).Equal(got.UpdatedAt.Truncate(time.Millisecond)) {
		t.Fatalf("unexpected UpdatedAt.  exp %v to be before %v", exp.UpdatedAt, got.UpdatedAt)
	}
}
