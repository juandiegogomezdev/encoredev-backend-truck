package appService

import (
	"net/http"

	"encore.app/appservice/appbusiness"
	"encore.app/appservice/appstore"
	"encore.app/pkg/myjwt"
	"encore.app/pkg/resendmailer"
	"encore.dev/config"
	"encore.dev/storage/sqldb"
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
	AuthToken     *http.Cookie `cookie:"auth_token"`
	Authorization string       `header:"Authorization"`
}
