package datastore_test

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/kr/pretty"
	"github.com/october93/engine/kit/globalid"
	model "github.com/october93/engine/model"
	datastore "github.com/october93/engine/store/datastore"
	"github.com/october93/engine/test"
	"github.com/pkg/errors"
)

func Test_Cards(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping - short test detected")
	}

	cfg := datastore.NewTestConfig()
	cfg.Database = "engine_datastore_cards_test"
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

	c := model.Card{
		OwnerID:             globalid.Next(),
		ID:                  globalid.Next(),
		Title:               "This is the title",
		Content:             "This is the content",
		URL:                 "http://someurl.com",
		BackgroundColor:     "indigo-blue",
		BackgroundImagePath: "https://october.news/images/cards/c06f34d5-92fe-4cd5-9c4a-1de69fed5b25.png",
		Anonymous:           true,
		AuthorToAlias:       model.IdentityMap{globalid.Next(): globalid.Next()},
	}

	u := model.User{
		ID:           c.OwnerID,
		Email:        "kafka@october.news",
		Username:     "kafka",
		DisplayName:  "Kafka",
		PasswordHash: "secret",
		PasswordSalt: "salty secret",
	}

	t.Run("create_card", func(t *testing.T) {
		if err := store.SaveUser(&u); err != nil {
			t.Fatal(err)
		}
		if err := store.SaveCard(&c); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("get_cards", func(t *testing.T) {
		cs, err := store.GetCards()
		if err != nil {
			t.Fatal(err)
		}
		if exp, got := 1, len(cs); exp != got {
			t.Fatalf("unexpected length of cards.  exp %d, got %d", exp, got)
		}
		exp := &c
		got := cs[0]
		compareCards(t, exp, got)
	})

	t.Run("get_card", func(t *testing.T) {
		cg, err := store.GetCard(c.ID)
		if err != nil {
			t.Fatal(err)
		}
		compareCards(t, &c, cg)
	})

	t.Run("get_thread", func(t *testing.T) {
		r1 := model.Card{
			OwnerID:       u.ID,
			ThreadRootID:  c.ID,
			ThreadReplyID: c.ID,
		}
		if err := store.SaveCard(&r1); err != nil {
			t.Fatal(err)
		}
		r2 := model.Card{
			OwnerID:       u.ID,
			ThreadRootID:  c.ID,
			ThreadReplyID: c.ID,
			CreatedAt:     time.Now().UTC().Add(10 * time.Hour),
		}
		if err := store.SaveCard(&r2); err != nil {
			t.Fatal(err)
		}
		r4 := model.Card{
			OwnerID:       u.ID,
			ThreadRootID:  c.ID,
			ThreadReplyID: r1.ID,
			CreatedAt:     time.Now().UTC().Add(3 * time.Hour),
		}
		if err := store.SaveCard(&r4); err != nil {
			t.Fatal(err)
		}
		r3 := model.Card{
			OwnerID:       u.ID,
			ThreadRootID:  c.ID,
			ThreadReplyID: r1.ID,
			CreatedAt:     time.Now().UTC().Add(2 * time.Hour),
		}
		if err := store.SaveCard(&r3); err != nil {
			t.Fatal(err)
		}
		cgs, err := store.GetThread(c.ID, globalid.Nil)
		if err != nil {
			t.Fatal(err)
		}
		if exp, got := 4, len(cgs); exp != got {
			t.Fatalf("unexpected length of cards.  exp %d, got %d", exp, got)
		}
		compareCards(t, &r1, cgs[0])
		compareCards(t, &r3, cgs[1])
		compareCards(t, &r4, cgs[2])
		compareCards(t, &r2, cgs[3])
	})

	t.Run("get_thread_of_reply", func(t *testing.T) {
		thread, err := store.GetThread(c.ID, globalid.Nil)
		if err != nil {
			t.Fatalf("GetThread(): %v", err)
		}
		thread, err = store.GetThread(thread[0].ID, globalid.Nil)
		if err != nil {
			t.Fatalf("GetThread(): %v", err)
		}
		if exp, got := 0, len(thread); exp != got {
			t.Fatalf("unexpected length of cards.  exp %d, got %d", exp, got)
		}
	})

	t.Run("get_thread_out_of_two", func(t *testing.T) {
		d := model.Card{
			OwnerID: u.ID,
		}
		if err := store.SaveCard(&d); err != nil {
			t.Fatalf("SaveCard(): %v", err)
		}
		r1 := model.Card{
			OwnerID:       u.ID,
			ThreadRootID:  d.ID,
			ThreadReplyID: d.ID,
		}
		if err := store.SaveCard(&r1); err != nil {
			t.Fatalf("SaveCard(): %v", err)
		}
		thread, err := store.GetThread(d.ID, globalid.Nil)
		if err != nil {
			t.Fatalf("GetThread(): %v", err)
		}
		if exp, got := 1, len(thread); exp != got {
			t.Fatalf("unexpected length of cards.  exp %d, got %d", exp, got)
		}
	})

	t.Run("update_card", func(t *testing.T) {
		c.Title = "New Title"
		c.Content = "New Content"
		c.URL = "http://newurl.com"
		c.BackgroundColor = "indigo-blue"
		c.BackgroundImagePath = "https://october.news/images/cards/c06f34d5-92fe-4cd5-9c4a-1de69fed5b25.png"
		c.Anonymous = false

		if err := store.SaveCard(&c); err != nil {
			t.Fatal(err)
		}

		c1, err := store.GetCard(c.ID)
		if err != nil {
			t.Fatal(err)
		}
		compareCards(t, &c, c1)
	})

	t.Run("get_cards_by_interval", func(t *testing.T) {
		// Clear the database
		if err := db.TruncateAll(); err != nil {
			t.Fatalf("trancating db failed: %s", err)
		}
		c1 := model.Card{
			OwnerID:             globalid.Next(),
			Title:               "This is the title",
			Content:             "This is the content",
			URL:                 "http://someurl.com",
			BackgroundColor:     "indigo-blue",
			BackgroundImagePath: "https://october.news/images/cards/c06f34d5-92fe-4cd5-9c4a-1de69fed5b25.png",
			Anonymous:           true,
			CreatedAt:           time.Date(2008, time.November, 10, 23, 0, 0, 0, time.UTC),
		}
		c2 := model.Card{
			OwnerID:             globalid.Next(),
			Title:               "title",
			Content:             "content",
			URL:                 "http://yahoo.com",
			BackgroundColor:     "red-pink",
			BackgroundImagePath: "https://october.news/images/cards/f1ba6683-ee27-4195-a161-567a7e9077de.png",
			Anonymous:           false,
			CreatedAt:           time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		}
		c1.ThreadRootID = globalid.Next()
		c1.ThreadReplyID = c1.ThreadRootID
		c2.ThreadRootID = globalid.Next()
		c2.ThreadReplyID = c2.ThreadRootID

		u1 := model.User{
			ID:           c1.OwnerID,
			Email:        "kafka@october.news",
			Username:     "kafka",
			DisplayName:  "Kafka",
			PasswordHash: "secret",
			PasswordSalt: "salty secret",
		}
		u2 := model.User{
			ID:           c2.OwnerID,
			Email:        "chad@october.news",
			Username:     "chad",
			DisplayName:  "Chad",
			PasswordHash: "secret",
			PasswordSalt: "salty secret",
		}
		r1 := model.Card{
			ID:      c1.ThreadRootID,
			OwnerID: u1.ID,
		}
		r2 := model.Card{
			ID:      c2.ThreadRootID,
			OwnerID: u2.ID,
		}
		if err := store.SaveUser(&u1); err != nil {
			t.Fatal(err)
		}
		if err := store.SaveUser(&u2); err != nil {
			t.Fatal(err)
		}
		if err := store.SaveCard(&r1); err != nil {
			t.Fatal(err)
		}
		if err := store.SaveCard(&r2); err != nil {
			t.Fatal(err)
		}
		if err := store.SaveCard(&c1); err != nil {
			t.Fatal(err)
		}
		if err := store.SaveCard(&c2); err != nil {
			t.Fatal(err)
		}
		from := time.Date(2008, time.November, 10, 23, 0, 0, 0, time.UTC)
		to := time.Date(2008, time.November, 10, 23, 0, 0, 0, time.UTC)
		cgs, err := store.GetCardsByInterval(from, to)
		if err != nil {
			t.Fatal(err)
		}
		if exp, got := 1, len(cgs); exp != got {
			t.Fatalf("unexpected length of cards.  exp %d, got %d", exp, got)
		}
		compareCards(t, &c1, cgs[0])
		from = time.Date(2000, time.November, 10, 23, 0, 0, 0, time.UTC)
		to = time.Date(2010, time.November, 10, 23, 0, 0, 0, time.UTC)
		cgs, err = store.GetCardsByInterval(from, to)
		if err != nil {
			t.Fatal(err)
		}
		if exp, got := 2, len(cgs); exp != got {
			t.Fatalf("unexpected length of cards.  exp %d, got %d", exp, got)
		}
		// Cards come back sorted DESC
		got := cgs[1]
		compareCards(t, &c1, got)
		got = cgs[0]
		compareCards(t, &c2, got)
	})

	t.Run("get_cards_by_node_in_range", func(t *testing.T) {
		// Clear the database
		if err := db.TruncateAll(); err != nil {
			t.Fatalf("trancating db failed: %s", err)
		}
		c1 := model.Card{
			OwnerID:             globalid.Next(),
			Title:               "This is the title",
			Content:             "This is the content",
			URL:                 "http://someurl.com",
			BackgroundColor:     "indigo-blue",
			BackgroundImagePath: "https://october.news/images/cards/c06f34d5-92fe-4cd5-9c4a-1de69fed5b25.png",
			Anonymous:           true,
			CreatedAt:           time.Date(2008, time.November, 10, 23, 0, 0, 0, time.UTC),
		}
		c2 := model.Card{
			OwnerID:             globalid.Next(),
			Title:               "title",
			Content:             "content",
			URL:                 "http://yahoo.com",
			BackgroundColor:     "red-pink",
			BackgroundImagePath: "https://october.news/images/cards/f1ba6683-ee27-4195-a161-567a7e9077de.png",
			Anonymous:           false,
			CreatedAt:           time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		}
		c1.ThreadRootID = globalid.Next()
		c1.ThreadReplyID = c1.ThreadRootID
		c2.ThreadRootID = globalid.Next()
		c2.ThreadReplyID = c2.ThreadRootID

		u1 := model.User{
			ID:           c1.OwnerID,
			Email:        "kafka@october.news",
			Username:     "kafka",
			DisplayName:  "Kafka",
			PasswordHash: "secret",
			PasswordSalt: "salty secret",
		}
		u2 := model.User{
			ID:           c2.OwnerID,
			Email:        "franz@october.news",
			Username:     "franz",
			DisplayName:  "Franz",
			PasswordHash: "secret",
			PasswordSalt: "salty secret",
		}
		r1 := model.Card{
			ID:      c1.ThreadRootID,
			OwnerID: u1.ID,
		}
		r2 := model.Card{
			ID:      c2.ThreadRootID,
			OwnerID: u2.ID,
		}
		if err := store.SaveUser(&u1); err != nil {
			t.Fatal(err)
		}
		if err := store.SaveUser(&u2); err != nil {
			t.Fatal(err)
		}
		if err := store.SaveCard(&r1); err != nil {
			t.Fatal(err)
		}
		if err := store.SaveCard(&r2); err != nil {
			t.Fatal(err)
		}
		if err := store.SaveCard(&c1); err != nil {
			t.Fatal(err)
		}
		if err := store.SaveCard(&c2); err != nil {
			t.Fatal(err)
		}
		from := time.Date(2008, time.November, 10, 23, 0, 0, 0, time.UTC)
		to := time.Date(2008, time.November, 10, 23, 0, 0, 0, time.UTC)
		cgs, err := store.GetCardsByNodeInRange(c1.OwnerID, from, to)
		if err != nil {
			t.Fatal(err)
		}
		if exp, got := 1, len(cgs); exp != got {
			t.Fatalf("unexpected length of cards.  exp %d, got %d", exp, got)
		}
		got := cgs[0]
		compareCards(t, &c1, got)
	})

	t.Run("delete_card", func(t *testing.T) {
		err := store.DeleteCard(c.ID)
		if err != nil {
			t.Fatalf("DeleteCard(): %v", err)
		}
		_, err = store.GetCard(c.ID)
		err = errors.Cause(err)
		if err != sql.ErrNoRows {
			t.Fatalf("GetCard(): expected error %v, actual: %v", sql.ErrNoRows, err)
		}
	})
}

func compareCards(t *testing.T, exp, got *model.Card) {
	t.Helper()
	if exp.ID != got.ID {
		t.Fatalf("unexpected ID.  exp %v, got %v", exp.ID, got.ID)
	}
	if exp.OwnerID != got.OwnerID {
		t.Fatalf("unexpected OwnerID.  exp %v, got %v", exp.OwnerID, got.OwnerID)
	}

	if exp.ThreadReplyID != got.ThreadReplyID {
		t.Fatalf("unexpected ThreadReplyID.  exp %v, got %v", exp.ThreadReplyID, got.ThreadReplyID)
	}
	if exp.ThreadRootID != got.ThreadRootID {
		t.Fatalf("unexpected ThreadRootID.  exp %v, got %v", exp.ThreadRootID, got.ThreadRootID)
	}
	if exp.Title != got.Title {
		t.Fatalf("unexpected Title.  exp %v, got %v", exp.Title, got.Title)
	}
	if exp.Content != got.Content {
		t.Fatalf("unexpected Content.  exp %v, got %v", exp.Content, got.Content)
	}
	if exp.URL != got.URL {
		t.Fatalf("unexpected URL.  exp %v, got %v", exp.URL, got.URL)
	}
	if exp.BackgroundColor != got.BackgroundColor {
		t.Fatalf("unexpected BackgroundColor.  exp %v, got %v", exp.BackgroundColor, got.BackgroundColor)
	}
	if exp.BackgroundImagePath != got.BackgroundImagePath {
		t.Fatalf("unexpected BackgroundImagePath.  exp %v, got %v", exp.BackgroundImagePath, got.BackgroundImagePath)
	}
	if exp.Anonymous != got.Anonymous {
		t.Fatalf("unexpected Anonymous.  exp %v, got %v", exp.Anonymous, got.Anonymous)
	}
	if exp, got := exp.CreatedAt.Truncate(time.Millisecond), got.CreatedAt.Truncate(time.Millisecond); !exp.Equal(got) {
		t.Fatalf("unexpected CreatedAt.  exp %v, got %v", exp, got)
	}
	if !exp.UpdatedAt.Truncate(time.Millisecond).Before(got.UpdatedAt) && !exp.UpdatedAt.Truncate(time.Millisecond).Equal(got.UpdatedAt.Truncate(time.Millisecond)) {
		t.Fatalf("unexpected UpdatedAt.  exp %v to be before %v", exp.UpdatedAt, got.UpdatedAt)
	}
	if !reflect.DeepEqual(exp.AuthorToAlias, got.AuthorToAlias) {
		t.Fatalf("unexpected AuthorToAlias.  diff: %v", pretty.Diff(exp.AuthorToAlias, got.AuthorToAlias))
	}
}
