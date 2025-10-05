package appstore

import (
	"context"
	"time"

	"encore.dev/types/uuid"
)

// Search if user exists by email
func (s *StoreApp) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	q := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	err := s.db.QueryRow(ctx, q, email).Scan(&exists)
	if err != nil {
		return true, err
	}
	return exists, nil
}

// Create a new user in the database
func (s *StoreApp) CreateUser(ctx context.Context, user *CreateUserStoreStruct, verification *CreateUserVerificationStruct) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	qUser := `
		INSERT INTO users (id, email, hashed_password)
		VALUES ($1, $2, $3)`
	_, err = tx.Exec(ctx, qUser, user.ID, user.Email, user.HashedPassword)

	if err != nil {
		tx.Rollback()
		return err
	}

	qVerification := `
		INSERT INTO user_login_codes (user_id, code, expires_at)
		VALUES ($1, $2, $3)`

	_, err = tx.Exec(ctx, qVerification, verification.UserID, verification.Code, verification.ExpiresAt)

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
}

type CreateUserVerificationStruct struct {
	UserID    uuid.UUID
	Code      string
	ExpiresAt time.Time
}
