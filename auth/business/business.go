package business

import (
	"encore.app/auth/store"
	"encore.app/pkg/myjwt"
	"encore.app/pkg/resendmailer"
)

type BusinessAuth struct {
	store     *store.AuthStore
	tokenizer myjwt.JWTTokenizer
	mailer    resendmailer.ResendMailer
}

func NewAuthBusiness(store *store.AuthStore, tokenizer myjwt.JWTTokenizer, mailer resendmailer.ResendMailer) *BusinessAuth {
	return &BusinessAuth{store: store, tokenizer: tokenizer, mailer: mailer}
}
