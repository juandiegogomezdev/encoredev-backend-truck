package appbusiness

import (
	"context"

	"encore.dev/types/uuid"
)

func (b *BusinessApp) GetAllUserMemberships(ctx context.Context, userID uuid.UUID) ([]Membership, error) {
	return b.store.GetAllUserOrganizations(ctx, userID)
}
