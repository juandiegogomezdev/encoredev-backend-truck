package auth

import "context"

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
