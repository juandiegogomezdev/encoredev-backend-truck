package appbusiness

import (
	"context"
	"fmt"

	"encore.dev/beta/errs"
	"encore.dev/types/uuid"
)

// Get all memberships of a user
// func (b *BusinessApp) GetAllUserMemberships(ctx context.Context, userID uuid.UUID) ([]appstore.ResUserOrganizationStore, error) {

// 	organizations, err := b.store.GetAllUserMemberships(ctx, userID)
// 	if err != nil {
// 		return nil, &errs.Error{
// 			Code:    errs.Internal,
// 			Message: "Error al obtener las organizaciones del usuario",
// 		}
// 	}

// 	return organizations, nil
// }

// Create a personal organization
func (b *BusinessApp) CreatePersonalOrganization(ctx context.Context, userID uuid.UUID, name string) error {
	// Check not exists organizations for this user
	existingOrgs, err := b.store.GetAllUserOrganizations(ctx, userID)
	if err != nil {
		fmt.Println("Error checking existing organizations:", err)
		return &errs.Error{
			Code:    errs.Internal,
			Message: "Error al crear la organización personal",
		}
	}

	// Only is posible create a personal organization if has zero organizations
	if len(existingOrgs) == 0 {
		err := b.store.CreateUserOrganization(ctx, userID, name, "personal")
		if err != nil {
			return &errs.Error{
				Code:    errs.Internal,
				Message: "Error al crear el trabajo como independiente",
			}
		}
		return nil
	}
	return &errs.Error{
		Code:    errs.InvalidArgument,
		Message: "No es posible crear un trabajo como independiente porque ya existe",
	}
}

// Create a company organization
func (b *BusinessApp) CreateCompanyOrganization(ctx context.Context, userID uuid.UUID, name string) error {
	// Get all organizations for this user
	existingOrgs, err := b.store.GetAllUserOrganizations(ctx, userID)
	if err != nil {
		fmt.Println("Error checking existing organizations:", err)
		return &errs.Error{
			Code:    errs.Internal,
			Message: "Error al crear la organización de la empresa",
		}
	}

	// Only is posible create 4 organizations per user
	if len(existingOrgs) >= 4 {
		return &errs.Error{
			Code:    errs.FailedPrecondition,
			Message: "No es posible crear más organizaciones. Límite alcanzado.",
		}
	}

	err = b.store.CreateUserOrganization(ctx, userID, name, "company")
	if err != nil {
		return &errs.Error{
			Code:    errs.Internal,
			Message: "Error al crear la organización de la empresa",
		}
	}
	return nil
}
