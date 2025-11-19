package api

import (
	"context"

	"ems.dev/backend/services/title/types"
)

// AssignMemberTitle assigns a title to a member within an organization
func (s *Api) AssignMemberTitle(ctx context.Context, memberTitle types.MemberTitle) error {
	return s.db.AssignMemberTitle(memberTitle)
}
