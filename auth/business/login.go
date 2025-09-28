package business

import (
	"context"
	"errors"
	"fmt"

	"encore.dev/beta/errs"
	"encore.dev/storage/sqldb"
	"golang.org/x/crypto/bcrypt"
)

// Validate if the provided password matches the stored hash.
func (b *BusinessAuth) validatePassword(storedHash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	return err == nil
}

// Login authenticates a user and returns a JWT token with his email if successful.
func (b *BusinessAuth) Login(ctx context.Context, email string, password string) (string, error) {

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
