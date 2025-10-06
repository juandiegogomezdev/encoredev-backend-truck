package shared

import "encore.dev/types/uuid"

type CreateOrganizationStruct struct {
	OrgID   uuid.UUID `db:"id"`
	OwnerID uuid.UUID `db:"owner_user_id"`
	Name    string    `db:"name"`
	Type    string    `db:"type"`
}
type CreateOwnerMembershipStruct struct {
	MemID     uuid.UUID `db:"id"`
	OrgID     uuid.UUID `db:"org_id"`
	UserID    uuid.UUID `db:"user_id"`
	RoleID    uuid.UUID `db:"role_id"`
	Status    string    `db:"status"`
	CreatedBy uuid.UUID `db:"created_by"`
}
