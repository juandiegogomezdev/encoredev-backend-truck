package business

import (
	"context"
	"fmt"
	"time"

	"encore.app/pkg/myjwt"
	"encore.app/sessions/store"
	"encore.dev/beta/errs"
	"encore.dev/types/uuid"
)

type BusinessSession struct {
	store *store.SessionStore

	tokenizer myjwt.JWTTokenizer
}

func NewBusinessSession(store *store.SessionStore, tokenizer myjwt.JWTTokenizer) *BusinessSession {
	return &BusinessSession{store: store, tokenizer: tokenizer}
}

// Check if is posible create a new session
func (b *BusinessSession) isPosibleCreateNewSession(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	// Count active sessions for the user
	count, err := b.store.CountSessionsByUserID(ctx, userID)
	if err != nil {
		fmt.Println("Error counting user sessions:", err)
		return uuid.Nil, &errs.Error{
			Code:    errs.Internal,
			Message: "Error al generar acceso al usuario",
		}
	}

	// Limit to 5 active sessions
	if count >= 5 {
		return uuid.Nil, &errs.Error{
			Code:    errs.PermissionDenied,
			Message: "Has alcanzado el límite de sesiones activas. Cierra sesión en otros dispositivos para continuar.",
		}
	}

	// Generate a new session ID
	sesionID, err := uuid.NewV4()
	if err != nil {
		fmt.Println("Error generating session ID:", err)
		return uuid.Nil, &errs.Error{
			Code:    errs.Internal,
			Message: "Error al generar acceso al usuario",
		}
	}

	return sesionID, nil
}

// Create a session for enter to org-select page
func (b *BusinessSession) CreateOrgSelectSession(ctx context.Context, userID uuid.UUID, deviceInfo string) (string, error) {
	// Check if is posible create a new session and create the sessionID
	sessionID, err := b.isPosibleCreateNewSession(ctx, userID)

	// Create the org select token
	tokenOrgSelect, err := b.tokenizer.GenerateOrgSelectToken(userID, sessionID)

	// New session parameters
	newSession := store.CreateUserSessionStruct{
		UserID:     userID,
		SessionID:  sessionID,
		DeviceInfo: deviceInfo,
		ExpiresAt:  time.Now().Add(25 * time.Hour),
	}

	// Save the new session in the database
	err = b.store.CreateUserSession(ctx, newSession)
	if err != nil {
		return "", err
	}

	return tokenOrgSelect, nil

}

// Create a session for enter to the app and use the apis
func (b *BusinessSession) CreateMembershipSession(ctx context.Context, membershipID uuid.UUID, sessionID uuid.UUID) (string, error) {

	// Create the org select token
	tokenMembership, err := b.tokenizer.GenerateMembershipToken(membershipID, sessionID)
	if err != nil {
		fmt.Println("Error generating membership token:", err)
		return "", &errs.Error{
			Code:    errs.Internal,
			Message: "Error al generar el acceso",
		}
	}

	return tokenMembership, nil
}

func (b *BusinessSession) DeleteUserSession(ctx context.Context, sessionID uuid.UUID) error {
	err := b.store.DeleteUserSession(ctx, sessionID)
	if err != nil {
		fmt.Println("Error deleting user session:", err)
		return &errs.Error{
			Code:    errs.Internal,
			Message: "Error al eliminar la sesión",
		}
	}
	return nil
}

// Check if a session is expired (Refresh token)
func (b *BusinessSession) CheckSessionIsActive(ctx context.Context, sessionID uuid.UUID) (bool, error) {
	isActive, err := b.store.IsActiveSession(ctx, sessionID)

	if err != nil {
		return false, err
	}

	return isActive, nil

}
