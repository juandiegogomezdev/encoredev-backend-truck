package appbusiness

import (
	"context"
	"fmt"
	"time"

	"encore.app/appservice/appstore"
	"encore.app/pkg/myjwt"
	"encore.dev/beta/errs"
	"encore.dev/types/uuid"
)

func (b *BusinessApp) GenerateConfirmRegisterToken(newEmail string) (string, error) {
	claims, err := b.tokenizer.GenerateConfirmRegisterToken(newEmail)
	if err != nil {
		fmt.Println("Error generating confirm register token:", err)
		return "", &errs.Error{
			Code:    errs.Internal,
			Message: "Error al generar el token de confirmación",
		}
	}
	return claims, nil
}

// Check if is posible create a new session and return the new session ID
func (b *BusinessApp) isPosibleCreateNewSession(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
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
func (b *BusinessApp) CreateOrgSelectSession(ctx context.Context, userID uuid.UUID, deviceInfo string) (string, error) {
	// Check if is posible create a new session and create the sessionID
	sessionID, err := b.isPosibleCreateNewSession(ctx, userID)

	// Create the org select token
	tokenOrgSelect, err := b.tokenizer.GenerateOrgSelectToken(userID, sessionID)

	// New session parameters
	newSession := appstore.CreateUserSessionStruct{
		UserID:     userID,
		SessionID:  sessionID,
		DeviceInfo: deviceInfo,
		ExpiresAt:  time.Now().Add(25 * time.Hour),
	}

	// Save the new session in the database
	err = b.store.CreateUserSession(ctx, newSession)
	if err != nil {
		fmt.Println("Error creating user session:", err)
		return "", &errs.Error{
			Code:    errs.Internal,
			Message: "Error al generar el acceso",
		}
	}

	return tokenOrgSelect, nil

}

// Create a session for enter to the app and use the apis
func (b *BusinessApp) CreateMembershipSession(ctx context.Context, userID, membershipID, sessionID uuid.UUID) (string, error) {

	// Create the org select token
	tokenMembership, err := b.tokenizer.GenerateMembershipToken(userID, membershipID, sessionID)
	if err != nil {
		fmt.Println("Error generating membership token:", err)
		return "", &errs.Error{
			Code:    errs.Internal,
			Message: "Error al generar el acceso",
		}
	}

	return tokenMembership, nil
}

// Delete the user session of the database.
func (b *BusinessApp) DeleteUserSession(ctx context.Context, tokenSession string) {
	// Parse the token to get the session ID
	claims, tokenStatus := b.tokenizer.ParseFullClaims(tokenSession)

	switch tokenStatus {
	case myjwt.TokenStatusValid, myjwt.TokenStatusExpired:
		sessionID := claims.SessionID
		err := b.store.DeleteUserSession(ctx, sessionID)
		if err != nil {
			// TODO : log the error
			fmt.Println("Error deleting user session:", err)
		}

	default:
		fmt.Println("Invalid token:", tokenStatus)
	}
}

// Check if a session is expired (Refresh token)
func (b *BusinessApp) CheckSessionIsActive(ctx context.Context, sessionID uuid.UUID) (bool, error) {
	isActive, err := b.store.IsActiveSession(ctx, sessionID)

	if err != nil {
		return false, err
	}

	return isActive, nil

}
