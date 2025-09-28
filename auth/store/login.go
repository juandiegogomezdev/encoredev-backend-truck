package store

import (
	"context"
	"fmt"
	"time"

	"encore.dev/types/uuid"
)

// GetUserByEmailResult holds the result of GetUserByEmail.
type GetUserByEmailResult struct {
	ID           uuid.UUID
	PasswordHash string
}

// GetUserByEmail retrieves a user's ID by their email.
func (s *AuthStore) GetUserByEmail(ctx context.Context, email string) (GetUserByEmailResult, error) {
	var result GetUserByEmailResult

	// Debug checks
	if s == nil {
		return GetUserByEmailResult{}, fmt.Errorf("AuthStore is nil")
	}
	if s.db == nil {
		return GetUserByEmailResult{}, fmt.Errorf("database is nil")
	}

	fmt.Printf("Executing query with email: %s\n", email)
	q := "SELECT id, hashed_password FROM users WHERE email = $1"
	err := s.db.QueryRow(ctx, q, email).Scan(&result.ID, &result.PasswordHash)
	if err != nil {
		return GetUserByEmailResult{}, err
	}
	return result, nil
}

// GetUserLoginCodeResult holds the result of GetUserLoginCodeByUserID.
type GetUserLoginCodeResult struct {
	Code      string
	ExpiresAt string
}

// GetUserLoginCodeByUserID retrieves a user's login code by their user ID.
func (s *AuthStore) GetUserLoginCodeByUserID(ctx context.Context, userID uuid.UUID) (GetUserLoginCodeResult, error) {
	var result GetUserLoginCodeResult
	q := "SELECT code, expires_at FROM user_login_codes WHERE user_id = $1"
	err := s.db.QueryRow(ctx, q, userID).Scan(&result.Code, &result.ExpiresAt)
	if err != nil {
		return GetUserLoginCodeResult{}, err
	}
	return result, nil
}

// CreateUserLoginCode creates a new login code for a user.
func (s *AuthStore) CreateUserLoginCode(ctx context.Context, userID uuid.UUID, code string, created_at time.Time, expiresAt time.Time) error {
	q := "INSERT INTO user_login_codes (user_id, code, expires_at) VALUES ($1, $2, $3)"
	_, err := s.db.Exec(ctx, q, userID, code, expiresAt)
	return err
}
