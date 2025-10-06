package appService

import (
	"context"

	"encore.app/pkg/utils"
)

//encore:api public method=POST path=/login
func (s *ServiceApp) Login(ctx context.Context, req *RequestLogin) (*ResponseLogin, error) {
	// Validate if the user exists and the password is correct
	userID, err := s.b.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	// Generate a token for the user to confirm the login via email code
	token, err := s.b.GenerateConfirmLoginToken(ctx, userID)
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

//encore:api public method=POST path=/login/confirm-code
func (s *ServiceApp) LoginConfirmCode(ctx context.Context, req *RequestLoginConfirmCode) (*ResponseLoginConfirmCode, error) {
	// Confirm the login with the code and generate a session token
	orgSelectTokenSession, err := s.b.LoginConfirm(ctx, req.Token, req.Code)
	if err != nil {
		return nil, err
	}

	cookie := utils.DefaultCookieOptions("auth_token", orgSelectTokenSession)
	switch req.ClientType {
	case "mobile":
		// Return token in response body
		return &ResponseLoginConfirmCode{Token: orgSelectTokenSession}, nil
	default:
		// Return token in HttpOnly cookie
		return &ResponseLoginConfirmCode{Token: "", SetCookie: cookie}, nil
	}
}

type RequestLoginConfirmCode struct {
	ClientType string `header:"X-Client-Type"` // Client type (mobile or web)
	Token      string `json:"token"`           // token with the userID
	Code       string `json:"code"`            // code to confirm the login
}

type ResponseLoginConfirmCode struct {
	SetCookie string `header:"Set-Cookie,omitempty"` // Set-Cookie header to set the session cookie
	Token     string `json:"token,omitempty"`        // token for mobile clients
}
