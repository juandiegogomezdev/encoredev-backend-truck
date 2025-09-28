// Service auth implements user login and register functionality.
package auth

import (
	"context"

	"encore.app/auth/business"
	"encore.app/auth/store"
	"encore.app/shared/myjwt"
	"encore.dev/storage/sqldb"
)

var sharedDB = sqldb.NewDatabase("shareddb", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})

// var jwtSecretKey = config.Secret("JWT_SECRET_KEY")

//encore:service
type ServiceAuth struct {
	b *business.BusinessAuth
}

func initServiceAuth() (*ServiceAuth, error) {
	s := store.NewAuthStore(sharedDB)
	tokenizer := myjwt.NewJWTTokenizer("MY_SECRET")
	b := business.NewAuthBusiness(s, tokenizer)
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

//encore:api public method=POST path=/auth/login
func (s *ServiceAuth) Login(ctx context.Context, req *RequestLogin) (*ResponseLogin, error) {
	token, err := s.b.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &ResponseLogin{Token: token}, nil
}

type RequestLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResponseLogin struct {
	Token string `json:"token"`
}
