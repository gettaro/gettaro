package api

import (
	"context"
	"errors"

	"ems.dev/backend/services/team/types"
)

func (a *Api) CreateTeam(ctx context.Context, team *types.Team) error {
	if team.OrganizationID == "" {
		return errors.New("organization ID is required")
	}

	org, err := a.orgApi.GetOrganizationByID(ctx, team.OrganizationID)
	if err != nil {
		return err
	}

	if org == nil {
		return errors.New("organization not found")
	}

	return a.db.CreateTeam(ctx, team)
}
