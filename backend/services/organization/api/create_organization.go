package api

import (
	"context"

	"ems.dev/backend/services/organization/types"
)

// CreateOrganization creates a new organization and sets the specified user as its owner
func (a *Api) CreateOrganization(ctx context.Context, org *types.Organization, userID string) error {
	return a.db.CreateOrganization(org, userID)
}
