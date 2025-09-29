package auth

import (
	"context"
)

//encore:api public method=POST path=/auth/register
func (s *ServiceAuth) Register(ctx context.Context, req *RequestRegisterUser) (*ResponseRegister, error) {
	token, err := s.b.CheckUser(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	return &ResponseRegister{
		Message: "Revisa tu correo para confirmar el registro",
		Token:   token,
	}, nil
}

type ResponseRegister struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}
type RequestRegisterUser struct {
	Email string `json:"email"`
}

//encore:api public method=POST path=/auth/confirm-register
func (s *ServiceAuth) ConfirmUserRegister(ctx context.Context, req *RequestConfirmRegisterUser) (*ResponseLogin, error) {
	token, err := s.b.CheckUser(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	return &ResponseLogin{Token: token}, nil
}

type RequestConfirmRegisterUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
