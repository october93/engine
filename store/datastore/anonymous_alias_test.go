package datastore_test

import (
	"reflect"
	"testing"

	"github.com/kr/pretty"
	"github.com/october93/engine/kit/globalid"
	model "github.com/october93/engine/model"
	datastore "github.com/october93/engine/store/datastore"
	"github.com/october93/engine/test"
)

func TestAnonymousAliases(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping - short test detected")
	}

	cfg := datastore.NewTestConfig()
	cfg.Database = "engine_datastore_anonymous_aliases_test"
	db := test.DBInit(t, cfg)
	// drop the database after the test is finished
	defer test.DBCleanup(t, db)

	store, err := datastore.New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	// close the store at the end
	defer func() {
		if e := store.Close(); e != nil {
			t.Fatalf("store.Close() failed.: %s", e)
		}
	}()

	user1 := &model.User{
		ID:          "0f673ea0-999f-40e4-9b1f-521f678ae7cf",
		Username:    "chad",
		Email:       "chad@october.news",
		DisplayName: "Chad Unicorn",
	}
	user2 := &model.User{
		ID:          "2226a76d-77d3-48cb-9bba-136ad7b87493",
		Username:    "Richard",
		Email:       "richard@october.news",
		DisplayName: "Richard Hendricks",
	}

	card1 := &model.Card{
		ID:            "6b34ec36-899a-4dc8-874f-4f8e191216c0",
		OwnerID:       user1.ID,
		AuthorToAlias: make(model.IdentityMap),
	}
	card2 := &model.Card{
		ID:      "290278ec-bcee-4b62-8ca4-c12f9fc8a337",
		OwnerID: user1.ID,
	}
	card3 := &model.Card{
		ID:      "9495e377-f4d1-48c1-ae23-7a5e90ab319b",
		OwnerID: user2.ID,
	}

	reaction1 := &model.UserReaction{
		CardID: card1.ID,
		UserID: user1.ID,
		Type:   model.ReactionLike,
	}
	reaction2 := &model.UserReaction{
		CardID: card1.ID,
		UserID: user2.ID,
		Type:   model.ReactionLike,
	}
	reaction3 := &model.UserReaction{
		CardID: card3.ID,
		UserID: user2.ID,
		Type:   model.ReactionLike,
	}

	setupUsers(t, store, user1, user2)
	setupCards(t, store, card1, card2)

	var aliases []*model.AnonymousAlias
	t.Run("generate_ids", func(t *testing.T) {
		err = store.CreateAnonymousAliases()
		if err != nil {
			t.Fatalf("GenerateIDs(): unexpected error %v", err)
		}
		aliases, err = store.GetAnonymousAliases()
		if err != nil {
			t.Fatalf("GetAnonymousAliases(): unexpected error %v", err)
		}
		if len(aliases) == 0 {
			t.Fatalf("GenerateIDs(): expected to generate and persist anonymous aliases")
		}
	})

	t.Run("get_unused_alias", func(t *testing.T) {
		var alias *model.AnonymousAlias
		alias, err = store.GetUnusedAlias(card1.ID)
		if err != nil {
			t.Fatalf("GetUnusedAlias(): unexpected error: %v", err)
		}
		if alias == nil {
			t.Fatalf("GetUnusedAlias(): did not return an alias")
		}
	})

	t.Run("get_unused_alias_with_used_alias", func(t *testing.T) {
		aliases, err = store.GetAnonymousAliases()
		if err != nil {
			t.Fatalf("GetAnonymousAliases(): unexpected error %v", err)
		}
		ids := make([]globalid.ID, len(aliases))
		for i, alias := range aliases {
			ids[i] = alias.ID
		}
		card1.AuthorToAlias[user1.ID] = aliases[3].ID
		err = store.SaveCard(card1)
		if err != nil {
			t.Fatal(err)
		}
		var alias *model.AnonymousAlias
		alias, err = store.GetUnusedAlias(card1.ID)
		if err != nil {
			t.Fatalf("GetUnusedAlias(): unexpected error: %v", err)
		}
		if alias == nil {
			t.Fatalf("GetUnusedAlias(): did not return an alias")
		}
	})

	t.Run("get_unused_alias_with_all_but_one_used_alias", func(t *testing.T) {
		aliases, err = store.GetAnonymousAliases()
		if err != nil {
			t.Fatalf("GetAnonymousAliases(): unexpected error %v", err)
		}
		ids := make([]globalid.ID, len(aliases))
		for i, alias := range aliases {
			ids[i] = alias.ID
		}
		// all aliases are used except the first
		for i := 1; i < len(aliases); i++ {
			card1.AuthorToAlias[globalid.Next()] = aliases[i].ID
		}
		err = store.SaveCard(card1)
		if err != nil {
			t.Fatal(err)
		}
		var alias *model.AnonymousAlias
		alias, err = store.GetUnusedAlias(card1.ID)
		if err != nil {
			t.Fatalf("GetUnusedAlias(): unexpected error: %v", err)
		}
		compareAnonymousAlias(t, aliases[0], alias)
	})
	t.Run("get_unused_alias_with_no_alias_left", func(t *testing.T) {
		aliases, err = store.GetAnonymousAliases()
		if err != nil {
			t.Fatalf("GetAnonymousAliases(): unexpected error %v", err)
		}
		for _, alias := range aliases {
			card1.AuthorToAlias[globalid.Next()] = alias.ID
		}
		err = store.SaveCard(card1)
		if err != nil {
			t.Fatal(err)
		}
		_, err = store.GetUnusedAlias(card1.ID)
		if err != datastore.ErrNoAnonymousAliasLeft {
			t.Fatalf("GetUnusedAlias(): expected error %v, actual %v", datastore.ErrNoAnonymousAliasLeft, err)
		}
	})

	t.Run("anonymous alias last used", func(t *testing.T) {
		// user 1 just posted card with his real identity
		lastUsed, err := store.GetAnonymousAliasLastUsed(user1.ID, card1.ID)
		if err != nil {
			t.Fatal(err)
		}
		if lastUsed {
			t.Fatalf("expected %t, actual %t", false, lastUsed)
		}
		reaction1.AliasID = aliases[0].ID
		err = store.SaveUserReaction(reaction1)
		if err != nil {
			t.Fatal(err)
		}
		// user 1 reacted to his own post anonymously
		lastUsed, err = store.GetAnonymousAliasLastUsed(user1.ID, card1.ID)
		if err != nil {
			t.Fatal(err)
		}
		if !lastUsed {
			t.Fatalf("expected %t, actual %t", !lastUsed, lastUsed)
		}
		// user 2 did not engage at all with this thread yet
		lastUsed, err = store.GetAnonymousAliasLastUsed(user2.ID, card1.ID)
		if err != nil {
			t.Fatal(err)
		}
		if lastUsed {
			t.Fatalf("expected %t, actual %t", !lastUsed, lastUsed)
		}
		reaction2.AliasID = aliases[1].ID
		err = store.SaveUserReaction(reaction2)
		if err != nil {
			t.Fatal(err)
		}
		// user 2 reacted anonymously
		lastUsed, err = store.GetAnonymousAliasLastUsed(user2.ID, card1.ID)
		if err != nil {
			t.Fatal(err)
		}
		if !lastUsed {
			t.Fatalf("expected %t, actual %t", !lastUsed, lastUsed)
		}
		// user 2 did not engage at all with this thread yet
		lastUsed, err = store.GetAnonymousAliasLastUsed(user2.ID, card2.ID)
		if err != nil {
			t.Fatal(err)
		}
		if lastUsed {
			t.Fatalf("expected %t, actual %t", !lastUsed, lastUsed)
		}
		// user 2 replies to card 1 with real identity
		card3.ReplyTo(card1)
		err = store.SaveCard(card3)
		if err != nil {
			t.Fatal(err)
		}
		lastUsed, err = store.GetAnonymousAliasLastUsed(user2.ID, card1.ID)
		if err != nil {
			t.Fatal(err)
		}
		if lastUsed {
			t.Fatalf("expected %t, actual %t", !lastUsed, lastUsed)
		}
		// user 2 reacts to own card in thread anonymously
		reaction3.AliasID = aliases[1].ID
		err = store.SaveUserReaction(reaction3)
		if err != nil {
			t.Fatal(err)
		}
		lastUsed, err = store.GetAnonymousAliasLastUsed(user2.ID, card1.ID)
		if err != nil {
			t.Fatal(err)
		}
		if !lastUsed {
			t.Fatalf("expected %t, actual %t", !lastUsed, lastUsed)
		}
	})
}

func compareAnonymousAlias(t *testing.T, exp, got *model.AnonymousAlias) {
	t.Helper()
	if !reflect.DeepEqual(exp, got) {
		t.Fatalf("unexpected anonymous alias, diff: %v", pretty.Diff(exp, got))
	}
}
