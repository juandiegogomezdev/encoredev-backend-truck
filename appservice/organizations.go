package appService

import (
	"context"

	"encore.app/appservice/shared"
	"encore.dev/beta/auth"
)

//encore:api public method=GET path=/org/hello
func (s *ServiceApp) Hello(ctx context.Context) (*responseHello, error) {
	return &responseHello{Message: "Hello, World!"}, nil
}

//encore:api public method=GET path=/org
func (s *ServiceApp) GetAllOrganizations(ctx context.Context) (responseGetAllOrganizations, error) {
	return responseGetAllOrganizations{}, nil
}

type responseGetAllOrganizations struct {
}

type responseHello struct {
	Message string `json:"message"`
}

// Create a personal organization for the user
// encore:api auth method=POST path=/org/personal
func (s *ServiceApp) CreatePersonalOrg(ctx context.Context) (*CreatePersonalOrgResponse, error) {
	data := auth.Data().(*AuthData)
	membership, err := s.b.CreatePersonalOrganization(ctx, data.UserID, "Trabajo independiente")
	if err != nil {
		return nil, err
	}
	return &CreatePersonalOrgResponse{Memberships: membership}, nil
}

type CreatePersonalOrgResponse struct {
	Memberships shared.Membership `json:"memberships"`
}
