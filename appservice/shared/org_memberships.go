package shared

import (
	"time"

	"encore.dev/types/uuid"
)

type Membership struct {
	ID          uuid.UUID  `json:"id" db:"id"`                     // Id of the membership
	Status      string     `json:"status" db:"status"`             // Status of the membership (active, revoked, etc)
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`     // Creation time of the membership
	FinalizedAt *time.Time `json:"finalized_at" db:"finalized_at"` // Finalization time of the membership
	OrgName     string     `json:"org_name" db:"org_name"`         // Name of the organization where the membership belongs
	OrgType     string     `json:"org_type" db:"org_type"`         // Type of the organization (e.g., company, team) where the membership belongs
	RoleName    string     `json:"role_name" db:"role_name"`       // Role of the user in the organization
}
