package api

import (
	"context"

	"ems.dev/backend/services/organization/types"
)

// GetMemberOrganizations returns all organizations a user is part of, with ownership information
func (a *Api) GetMemberOrganizations(ctx context.Context, userID string) ([]types.Organization, error) {
	return a.db.GetMemberOrganizations(userID)
}

// GetOrganizationByID returns an organization by its ID
func (a *Api) GetOrganizationByID(ctx context.Context, id string) (*types.Organization, error) {
	return a.db.GetOrganizationByID(id)
}
