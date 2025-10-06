package appService

import (
	"context"
	"net/http"
	"strings"

	"encore.app/appservice/appbusiness"
	"encore.app/appservice/appstore"
	"encore.app/pkg/myjwt"
	"encore.app/pkg/resendmailer"
	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"encore.dev/config"
	"encore.dev/storage/sqldb"
	"encore.dev/types/uuid"
	"github.com/jmoiron/sqlx"
)

var secrets struct {
	JWT_SECRET_KEY string
	RESEND_API_KEY string
}

var appDB = sqldb.NewDatabase("db_app", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})

type ServiceConfig struct {
	BaseUrl string
}

var cfg *ServiceConfig = config.Load[*ServiceConfig]()

//encore:service
type ServiceApp struct {
	b *appbusiness.BusinessApp
}

func initServiceApp() (*ServiceApp, error) {
	appDBX := sqlx.NewDb(appDB.Stdlib(), "postgres")
	tokenizer := myjwt.NewJWTTokenizer(secrets.JWT_SECRET_KEY)

	// Initialize the resend mailer
	m := resendmailer.NewResendMailer(secrets.RESEND_API_KEY, "Acme <onboarding@resend.dev>")
	s := appstore.NewStoreApp(appDB, appDBX)
	b := appbusiness.NewAppBusiness(s, tokenizer, m)
	return &ServiceApp{b: b}, nil
}

type MyAuthParams struct {
	// Extract the auth token from either the cookie or the Authorization header
	SessionCookie *http.Cookie `cookie:"auth_token"`
	// Extract the authorization header
	AuthorizationHeader string `header:"Authorization"`
}

//encore:authhandler
func (s *ServiceApp) AuthHandler(ctx context.Context, p *MyAuthParams) (auth.UID, *AuthData, error) {
	// Extract the token from the request
	token, err := extractToken(p)
	if err != nil {
		return "", nil, err
	}

	// Parse the token and get the claims
	claims, err := s.b.ParseMembershipToken(token)
	if err != nil {
		return "", nil, err
	}

	authData := &AuthData{
		UserID:       claims.UserID,
		SessionID:    claims.SessionID,
		MembershipID: claims.MembershipID,
	}
	return auth.UID(claims.UserID.String()), authData, nil
}

type AuthData struct {
	UserID       uuid.UUID
	SessionID    uuid.UUID
	MembershipID uuid.UUID
}

// Extract token from AuthParams
func extractToken(p *MyAuthParams) (token string, err error) {
	// Verify if the token is in the cookie
	if p.SessionCookie != nil {
		return p.SessionCookie.Value, nil
	}

	// Verify if the token is in the Authorization header
	if p.AuthorizationHeader != "" {
		if after, found := strings.CutPrefix(p.AuthorizationHeader, "Bearer "); found {
			token = strings.TrimSpace(after)
			if token != "" {
				return token, nil
			}
		}
		return "", &errs.Error{
			Code:    errs.Unauthenticated,
			Message: "Formato de header de autorizacion invalido",
		}
	}

	return "", &errs.Error{
		Code:    errs.Unauthenticated,
		Message: "se requiere autenticacion: proporciona una cookie o un header de autorizacion",
	}
}
