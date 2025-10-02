package myjwt

import (
	"fmt"
	"time"

	"encore.dev/types/uuid"
	"github.com/golang-jwt/jwt/v5"
)

// Sign the token with the secret key
func (t *jwtTokenizer) signToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(t.secretKey))
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return tokenString, nil
}

// This token is used to confirm a new email during registration.
func (t *jwtTokenizer) GenerateConfirmRegisterToken(newEmail string) (string, error) {
	claims := ConfirmRegisterClaims{
		NewEmail: newEmail,
		BaseClaims: BaseClaims{
			TokenType: TokenTypeConfirmRegister,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
			},
		},
	}

	return t.signToken(claims)
}

// This token is used to confirm a login via email code.
func (t *jwtTokenizer) GenerateConfirmLoginToken(email string) (string, error) {
	claims := ConfirmLoginClaims{
		Email: email,
		BaseClaims: BaseClaims{
			TokenType: TokenTypeConfirmLogin,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
			},
		},
	}

	return t.signToken(claims)
}

// This token is used to select an organization after login.
func (t *jwtTokenizer) GenerateOrgSelectToken(userID uuid.UUID, sessionID uuid.UUID) (string, error) {
	claims := OrgSelectClaims{
		UserID:    userID,
		SessionID: sessionID,
		BaseClaims: BaseClaims{
			TokenType: TokenTypeOrgSelect,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
			},
		},
	}

	return t.signToken(claims)
}

// This token is used to access the api as a member of an organization.
func (t *jwtTokenizer) GenerateMembershipToken(membershipID uuid.UUID, sessionID uuid.UUID) (string, error) {
	claims := MembershipClaims{
		MembershipID: membershipID,
		SessionID:    sessionID,
		BaseClaims: BaseClaims{
			TokenType: TokenTypeMembership,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
			},
		},
	}

	return t.signToken(claims)
}
