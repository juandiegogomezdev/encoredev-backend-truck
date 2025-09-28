package myjwt

import (
	"fmt"

	"encore.dev/types/uuid"
	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	TokenTypeLogin      TokenType = "login"
	TokenTypeConfirmed  TokenType = "confirmed"
	TokenTypeMembership TokenType = "membership"
)

type LoginClaims struct {
	Email     string    `json:"email"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

type ConfirmedClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

type MembershipClaim struct {
	MembershipID uuid.UUID `json:"membership_id"`
	TokenType    TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

func ParseLoginToken(tokenString string) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &LoginClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-very-secret-key"), nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token.Claims, nil
}
