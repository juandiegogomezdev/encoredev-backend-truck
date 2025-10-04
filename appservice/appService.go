package appService

import (
	"fmt"
	"net/http"

	"encore.app/appservice/appbusiness"
	"encore.app/appservice/appstore"
	"encore.app/pkg/myjwt"
	"encore.app/pkg/resendmailer"
	"encore.dev/config"
	"encore.dev/storage/sqldb"
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

	fmt.Println("secrets.JWT_SECRET_KEY:", secrets.JWT_SECRET_KEY)
	fmt.Println("cfg", cfg.BaseUrl)

	// Initialize the resend mailer
	m := resendmailer.NewResendMailer(secrets.RESEND_API_KEY, "Acme <onboarding@resend.dev>")

	s := appstore.NewAppStore(appDB)
	tokenizer := myjwt.NewJWTTokenizer(secrets.JWT_SECRET_KEY)
	b := appbusiness.NewAppBusiness(s, tokenizer, m)
	return &ServiceApp{b: b}, nil
}

type MyAuthParams struct {
	AuthToken     *http.Cookie `cookie:"auth_token"`
	Authorization string       `header:"Authorization"`
}
