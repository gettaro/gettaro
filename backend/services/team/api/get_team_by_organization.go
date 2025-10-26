package api

import (
	"context"
	"errors"

	"ems.dev/backend/services/team/types"
)

func (a *Api) GetTeamByOrganization(ctx context.Context, teamID, organizationID string) (*types.Team, error) {
	team, err := a.db.GetTeam(ctx, teamID)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, errors.New("team not found")
	}

	// Validate that the team belongs to the organization
	if team.OrganizationID != organizationID {
		return nil, errors.New("team not found")
	}

	return team, nil
}
