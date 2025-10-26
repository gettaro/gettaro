package api

import (
	"context"
	"errors"
)

func (a *Api) RemoveTeamMember(ctx context.Context, teamID, memberID string) error {
	// Validate that the team exists
	team, err := a.db.GetTeam(ctx, teamID)
	if err != nil {
		return err
	}
	if team == nil {
		return errors.New("team not found")
	}

	return a.db.RemoveTeamMember(ctx, teamID, memberID)
}
