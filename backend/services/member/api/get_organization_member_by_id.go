package api

import (
	"context"

	"ems.dev/backend/services/member/types"
)

// GetOrganizationMemberByID retrieves a member by their ID
func (a *Api) GetOrganizationMemberByID(ctx context.Context, memberID string) (*types.OrganizationMember, error) {
	return a.db.GetOrganizationMemberByID(ctx, memberID)
}
