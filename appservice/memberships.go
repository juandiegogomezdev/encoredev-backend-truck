package appService

import (
	"context"

	"encore.app/appservice/shared"
)

// Send the memberships of an user
// encore:api auth method=GET path=/memberships
func (s *ServiceApp) GetUserMemberships(ctx context.Context) (shared.Membership, error) {
	return shared.Membership{}, nil
}
