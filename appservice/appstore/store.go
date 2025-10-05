package appstore

import (
	"encore.dev/storage/sqldb"
	"github.com/jmoiron/sqlx"
)

type StoreApp struct {
	db  *sqldb.Database
	dbx *sqlx.DB
}

func NewStoreApp(db *sqldb.Database, dbx *sqlx.DB) *StoreApp {
	if db == nil {
		panic("database is nil")
	}
	return &StoreApp{db: db, dbx: dbx}
}
