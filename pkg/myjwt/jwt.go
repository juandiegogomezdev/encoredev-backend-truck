package myjwt

import (
	"encore.dev/types/uuid"
	"github.com/golang-jwt/jwt/v5"
)

type JWTTokenizer interface {
	// Generation
	signToken(claims jwt.Claims) (string, error)

	GenerateConfirmRegisterToken(newEmail string) (string, error)
	GenerateConfirmLoginToken(email string) (string, error)
	GenerateOrgSelectToken(userID uuid.UUID, sessionID uuid.UUID) (string, error)
	GenerateMembershipToken(membershipID uuid.UUID, sessionID uuid.UUID) (string, error)

	// Parsing
	parseToken(tokenString string, claims jwt.Claims) (jwt.Claims, error)
	// ParseFullClaims(tokenString string) (*FullClaims, error)

	ParseBaseClaims(tokenString string) (*BaseClaims, error)
	ParseConfirmRegisterToken(tokenString string) (*ConfirmRegisterClaims, error)
	ParseConfirmLoginToken(tokenString string) (*ConfirmLoginClaims, error)
	ParseOrgSelectToken(tokenString string) (*OrgSelectClaims, error)
	ParseMembershipToken(tokenString string) (*MembershipClaims, error)
}

type jwtTokenizer struct {
	secretKey string
}

func NewJWTTokenizer(secretKey string) JWTTokenizer {
	return &jwtTokenizer{secretKey: secretKey}
}

type TokenType string

const (
	TokenTypeConfirmRegister TokenType = "confirm_register"
	TokenTypeConfirmLogin    TokenType = "confirm_login"
	TokenTypeOrgSelect       TokenType = "org_select"
	TokenTypeMembership      TokenType = "membership"
	TokenTypeUnknown         TokenType = "unknown"
)

// BaseClaims holds common JWT claims.
type BaseClaims struct {
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

// Used for confirm register
type ConfirmRegisterClaims struct {
	NewEmail string `json:"new_email"`
	BaseClaims
}

// Used for confirm login
type ConfirmLoginClaims struct {
	Email string `json:"email"`
	BaseClaims
}

// Used for allow selecting an organization
type OrgSelectClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	SessionID uuid.UUID `json:"session_id"`
	BaseClaims
}

// Used for use the app as a member of an organization
type MembershipClaims struct {
	MembershipID uuid.UUID `json:"membership_id"`
	SessionID    uuid.UUID `json:"session_id"`
	BaseClaims
}

// Used for access static resources.
// This token contains all the properties of the previous tokens
type FullClaims struct {
	NewEmail     string    `json:"new_email"`
	Email        string    `json:"email"`
	UserID       uuid.UUID `json:"user_id"`
	MembershipID uuid.UUID `json:"membership_id"`
	SessionID    uuid.UUID `json:"session_id"`
	BaseClaims
}
