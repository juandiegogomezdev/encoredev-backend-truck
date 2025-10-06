package utils

import (
	"fmt"

	"encore.dev/beta/errs"
	"encore.dev/types/uuid"
)

func MustNewUUID() (uuid.UUID, error) {
	newID, err := uuid.NewV4()
	if err != nil {
		fmt.Println("Error generating UUID:", err)
		return uuid.UUID{}, &errs.Error{
			Code:    errs.Internal,
			Message: "Error en el sistema",
		}
	}

	return newID, nil
}
