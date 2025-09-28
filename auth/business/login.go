package business

import (
	"context"
	"fmt"

	"encore.dev/beta/errs"
	"encore.dev/storage/sqldb"
	"golang.org/x/crypto/bcrypt"
)

func (b *BusinessAuth) validatePassword(storedHash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	return err == nil
}

// Login authenticates a user and returns a JWT token if successful.
func (b *BusinessAuth) Login(ctx context.Context, email string, password string) (string, error) {
	fmt.Println("Iniciando proceso de login para:", email)
	fmt.Println("Contraseña proporcionada:", password)

	// Debug checks
	if b == nil {
		fmt.Println("Error: BusinessAuth is nil")
		return "", fmt.Errorf("BusinessAuth is nil")
	}
	if b.store == nil {
		fmt.Println("Error: store no inicializado")
		return "", fmt.Errorf("store is nil")
	}

	fmt.Printf("Store: %+v\n", b.store)
	// Fetch user by email
	user, err := b.store.GetUserByEmail(ctx, email)
	if err != nil {

		fmt.Printf("Error getting user by email: %v\n", err)
		if err == sqldb.ErrNoRows {
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
