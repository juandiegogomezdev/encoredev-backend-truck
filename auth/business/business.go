package business

import "encore.app/auth/store"

type BusinessAuth struct {
	store *store.AuthStore
}

func NewAuthBusiness(store *store.AuthStore) *BusinessAuth {
	return &BusinessAuth{store: store}
}
