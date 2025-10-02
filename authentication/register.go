package authentication

import (
	"context"

	"encore.app/pkg/utils"
	"encore.app/sessions"
)

//encore:api public method=POST path=/auth/register
func (s *ServiceAuth) Register(ctx context.Context, req *RequestRegisterUser) (*ResponseRegister, error) {
	err := s.b.CheckUserExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	token, err := s.b.GenerateConfirmRegisterToken(req.Email)
	if err != nil {
		return nil, err
	}

	// Send email with token to confirm registration
	s.b.SendConfirmRegisterEmail(ctx, req.Email, token)

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
	resCreateOrgSelectSession, err := sessions.CreateOrgSelectSession(ctx, sessions.RequestCreateOrgSelectSession{
		UserID:     userID,
		DeviceInfo: "device info", // TODO: Get device info from request header
	})
	if err != nil {
		return nil, err
	}

	// Save session in the database

	// If is mobile, return token in response body
	if req.ClientType == "mobile" {
		// Return token in response body
		return &ResponseConfirmRegister{
			Message: "Usuario registrado correctamente",
			Token:   resCreateOrgSelectSession.OrgSelectSessionToken,
		}, nil

	}

	// If is web, return token in HttpOnly cookie
	cookieOptions := utils.DefaultCookieOptions("auth_token", resCreateOrgSelectSession.OrgSelectSessionToken)
	return &ResponseConfirmRegister{
		Message:   "Usuario registrado correctamente",
		Token:     "",
		SetCookie: cookieOptions,
	}, nil
}

type RequestConfirmRegisterUser struct {
	// Body params
	Password string `json:"password"`
	// Query param
	Token string `json:"token"`
	// header param for client type (web or mobile)
	ClientType string `header:"X-Client-Type"`
}

type ResponseConfirmRegister struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`

	SetCookie string `header:"Set-Cookie,omitempty"`
}
