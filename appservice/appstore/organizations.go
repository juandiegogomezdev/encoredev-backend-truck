package appstore

import (
	"context"
	"fmt"
	"time"

	"encore.app/appservice/shared"
	"encore.dev/types/uuid"
)

// Get all organizations for a user
func (s *StoreApp) GetAllUserOrganizations(ctx context.Context, userID uuid.UUID) ([]ResUserOrganizationStore, error) {

	query := `
		SELECT id, name, type, created_at
		FROM organizations
		WHERE owner_user_id = $1
	`
	var organizations []ResUserOrganizationStore
	if err := s.dbx.SelectContext(ctx, &organizations, query, userID); err != nil {
		return nil, err
	}
	return organizations, nil
}

type ResUserOrganizationStore struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Type      string    `json:"type" db:"type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Create and organization and assign the owner as admin
func (s *StoreApp) CreateOrgAndMembership(ctx context.Context, newOrg shared.CreateOrganizationStruct, newMembership shared.CreateOwnerMembershipStruct) error {
	tx, err := s.dbx.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	orgQuery := `
		INSERT INTO organizations (id, owner_user_id, name, type)
		VALUES (:id, :owner_user_id, :name, :type)
		`

	if _, err := tx.NamedExecContext(ctx, orgQuery, newOrg); err != nil {
		tx.Rollback()
		return fmt.Errorf("error creating organization: %w", err)
	}

	// Find the role ID for "owner"
	var roleID uuid.UUID
	roleQuery := `SELECT id FROM roles WHERE name = 'owner'`
	if err := tx.GetContext(ctx, &roleID, roleQuery); err != nil {
		tx.Rollback()
		return fmt.Errorf("error getting role ID: %w", err)
	}

	newMembership.RoleID = roleID

	membershipQuery := `
		INSERT INTO org_memberships (id, org_id, user_id, role_id, status, created_by)
		VALUES (:id, :org_id, :user_id, :role_id, :status, :created_by)
	`

	if _, err := tx.NamedExecContext(ctx, membershipQuery, newMembership); err != nil {
		tx.Rollback()
		return fmt.Errorf("error creating membership: %w", err)
	}

	tx.Commit()

	return nil

}
