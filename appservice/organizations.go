package appService

import "context"

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
