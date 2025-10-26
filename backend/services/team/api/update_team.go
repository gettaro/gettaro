package api

import (
	"context"
	"errors"

	"ems.dev/backend/services/team/types"
)

func (a *Api) UpdateTeam(ctx context.Context, id string, team *types.Team) error {
	// Validate that the team exists
	existingTeam, err := a.db.GetTeam(ctx, id)
	if err != nil {
		return err
	}
	if existingTeam == nil {
		return errors.New("team not found")
	}

	return a.db.UpdateTeam(ctx, id, team)
}
