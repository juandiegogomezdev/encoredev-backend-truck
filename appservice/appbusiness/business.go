package appbusiness

import (
	"encore.app/appservice/appstore"
	"encore.app/pkg/myjwt"
	"encore.app/pkg/resendmailer"
)

type BusinessApp struct {
	store     *appstore.StoreApp
	tokenizer myjwt.JWTTokenizer
	mailer    resendmailer.ResendMailer
}

func NewAppBusiness(store *appstore.StoreApp, tokenizer myjwt.JWTTokenizer, mailer resendmailer.ResendMailer) *BusinessApp {
	return &BusinessApp{store: store, tokenizer: tokenizer, mailer: mailer}
}
