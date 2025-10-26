package api

import (
	"context"
	"errors"
)

func (a *Api) DeleteTeam(ctx context.Context, id string) error {
	// Validate that the team exists
	team, err := a.db.GetTeam(ctx, id)
	if err != nil {
		return err
	}
	if team == nil {
		return errors.New("team not found")
	}

	return a.db.DeleteTeam(ctx, id)
}
