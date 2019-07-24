package datastore

import "github.com/jmoiron/sqlx"

type Tx struct {
	*sqlx.Tx
}
