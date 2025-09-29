// Service auth implements user login and register functionality.
package auth

import (
	"context"
	"fmt"

	"encore.app/auth/business"
	"encore.app/auth/store"
	"encore.app/pkg/myjwt"
	"encore.app/pkg/resendmailer"
	"encore.dev/config"
	"encore.dev/storage/sqldb"
)

var sharedDB = sqldb.NewDatabase("shareddb", sqldb.DatabaseConfig{
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

//encore:api public path=/saludo/:name
func World(ctx context.Context, name string) (*Response, error) {
	msg := "Hello you are logging in, " + name + "!"
	return &Response{Message: msg}, nil
}

type Response struct {
	Message string
}
