package api

import (
	"context"
)

// IsOrganizationOwner checks if a user is the owner of an organization
func (a *Api) IsOrganizationOwner(ctx context.Context, orgID string, userID string) (bool, error) {
	return a.db.IsOrganizationOwner(orgID, userID)
}
