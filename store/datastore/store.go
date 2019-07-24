package datastore

import (
	"fmt"
	"os"

	"github.com/gobuffalo/pop"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Store struct {
	*sqlx.DB
	config Config
}

// New returns a new instance of Store. It connects to the default database
// postgres to check whether the application database exists. If it does not,
// it creates the database and uses the current schema.
func New(c Config) (*Store, error) {
	store := &Store{config: c}
	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable", c.User, c.Password, c.Database, c.Host, c.Port))
	if err != nil {
		return nil, err
	}
	db.DB.SetMaxOpenConns(c.MaxConnections)
	store.DB = db
	return store, nil
}

func (store *Store) MustBegin() *Tx {
	tx := store.DB.MustBegin()
	return &Tx{Tx: tx}
}

func (store *Store) Close() error {
	return store.DB.Close()
}

// DropDatabase deletes the existing database.
func DropDatabase(c Config) error {
	db, err := pop.Connect(c.Environment)
	if err != nil {
		return err
	}
	exists, err := exists(c)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	return db.Dialect.DropDB()
}

func SetupDatabase(c Config) error {
	db, err := pop.Connect(c.Environment)
	if err != nil {
		return err
	}
	exists, err := exists(c)
	if err != nil {
		return err
	}
	if !exists {
		err = db.Dialect.CreateDB()
		if err != nil {
			return err
		}
		var file *os.File
		file, err = os.Open("migrations/base.sql")
		if err != nil {
			return err
		}
		err = db.Dialect.LoadSchema(file)
		if err != nil {
			return err
		}
		err = file.Close()
		if err != nil {
			return err
		}
	}
	mig, err := pop.NewFileMigrator(c.MigrationPath, db)
	if err != nil {
		return err
	}
	err = mig.Up()
	if err != nil {
		return err
	}
	return db.Close()
}

func exists(c Config) (bool, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname=postgres host=%s port=%d sslmode=disable", c.User, c.Password, c.Host, c.Port))
	if err != nil {
		return false, err
	}
	existsQuery := `SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = $1)`
	var exists bool
	err = db.Get(&exists, existsQuery, c.Database)
	if err != nil {
		return false, err
	}
	return exists, db.Close()
}
