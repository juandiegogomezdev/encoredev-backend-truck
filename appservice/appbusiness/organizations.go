package appbusiness

import (
	"context"
	"fmt"
	"time"

	"encore.app/appservice/shared"
	"encore.app/pkg/utils"
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

// Create a personal organization and return the membership
func (b *BusinessApp) CreatePersonalOrganization(ctx context.Context, userID uuid.UUID, name string) (shared.Membership, error) {
	// Get all organizations for this user
	existingOrgs, err := b.store.GetAllUserOrganizations(ctx, userID)
	if err != nil {
		fmt.Println("Error checking existing organizations:", err)
		return shared.Membership{}, &errs.Error{
			Code:    errs.Internal,
			Message: "Error al crear la organización personal",
		}
	}

	// Only is posible create 1 personal organization if not exists
	if len(existingOrgs) != 0 {
		return shared.Membership{}, &errs.Error{
			Code:    errs.FailedPrecondition,
			Message: "No es posible crear más organizaciones personales. Límite alcanzado.",
		}
	}

	// Generate the new organization ID and membership ID
	orgID, err := utils.MustNewUUID()
	if err != nil {
		return shared.Membership{}, err
	}

	memID, err := utils.MustNewUUID()
	if err != nil {
		return shared.Membership{}, err
	}

	// Create the personal organization struct
	createOrganization := shared.CreateOrganizationStruct{
		Name:    name,
		OrgID:   orgID,
		Type:    "personal",
		OwnerID: userID,
	}

	// Create the membership organization struct
	createMembership := shared.CreateOwnerMembershipStruct{
		MemID:     memID,
		OrgID:     orgID,
		UserID:    userID,
		Status:    "active",
		CreatedBy: userID,
		RoleID:    uuid.Nil,
	}

	// Create the organization and the membership in a transaction
	err = b.store.CreateOrgAndMembership(ctx, createOrganization, createMembership)
	if err != nil {
		fmt.Println("Error creating organization and membership:", err)
		return shared.Membership{}, &errs.Error{
			Code:    errs.Internal,
			Message: "Error al crear el trabajo como independiente",
		}
	}

	membership := shared.Membership{
		ID:        memID,
		Status:    "active",
		CreatedAt: time.Now(),
		OrgName:   name,
		RoleName:  "owner",
		OrgType:   "personal",
	}
	return membership, nil
}

// Create a company organization
func (b *BusinessApp) CreateCompanyOrganization(ctx context.Context, userID uuid.UUID, name string) (shared.Membership, error) {
	// Get all organizations for this user
	existingOrgs, err := b.store.GetAllUserOrganizations(ctx, userID)
	if err != nil {
		fmt.Println("Error checking existing organizations:", err)
		return shared.Membership{}, &errs.Error{
			Code:    errs.Internal,
			Message: "Error al crear la organización de la empresa",
		}
	}

	// Only is posible create 4 organizations per user and the user need to have at least the
	// personal organization created
	if len(existingOrgs) >= 4 && len(existingOrgs) != 0 {
		return shared.Membership{}, &errs.Error{
			Code:    errs.FailedPrecondition,
			Message: "No es posible crear más organizaciones. Límite alcanzado.",
		}
	}

	// Generate the new organization ID and membership ID
	orgID, err := utils.MustNewUUID()
	if err != nil {
		return shared.Membership{}, err
	}
	memID, err := utils.MustNewUUID()
	if err != nil {
		return shared.Membership{}, err
	}

	// Create the company organization struct
	createOrganization := shared.CreateOrganizationStruct{
		Name:    name,
		OrgID:   orgID,
		Type:    "company",
		OwnerID: userID,
	}

	// Create the membership organization struct
	createMembership := shared.CreateOwnerMembershipStruct{
		MemID:     memID,
		OrgID:     orgID,
		UserID:    userID,
		Status:    "active",
		CreatedBy: userID,
		RoleID:    uuid.Nil,
	}

	// Create the organization and the membership in a transaction
	err = b.store.CreateOrgAndMembership(ctx, createOrganization, createMembership)
	if err != nil {
		return shared.Membership{}, &errs.Error{
			Code:    errs.Internal,
			Message: "Error al crear el trabajo como independiente",
		}
	}

	membership := shared.Membership{
		ID:        memID,
		Status:    "active",
		CreatedAt: time.Now(),
		OrgName:   name,
		RoleName:  "owner",
		OrgType:   "company",
	}
	return membership, nil
}
