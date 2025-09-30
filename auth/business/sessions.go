package business

import (
	"fmt"

	"encore.dev/beta/errs"
	"encore.dev/types/uuid"
)

func (b *BusinessAuth) GenerateOrgAccessToken(userID uuid.UUID) (string, error) {
	claims, err := b.tokenizer.GenerateOrgSelectToken(userID)
	if err != nil {
		fmt.Println("Error generating org access token:", err)
		return "", &errs.Error{
			Code:    errs.Internal,
			Message: "Error al generar el acceso",
		}
	}
	return claims, nil
}

func (b *BusinessAuth) GenerateConfirmRegisterToken(newEmail string) (string, error) {
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
