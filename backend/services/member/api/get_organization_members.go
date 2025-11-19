package api

import (
	"context"

	"ems.dev/backend/services/member/types"
)

// GetOrganizationMembers returns all members of an organization
func (a *Api) GetOrganizationMembers(ctx context.Context, orgID string, params *types.OrganizationMemberParams) ([]types.OrganizationMember, error) {
	members, err := a.db.GetOrganizationMembers(orgID, params)
	if err != nil {
		return nil, err
	}

	// Populate manager information for each member
	for i := range members {
		manager, err := a.directsApi.GetMemberManager(ctx, members[i].ID, orgID)
		if err == nil && manager != nil {
			members[i].ManagerID = &manager.ManagerMemberID
		}
	}

	return members, nil
}
