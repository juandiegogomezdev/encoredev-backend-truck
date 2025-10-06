package appstore

import (
	"context"

	"encore.app/appservice/shared"
	"encore.dev/types/uuid"
)

// Check if a user has an active membership in an organization
func (s *StoreApp) HasActiveMembership(ctx context.Context, userID uuid.UUID, orgID uuid.UUID) (bool, error) {
	q := `
		SELECT EXISTS (
			SELECT 1 FROM org_memberships
			WHERE user_id = $1 AND org_id = $2 AND status = 'active'
		)
	`
	var exists bool
	if err := s.dbx.GetContext(ctx, &exists, q, userID, orgID); err != nil {
		return false, err
	}
	return exists, nil
}

// Get all memberships of a user
func (s *StoreApp) GetUserMemberships(ctx context.Context, userID uuid.UUID) ([]shared.Membership, error) {

	q := `
		SELECT
			om.id,
			om.status,
			om.created_at,
			om.finalized_at,
			o.name AS org_name,
			o.type AS org_type,
			r.name AS role_name
		FROM org_memberships om
		JOIN organizations o ON om.org_id = o.id
		JOIN roles r ON om.role_id = r.id
		WHERE om.user_id = $1
	`
	memberships := []shared.Membership{}
	if err := s.dbx.SelectContext(ctx, &memberships, q, userID); err != nil {
		return nil, err
	}

	return memberships, nil
}

// Get the id of a user by membership
func (s *StoreApp) GetUserIDByMembership(ctx context.Context, memID uuid.UUID) (uuid.UUID, error) {
	q := `
		SELECT user_id FROM org_memberships WHERE id = $1
	`
	var userID uuid.UUID
	if err := s.dbx.GetContext(ctx, &userID, q, memID); err != nil {
		return uuid.UUID{}, err
	}
	return userID, nil
}

// Create a new organization membership for a user
func (s *StoreApp) CreateOrgMembership(ctx context.Context, newMembership CreateOrgMembershipStruct) error {
	q := `
		INSERT INTO org_memberships (id, org_id, user_id, role_id, status, created_by)
		VALUES (:id, :org_id, :user_id, :role_id, :status, :created_by)
		`

	_, err := s.dbx.NamedExecContext(ctx, q, newMembership)
	return err
}

type CreateOrgMembershipStruct struct {
	ID        uuid.UUID `db:"id"`
	OrgID     uuid.UUID `db:"org_id"`
	UserID    uuid.UUID `db:"user_id"`
	RoleID    uuid.UUID `db:"role_id"`
	Status    string    `db:"status"`
	CreatedBy uuid.UUID `db:"created_by"`
}

// Get a membership status by ID
func (s *StoreApp) GetMembershipStatus(ctx context.Context, memID uuid.UUID) (string, error) {
	q := `
		SELECT status FROM org_memberships WHERE id = $1
	`
	var status string
	if err := s.dbx.GetContext(ctx, &status, q, memID); err != nil {
		return "", err
	}
	return status, nil
}

// Change the status of a membership
func (s *StoreApp) ChangeMembershipStatus(ctx context.Context, memID uuid.UUID, status string) error {
	q := `
		UPDATE org_memberships SET status = $1 WHERE id = $2
	`
	_, err := s.dbx.ExecContext(ctx, q, status, memID)
	return err
}

// Finalize the membership
func (s *StoreApp) FinalizeMembership(ctx context.Context, finalizeMembership FinalizeMembershipStruct) error {
	q := `
		UPDATE org_memberships
		SET status = 'ended', finalized_by = $1, finalized_at = now()
		WHERE id = $2
	`
	_, err := s.dbx.NamedExecContext(ctx, q, finalizeMembership)
	return err
}

type FinalizeMembershipStruct struct {
	MemID       uuid.UUID `db:"mem_id"`
	FinalizedBy uuid.UUID `db:"finalized_by"`
}
