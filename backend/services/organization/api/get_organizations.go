package api

import (
	"context"

	"ems.dev/backend/services/organization/types"
)

// GetUserOrganizations returns all organizations a user is part of, with ownership information
func (a *Api) GetUserOrganizations(ctx context.Context, userID string) ([]types.Organization, error) {
	return a.db.GetUserOrganizations(userID)
}

// GetOrganizationByID returns an organization by its ID
func (a *Api) GetOrganizationByID(ctx context.Context, id string) (*types.Organization, error) {
	return a.db.GetOrganizationByID(id)
}
