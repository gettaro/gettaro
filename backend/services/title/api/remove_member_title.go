package api

import (
	"context"
)

// RemoveMemberTitle removes a member's title assignment
func (s *Api) RemoveMemberTitle(ctx context.Context, memberID string, orgID string) error {
	return s.db.RemoveMemberTitle(memberID, orgID)
}
