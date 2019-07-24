package test

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/logging"
	datastore "github.com/october93/engine/store/datastore"
)

var muteLogs = func(lvl logging.Level, s string, args ...interface{}) {}

func init() {
	pop.SetLogger(muteLogs)
}

func DBInit(t *testing.T, cfg datastore.Config) *pop.Connection {
	t.Helper()
	db, err := pop.NewConnection(&pop.ConnectionDetails{
		Dialect:  "postgres",
		Port:     strconv.Itoa(cfg.Port),
		Host:     cfg.Host,
		Database: cfg.Database,
		User:     cfg.User,
		Password: cfg.Password,
	})
	if err != nil {
		t.Fatal(err)
	}
	// drop the test db. ignore error as it may not exist.
	if e := db.Dialect.DropDB(); e != nil {
		t.Logf("initial drop database failed, should not exist anyway so this is expected: %s", e)
	}

	// create the test db:
	if err = db.Dialect.CreateDB(); err != nil {
		t.Errorf("error creating db: %s", err)
	}

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	projectpath := strings.Replace(basepath, "/test", "", -1)
	err = os.Chdir(projectpath)
	if err != nil {
		t.Fatal(err)
	}

	var file *os.File
	file, err = os.Open("migrations/base.sql")
	if err != nil {
		t.Fatal(err)
	}
	err = db.Dialect.LoadSchema(file)
	if err != nil {
		t.Fatal(err)
	}
	err = file.Close()
	if err != nil {
		t.Fatal(err)
	}
	mig, err := pop.NewFileMigrator("migrations", db)
	if err != nil {
		t.Fatal(err)
	}
	err = mig.Up()
	if err != nil {
		t.Fatal(err)
	}

	err = mig.Connection.Close()
	if err != nil {
		t.Fatal(err)
	}

	db, err = pop.NewConnection(&pop.ConnectionDetails{
		Dialect:  "postgres",
		Port:     strconv.Itoa(cfg.Port),
		Host:     cfg.Host,
		Database: cfg.Database,
		User:     cfg.User,
		Password: cfg.Password,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Open(); err != nil {
		t.Fatal(err)
	}

	return db
}

func DBCleanup(t *testing.T, conn *pop.Connection) {
	t.Helper()
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
	if err := conn.Dialect.DropDB(); err != nil {
		t.Fatal(err)
	}
}
