package appstore

import (
	"context"
	"time"

	"encore.dev/types/uuid"
)

// Get all organizations for a user
func (s *AppStore) GetAllUserOrganizations(ctx context.Context, userID uuid.UUID) ([]resUserOrganizationStore, error) {

	query := `
		SELECT id, name, type, created_at
		FROM organizations
		WHERE user_id = $1
	`
	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var organizations []resUserOrganizationStore
	for rows.Next() {
		var org resUserOrganizationStore
		if err := rows.Scan(&org.ID, &org.Name, &org.Type, &org.CreatedAt); err != nil {
			return nil, err
		}
		organizations = append(organizations, org)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return organizations, nil
}

type resUserOrganizationStore struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

// Create a new organization for a user
func (s *AppStore) CreateUserOrganization(ctx context.Context, ownerId uuid.UUID, name string, orgType string) (resCreateUserOrganizationStore, error) {
	var createdAt time.Time
	var orgId uuid.UUID
	query := `
		INSERT INTO organizations (owner_id, name, type)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	err := s.db.QueryRow(ctx, query, ownerId, name, orgType).Scan(&orgId, &createdAt)
	if err != nil {
		return resCreateUserOrganizationStore{}, err
	}
	return resCreateUserOrganizationStore{
		ID:        orgId,
		CreatedAt: createdAt,
	}, nil
}

type resCreateUserOrganizationStore struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

//
