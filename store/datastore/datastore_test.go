package datastore_test

import (
	"sync"
	"testing"

	datastore "github.com/october93/engine/store/datastore"
	"github.com/october93/engine/test"
)

func TestConnectionPool(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping - short test detected")
	}

	cfg := datastore.NewTestConfig()
	cfg.MaxConnections = 1
	cfg.Database = "engine_datastore_connection_pool"
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

	requests := 1000

	var wg sync.WaitGroup
	wg.Add(requests)
	for i := 0; i < requests; i++ {
		go func() {
			_, err := store.GetCards()
			if err != nil {
				t.Errorf("GetCards(): %v", err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
