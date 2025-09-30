package business

import (
	"fmt"

	"encore.dev/beta/errs"
	"encore.dev/types/uuid"
)

func (b *BusinessAuth) GenerateOrgAccess(userID uuid.UUID) (string, error) {
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
