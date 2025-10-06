package appstore

import (
	"context"

	"encore.dev/types/uuid"
)

func (s *StoreApp) GetRoleIDByName(ctx context.Context, roleName string) (uuid.UUID, error) {
	q := `
		SELECT id FROM roles WHERE name = $1
	`
	var roleID uuid.UUID
	if err := s.dbx.GetContext(ctx, &roleID, q, roleName); err != nil {
		return uuid.Nil, err
	}
	return roleID, nil
}
