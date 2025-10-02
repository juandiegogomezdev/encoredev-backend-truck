package store

import (
	"encore.dev/storage/sqldb"
)

type AuthStore struct {
	db *sqldb.Database
}

func NewAuthStore(db *sqldb.Database) *AuthStore {
	if db == nil {
		panic("database is nil")
	}
	return &AuthStore{db: db}
}
