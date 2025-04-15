package api

import (
	"context"

	"ems.dev/backend/services/organization/types"
)

// UpdateOrganization updates an existing organization
func (a *Api) UpdateOrganization(ctx context.Context, org *types.Organization) error {
	return a.db.UpdateOrganization(org)
}
