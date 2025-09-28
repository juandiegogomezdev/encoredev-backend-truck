package business

import (
	"encore.app/auth/store"
	"encore.app/shared/myjwt"
)

type BusinessAuth struct {
	store     *store.AuthStore
	tokenizer myjwt.JWTTokenizer
}

func NewAuthBusiness(store *store.AuthStore, tokenizer myjwt.JWTTokenizer) *BusinessAuth {
	return &BusinessAuth{store: store, tokenizer: tokenizer}
}
