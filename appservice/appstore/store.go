package appstore

import (
	"encore.dev/storage/sqldb"
)

type AppStore struct {
	db *sqldb.Database
}

func NewAppStore(db *sqldb.Database) *AppStore {
	if db == nil {
		panic("database is nil")
	}
	return &AppStore{db: db}
}
