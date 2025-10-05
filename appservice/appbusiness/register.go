package appbusiness

import (
	"context"
	"fmt"
	"time"

	"encore.app/appservice/appstore"
	"encore.app/pkg/myjwt"
	"encore.dev/beta/errs"
	"encore.dev/types/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Hash the password using bcrypt.
func GenerateHashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// This function extracts the new email from the provided token.
func (b *BusinessApp) ExtractNewEmailFromToken(ctx context.Context, token string) (string, error) {
	fmt.Println("Extracting new email from token:", token)
	claims, tokenStatus := b.tokenizer.ParseConfirmRegisterToken(token)
	fmt.Println("token", claims)
	switch tokenStatus {
	case myjwt.TokenStatusValid:
		// Token is valid, proceed
		return claims.NewEmail, nil
	case myjwt.TokenStatusExpired:
		return "", &errs.Error{
			Code:    errs.Unauthenticated,
			Message: "El link de confirmación ha expirado. Por favor, regístrate de nuevo.",
		}
	case myjwt.TokenStatusInvalid:
		return "", &errs.Error{
			Code:    errs.Unauthenticated,
			Message: "Link de confirmación inválido. Por favor, verifica el enlace o solicita uno nuevo.",
		}
	default:
		return "", &errs.Error{
			Code:    errs.Internal,
			Message: "Error al procesar el token de confirmación.",
		}
	}
}

// Check if the user exists
func (b *BusinessApp) CheckUserExists(ctx context.Context, email string) error {
	exists, err := b.store.UserExistsByEmail(ctx, email)
	if err != nil {
		return &errs.Error{
			Code:    errs.Internal,
			Message: "Error al comprobar si el usuario existe",
		}
	}
	if exists {
		return &errs.Error{
			Code:    errs.AlreadyExists,
			Message: "El usuario ya existe",
		}
	}
	return nil
}

// Send email with token to confirm registration
func (b *BusinessApp) SendConfirmRegisterEmail(ctx context.Context, email string, token string) {
	// Send email with the token (using the mailer)

	body := fmt.Sprint(`
		<h1> Bienvenido </h1>
		<p> Gracias por registrarte. Por favor, confirma tu correo electrónico haciendo clic en el siguiente enlace: </p>
		<a href="` + "http://localhost:4000/static/confirm-register?token=" + token + `"> Confirmar correo </a>
		<p> Si no te has registrado, ignora este correo. </p>
	`)

	// TODO: Improve error handling
	go func() {
		b.mailer.Send(email, "Confirm your registration", body)
	}()

}

// Create user in the database
func (b *BusinessApp) CreateUser(ctx context.Context, newEmail string, password string) (uuid.UUID, error) {
	userID, err := uuid.NewV4()

	if err != nil {
		fmt.Println("Error generating UUID:", err)
		return uuid.Nil, &errs.Error{
			Code:    errs.Internal,
			Message: "Error al registrar el usuario",
		}
	}

	hashedPassword, err := GenerateHashPassword(password)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return uuid.Nil, &errs.Error{
			Code:    errs.Internal,
			Message: "Error al registrar el usuario",
		}
	}

	newUser := appstore.CreateUserStoreStruct{
		ID:             userID,
		Email:          newEmail,
		HashedPassword: hashedPassword,
	}

	newUserVerification := appstore.CreateUserVerificationStruct{
		UserID:    userID,
		Code:      b.generateCodeLogin(6),
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	err = b.store.CreateUser(ctx, &newUser, &newUserVerification)
	if err != nil {
		fmt.Println("Error creating user:", err)
		return uuid.Nil, &errs.Error{
			Code:    errs.Internal,
			Message: "Error al registrar el usuario",
		}
	}

	return userID, nil
}
