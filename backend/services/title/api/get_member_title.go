package api

import (
	"context"

	"ems.dev/backend/services/title/types"
)

// GetMemberTitle retrieves a member's title assignment
func (s *Api) GetMemberTitle(ctx context.Context, memberID string, orgID string) (*types.MemberTitle, error) {
	return s.db.GetMemberTitle(memberID, orgID)
}
