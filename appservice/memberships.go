package appService

import (
	"context"

	"encore.app/appservice/shared"
	"encore.dev/beta/auth"
)

// Send the memberships of an user
// encore:api auth method=GET path=/memberships
func (s *ServiceApp) GetUserMemberships(ctx context.Context) (GetUserMembershipsResponse, error) {
	memberships, err := s.b.GetAllUserMemberships(ctx, auth.Data().(*AuthData).UserID)
	if err != nil {
		return GetUserMembershipsResponse{}, err
	}
	return GetUserMembershipsResponse{Memberships: memberships}, nil
}

type GetUserMembershipsResponse struct {
	Memberships []shared.Membership `json:"memberships"`
}
