package api

import (
	"context"
)

// RemoveOrganizationMember removes a user from an organization
func (a *Api) RemoveOrganizationMember(ctx context.Context, orgID string, userID string) error {
	return a.db.RemoveOrganizationMember(orgID, userID)
}
