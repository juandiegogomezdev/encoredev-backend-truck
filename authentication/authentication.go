// Service auth implements user login and register functionality.
package authentication

import (
	"fmt"
	"net/http"

	"encore.app/authentication/business"
	"encore.app/authentication/store"
	"encore.app/pkg/myjwt"
	"encore.app/pkg/resendmailer"
	"encore.dev/config"
	"encore.dev/storage/sqldb"
)

var sharedDB = sqldb.NewDatabase("db_authentication", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})

var secrets struct {
	JWT_SECRET_KEY string
	RESEND_API_KEY string
}

type ServiceConfig struct {
	BaseUrl string
}

var cfg *ServiceConfig = config.Load[*ServiceConfig]()

//encore:service
type ServiceAuth struct {
	b *business.BusinessAuth
}

func initServiceAuth() (*ServiceAuth, error) {

	fmt.Println("secrets.JWT_SECRET_KEY:", secrets.JWT_SECRET_KEY)
	fmt.Println("cfg", cfg.BaseUrl)

	// Initialize the resend mailer
	m := resendmailer.NewResendMailer(secrets.RESEND_API_KEY, "Acme <onboarding@resend.dev>")

	s := store.NewAuthStore(sharedDB)
	tokenizer := myjwt.NewJWTTokenizer(secrets.JWT_SECRET_KEY)
	b := business.NewAuthBusiness(s, tokenizer, m)
	return &ServiceAuth{b: b}, nil
}

type MyAuthParams struct {
	AuthToken     *http.Cookie `cookie:"auth_token"`
	Authorization string       `header:"Authorization"`
}
