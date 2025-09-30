package myjwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func (t *jwtTokenizer) parseToken(tokenString string, claims jwt.Claims) (jwt.Claims, error) {

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.secretKey), nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token.Claims, nil

}

// Parse the token for confirm login
func (t *jwtTokenizer) ParseConfirmLoginToken(tokenString string) (*ConfirmLoginClaims, error) {
	parsedClaims, err := t.parseToken(tokenString, &ConfirmLoginClaims{})
	if err != nil {
		return nil, err
	}
	return parsedClaims.(*ConfirmLoginClaims), nil
}

// Parse the token for confirm register
func (t *jwtTokenizer) ParseConfirmRegisterToken(tokenString string) (*ConfirmRegisterClaims, error) {
	parsedClaims, err := t.parseToken(tokenString, &ConfirmRegisterClaims{})
	if err != nil {
		return nil, err
	}
	return parsedClaims.(*ConfirmRegisterClaims), nil
}

// Parse the token for organization selection
func (t *jwtTokenizer) ParseOrgSelectToken(tokenString string) (*OrgSelectClaims, error) {
	parsedClaims, err := t.parseToken(tokenString, &OrgSelectClaims{})
	if err != nil {
		return nil, err
	}
	return parsedClaims.(*OrgSelectClaims), nil
}

// Parse the token for membership
func (t *jwtTokenizer) ParseMembershipToken(tokenString string) (*MembershipClaims, error) {
	parsedClaims, err := t.parseToken(tokenString, &MembershipClaims{})
	if err != nil {
		return nil, err
	}
	return parsedClaims.(*MembershipClaims), nil
}

func (t *jwtTokenizer) ParseBaseClaims(tokenString string) (*BaseClaims, error) {
	parsedClaims, err := t.parseToken(tokenString, &BaseClaims{})
	if err != nil {
		return nil, err
	}
	return parsedClaims.(*BaseClaims), nil
}

// // Parse anyone token and return the claims
// func (t *jwtTokenizer) ParseFullClaims(tokenString string) (*FullClaims, error) {
// 	parsedClaims, err := t.parseToken(tokenString, &FullClaims{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return parsedClaims.(*FullClaims), nil
// }
