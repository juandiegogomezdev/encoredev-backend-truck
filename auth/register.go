package auth

import (
	"context"

	"encore.dev"
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
func (s *ServiceAuth) ConfirmUserRegister(ctx context.Context, req *RequestConfirmRegisterUser) (*ResponseConfirmRegister, error) {
	// Get token with email
	tokenConfirmEmail := req.Token

	// Extract new email from the token
	newEmail, err := s.b.ExtractNewEmailFromToken(ctx, tokenConfirmEmail)
	if err != nil {
		return nil, err
	}

	// Check if user already exists. If exists, return error
	err = s.b.CheckUserExists(ctx, newEmail)
	if err != nil {
		return nil, err
	}

	// Create user in the database
	userID, err := s.b.CreateUser(ctx, newEmail, req.Password)
	if err != nil {
		return nil, err
	}

	// Generate login token
	token, err := s.b.GenerateOrgAccess(userID)
	if err != nil {
		return nil, err
	}

	if req.ClientType == "web" {
		currentReq := encore.CurrentRequest()

	}
	return &ResponseLogin{Token: token}, nil
}

type RequestConfirmRegisterUser struct {
	// Body params
	Password string `json:"password"`

	// Query param
	Token string `query:"token"`

	// header param for client type (web or mobile)
	ClientType string `header:"X-Client-Type"`
}

type ResponseConfirmRegister struct {
	Message   string  `json:"message"`
	Token     *string `json:"token,omitempty"`
	SetCookie *string `header:"Set-Cookie,omitempty"`
}
