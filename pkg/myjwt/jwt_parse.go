package myjwt

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type TokenStatus string

const (
	TokenStatusValid   TokenStatus = "token_valid"
	TokenStatusExpired TokenStatus = "token_expired"
	TokenStatusInvalid TokenStatus = "token_invalid"
)

func (t *jwtTokenizer) parseToken(tokenString string, claims jwt.Claims) (jwt.Claims, TokenStatus) {

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.secretKey), nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return token.Claims, TokenStatusExpired
		case errors.Is(err, jwt.ErrTokenMalformed),
			errors.Is(err, jwt.ErrTokenSignatureInvalid),
			errors.Is(err, jwt.ErrTokenNotValidYet),
			errors.Is(err, jwt.ErrTokenInvalidClaims):

			fmt.Println("Error parsing token:", err)
			return nil, TokenStatusInvalid
		default:
			return nil, TokenStatusInvalid
		}
	}

	if !token.Valid {
		return nil, TokenStatusInvalid
	}

	return token.Claims, TokenStatusValid

}

// Parse the token for confirm login
func (t *jwtTokenizer) ParseConfirmLoginToken(tokenString string) (*ConfirmLoginClaims, TokenStatus) {
	parsedClaims, tokenStatus := t.parseToken(tokenString, &ConfirmLoginClaims{})

	switch tokenStatus {
	case TokenStatusInvalid:
		return nil, TokenStatusInvalid
	default:
		return parsedClaims.(*ConfirmLoginClaims), tokenStatus
	}
}

// Parse the token for confirm register
func (t *jwtTokenizer) ParseConfirmRegisterToken(tokenString string) (*ConfirmRegisterClaims, TokenStatus) {
	parsedClaims, tokenStatus := t.parseToken(tokenString, &ConfirmRegisterClaims{})
	switch tokenStatus {
	case TokenStatusInvalid:
		return nil, TokenStatusInvalid
	default:
		return parsedClaims.(*ConfirmRegisterClaims), tokenStatus
	}
}

// Parse the token for organization selection
func (t *jwtTokenizer) ParseOrgSelectToken(tokenString string) (*OrgSelectClaims, TokenStatus) {
	parsedClaims, tokenStatus := t.parseToken(tokenString, &OrgSelectClaims{})
	switch tokenStatus {
	case TokenStatusInvalid:
		return nil, TokenStatusInvalid
	default:
		return parsedClaims.(*OrgSelectClaims), tokenStatus
	}
}

// Parse the token for membership
func (t *jwtTokenizer) ParseMembershipToken(tokenString string) (*MembershipClaims, TokenStatus) {
	parsedClaims, tokenStatus := t.parseToken(tokenString, &MembershipClaims{})
	switch tokenStatus {
	case TokenStatusInvalid:
		return nil, TokenStatusInvalid
	default:
		return parsedClaims.(*MembershipClaims), tokenStatus
	}
}

func (t *jwtTokenizer) ParseBaseClaims(tokenString string) (*BaseClaims, TokenStatus) {
	parsedClaims, tokenStatus := t.parseToken(tokenString, &BaseClaims{})
	switch tokenStatus {
	case TokenStatusInvalid:
		return nil, TokenStatusInvalid
	default:
		return parsedClaims.(*BaseClaims), tokenStatus
	}
}

// Parse anyone token and return the claims
func (t *jwtTokenizer) ParseFullClaims(tokenString string) (*FullClaims, TokenStatus) {
	parsedClaims, tokenStatus := t.parseToken(tokenString, &FullClaims{})
	switch tokenStatus {
	case TokenStatusInvalid:
		return nil, TokenStatusInvalid
	default:
		return parsedClaims.(*FullClaims), tokenStatus
	}
}
