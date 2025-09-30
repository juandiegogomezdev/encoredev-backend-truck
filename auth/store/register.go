package store

import (
	"context"
	"time"

	"encore.dev/types/uuid"
)

// Search if user exists by email
func (s *AuthStore) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	q := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	err := s.db.QueryRow(ctx, q, email).Scan(&exists)
	if err != nil {
		return true, err
	}
	return exists, nil
}

// Create a new user in the database
func (s *AuthStore) CreateUser(ctx context.Context, user *CreateUserStoreStruct, verification *CreateUserVerificationStruct) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	qUser := `
		INSERT INTO users (id, email, hashed_password, created_at)
		VALUES ($1, $2, $3, $4)`
	_, err = tx.Exec(ctx, qUser, user.ID, user.Email, user.HashedPassword, user.CreatedAt)

	if err != nil {
		tx.Rollback()
		return err
	}

	qVerification := `
		INSERT INTO user_login_codes (user_id, code, created_at, expires_at)
		VALUES ($1, $2, $3, $4)`

	_, err = tx.Exec(ctx, qVerification, verification.UserID, verification.Code, verification.CreatedAt, verification.ExpiresAt)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil

}

type CreateUserStoreStruct struct {
	ID             uuid.UUID
	Email          string
	HashedPassword string
	CreatedAt      time.Time
}

type CreateUserVerificationStruct struct {
	UserID    uuid.UUID
	Code      string
	CreatedAt time.Time
	ExpiresAt time.Time
}
