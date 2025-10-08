package appbusiness

import (
	"encore.app/appservice/appstore"
	"encore.app/appservice/shared"
	"encore.app/pkg/myjwt"
	"encore.app/pkg/resendmailer"
	"encore.dev/beta/errs"
)

type BusinessApp struct {
	store     *appstore.StoreApp
	tokenizer myjwt.JWTTokenizer
	mailer    resendmailer.ResendMailer
}

func NewAppBusiness(store *appstore.StoreApp, tokenizer myjwt.JWTTokenizer, mailer resendmailer.ResendMailer) *BusinessApp {
	return &BusinessApp{store: store, tokenizer: tokenizer, mailer: mailer}
}

func (b *BusinessApp) ParseMembershipToken(token string) (*myjwt.MembershipClaims, error) {
	claims, status := b.tokenizer.ParseMembershipToken(token)
	switch status {
	case myjwt.TokenStatusExpired:
		return nil, &errs.Error{
			Code:    errs.Unauthenticated,
			Message: "El token ha expirado",
			Details: shared.ErrorDetailsToken{TokenStatus: string(myjwt.TokenStatusExpired)},
		}
	case myjwt.TokenStatusInvalid:
		return nil, &errs.Error{
			Code:    errs.Unauthenticated,
			Message: "El token no es válido",
			Details: shared.ErrorDetailsToken{TokenStatus: string(myjwt.TokenStatusInvalid)},
		}
	case myjwt.TokenStatusValid:
		return claims, nil
	}
	return nil, &errs.Error{
		Code:    errs.Unauthenticated,
		Message: "Error al validar el token",
		Details: shared.ErrorDetailsToken{TokenStatus: string(myjwt.TokenStatusInvalid)},
	}
}
