package business

import (
	"context"
	"fmt"
	"time"

	"encore.app/authentication/store"
	"encore.dev/beta/errs"
	"encore.dev/types/uuid"
)

// This function extracts the new email from the provided token.
func (b *BusinessAuth) ExtractNewEmailFromToken(ctx context.Context, token string) (string, error) {
	claims, err := b.tokenizer.ParseConfirmRegisterToken(token)
	if err != nil {
		return "", &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "Link de confirmacion inválido o vencido",
		}
	}
	return claims.NewEmail, nil
}

// Check if the user exists
func (b *BusinessAuth) CheckUserExists(ctx context.Context, email string) error {
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
func (b *BusinessAuth) SendConfirmRegisterEmail(ctx context.Context, email string, token string) {
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
func (b *BusinessAuth) CreateUser(ctx context.Context, newEmail string, password string) (uuid.UUID, error) {
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

	newUser := store.CreateUserStoreStruct{
		ID:             userID,
		Email:          newEmail,
		HashedPassword: hashedPassword,
		CreatedAt:      time.Now(),
	}

	newUserVerification := store.CreateUserVerificationStruct{
		UserID:    userID,
		Code:      b.generateCodeLogin(6),
		CreatedAt: time.Now(),
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
