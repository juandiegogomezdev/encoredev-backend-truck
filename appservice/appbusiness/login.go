package appbusiness

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"encore.app/pkg/myjwt"
	"encore.dev/beta/errs"
	"encore.dev/storage/sqldb"
	"encore.dev/types/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Validate if the provided password matches the stored hash.
func (b *BusinessApp) validatePassword(storedHash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	return err == nil
}

// Generate a code with n digits
func (b *BusinessApp) generateCodeLogin(n int) string {
	numbers := "0123456789"
	code := make([]byte, n)
	for i := range code {
		code[i] = numbers[rand.Intn(len(numbers))]
	}
	return string(code)
}

// Login authenticates a user and returns a JWT token with his userID if successful.
func (b *BusinessApp) Login(ctx context.Context, email string, password string) (uuid.UUID, error) {

	// Fetch user by email
	user, err := b.store.GetUserByEmail(ctx, email)
	if err != nil {

		if errors.Is(err, sqldb.ErrNoRows) {
			return uuid.Nil, &errs.Error{
				Code:    errs.NotFound,
				Message: "Usuario no encontrado",
			}
		}
		return uuid.Nil, &errs.Error{
			Code:    errs.Internal,
			Message: "Error al obtener el usuario",
		}
	}

	// Verify password
	if !b.validatePassword(user.PasswordHash, password) {
		return uuid.Nil, &errs.Error{
			Code:    errs.Unauthenticated,
			Message: "Contraseña incorrecta",
		}
	}

	return user.ID, nil
}

// Generate a token to confirm the login via email code
func (b *BusinessApp) GenerateConfirmLoginToken(ctx context.Context, userID uuid.UUID) (string, error) {
	tokenConfirmLogin, err := b.tokenizer.GenerateConfirmLoginToken(userID)
	if err != nil {
		fmt.Println("Error generating confirm login token:", err)
		return "", &errs.Error{
			Code:    errs.Internal,
			Message: "Error al generar el token de confirmación",
		}
	}

	return tokenConfirmLogin, nil
}

// Confirm the login with a code and generate a session
func (b *BusinessApp) LoginConfirm(ctx context.Context, token, code string) (string, error) {
	// Validate the token
	claims, tokenStatus := b.tokenizer.ParseConfirmLoginToken(token)

	// Check the token status
	switch tokenStatus {
	case myjwt.TokenStatusValid:
		// Valid token,
	case myjwt.TokenStatusInvalid:
		return "", &errs.Error{
			Code:    errs.Unauthenticated,
			Message: "Retrocede e ingresa tu correo nuevamente.",
		}
	case myjwt.TokenStatusExpired:
		return "", &errs.Error{
			Code:    errs.Unauthenticated,
			Message: "Retrocede e ingresa tu correo nuevamente.",
		}
	default:
		fmt.Println("Case no handled:", tokenStatus)
		return "", &errs.Error{
			Code:    errs.Unauthenticated,
			Message: "Retrocede e ingresa tu correo nuevamente.",
		}
	}

	// Get the code from the database
	userLoginCode, err := b.store.GetUserLoginCodeByUserID(ctx, claims.UserID)
	if err != nil {
		return "", &errs.Error{
			Code:    errs.Internal,
			Message: "Error en el sistema",
		}
	}

	// Check if the code is expired
	if !userLoginCode.ExpiresAt.Before(time.Now()) {
		return "", &errs.Error{
			Code:    errs.Unauthenticated,
			Message: "El código ha expirado. Por favor, inicia sesión nuevamente.",
		}
	}

	// Check if the code matches
	if userLoginCode.Code != code {
		return "", &errs.Error{
			Code:    errs.Unauthenticated,
			Message: "Código incorrecto",
		}
	}

	// Create a session for the user
	orgSelectToken, err := b.CreateOrgSelectSession(ctx, claims.UserID, "web")
	if err != nil {
		return "", err
	}

	return orgSelectToken, nil

}
