package appstore

import (
	"context"
	"time"

	"encore.dev/types/uuid"
)

// Get all organizations for a user
func (s *StoreApp) GetAllUserOrganizations(ctx context.Context, userID uuid.UUID) ([]ResUserOrganizationStore, error) {

	query := `
		SELECT id, name, type, created_at
		FROM organizations
		WHERE user_id = $1
	`
	var organizations []ResUserOrganizationStore
	if err := s.dbx.SelectContext(ctx, &organizations, query, userID); err != nil {
		return nil, err
	}
	return organizations, nil
}

type ResUserOrganizationStore struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

// Create a new organization for a user
func (s *StoreApp) CreateUserOrganization(ctx context.Context, ownerId uuid.UUID, name string, orgType string) error {
	var createdAt time.Time
	var orgId uuid.UUID
	query := `
		INSERT INTO organizations (owner_id, name, type)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	err := s.db.QueryRow(ctx, query, ownerId, name, orgType).Scan(&orgId, &createdAt)
	if err != nil {
		return err
	}
	return nil
}

// Create and organization and assign the owner as admin
func (s *StoreApp) CreateOrgAndMembership(ctx context.Context, org CreateOrganizationStruct, membership CreateOrgMembershipStruct) error {
	tx, err := s.dbx.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	orgQuery := `
		INSERT INTO organizations (id, owner_id, name, type)
		VALUES (:id, :owner_id, :name, :type)
		`

	if _, err := tx.NamedExecContext(ctx, orgQuery, org); err != nil {
		tx.Rollback()
		return err
	}

	// Find the role ID for "owner"
	var roleID uuid.UUID
	roleQuery := `SELECT id FROM roles WHERE name = 'owner'`
	if err := tx.GetContext(ctx, &roleID, roleQuery); err != nil {
		tx.Rollback()
		return err
	}

	membership.RoleID = roleID

	membershipQuery := `
		INSERT INTO org_memberships (id, org_id, user_id, role_id, status, created_by)
		VALUES (:id, :org_id, :user_id, :role_id, :status, :created_by)
	`

	if _, err := tx.NamedExecContext(ctx, membershipQuery, membership); err != nil {
		tx.Rollback()
		return err
	}

	return nil

}

type CreateOrganizationStruct struct {
	OrgID   uuid.UUID `db:"id"`
	OwnerID uuid.UUID `db:"owner_id"`
	Name    string    `db:"name"`
	Type    string    `db:"type"`
}
type CreateOwnerMembershipStruct struct {
	ID        uuid.UUID `db:"id"`
	OrgID     uuid.UUID `db:"org_id"`
	UserID    uuid.UUID `db:"user_id"`
	RoleID    uuid.UUID `db:"role_id"`
	Status    string    `db:"status"`
	CreatedBy uuid.UUID `db:"created_by"`
}
