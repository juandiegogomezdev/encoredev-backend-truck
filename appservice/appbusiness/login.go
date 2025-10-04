package appbusiness

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"encore.dev/beta/errs"
	"encore.dev/storage/sqldb"
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

// Login authenticates a user and returns a JWT token with his email if successful.
func (b *BusinessApp) Login(ctx context.Context, email string, password string) (string, error) {

	// Fetch user by email
	user, err := b.store.GetUserByEmail(ctx, email)
	if err != nil {

		if errors.Is(err, sqldb.ErrNoRows) {
			return "", &errs.Error{
				Code:    errs.NotFound,
				Message: "Usuario no encontrado",
			}
		}
		return "", &errs.Error{
			Code:    errs.Internal,
			Message: "Error al obtener el usuario",
		}
	}

	// Verify password
	if !b.validatePassword(user.PasswordHash, password) {
		return "", &errs.Error{
			Code:    errs.Unauthenticated,
			Message: "Contraseña incorrecta",
		}
	}

	fmt.Println("Usuario autenticado:", user.ID)

	return "token-here", nil

}
