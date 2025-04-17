package api

import (
	"context"

	"ems.dev/backend/services/team/types"
)

func (a *Api) AddTeamMember(ctx context.Context, teamID string, member *types.TeamMember) error {
	return a.db.AddTeamMember(ctx, teamID, member)
}
