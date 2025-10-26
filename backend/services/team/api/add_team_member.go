package api

import (
	"context"
	"errors"

	"ems.dev/backend/services/team/types"
)

func (a *Api) AddTeamMember(ctx context.Context, teamID string, member *types.TeamMember) error {
	// Validate that the team exists
	team, err := a.db.GetTeam(ctx, teamID)
	if err != nil {
		return err
	}
	if team == nil {
		return errors.New("team not found")
	}

	// Set the team ID on the member
	member.TeamID = teamID

	return a.db.AddTeamMember(ctx, teamID, member)
}
