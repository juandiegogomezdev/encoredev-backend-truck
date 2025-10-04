package appstore

import (
	"context"
	"database/sql"
	"time"

	"encore.dev/types/uuid"
)

func (s *AppStore) IsActiveSession(ctx context.Context, sessionID uuid.UUID) (bool, error) {
	var isActive bool
	q := `
		SELECT is_active
		FROM user_sessions
		WHERE session_id = $1
	`

	err := s.db.QueryRow(ctx, q, sessionID).Scan(&isActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return isActive, nil
}

// Count sessions of a user
func (s *AppStore) CountSessionsByUserID(ctx context.Context, userID uuid.UUID) (int8, error) {
	var count int8
	q := `SELECT COUNT(*) FROM user_sessions WHERE user_id = $1`
	err := s.db.QueryRow(ctx, q, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Create a sesssion for a user
func (s *AppStore) CreateUserSession(ctx context.Context, newSession CreateUserSessionStruct) error {

	q := `
		INSERT INTO user_sessions (user_id, session_id, device_info, expires_at)
		VALUES ($1, $2, $3, $4)`

	_, err := s.db.Exec(ctx, q, newSession.UserID, newSession.SessionID, newSession.DeviceInfo, newSession.ExpiresAt)
	return err
}

// Delete a session of the user
func (s *AppStore) DeleteUserSession(ctx context.Context, sessionID uuid.UUID) error {
	q := `DELETE FROM user_sessions WHERE session_id = $1		`

	_, err := s.db.Exec(ctx, q, sessionID)
	return err
}

type CreateUserSessionStruct struct {
	UserID     uuid.UUID
	SessionID  uuid.UUID
	DeviceInfo string
	ExpiresAt  time.Time
}
