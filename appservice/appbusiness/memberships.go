package appbusiness

import (
	"context"

	"encore.app/appservice/shared"
	"encore.dev/beta/errs"
	"encore.dev/types/uuid"
)

func (b *BusinessApp) GetAllUserMemberships(ctx context.Context, userID uuid.UUID) ([]shared.Membership, error) {
	memberships, err := b.store.GetUserMemberships(ctx, userID)
	if err != nil {
		return nil, &errs.Error{
			Code:    errs.Internal,
			Message: "failed to get user memberships",
		}
	}

	return memberships, nil
}
