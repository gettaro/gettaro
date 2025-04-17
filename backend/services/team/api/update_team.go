package api

import (
	"context"

	"ems.dev/backend/services/team/types"
)

func (a *Api) UpdateTeam(ctx context.Context, id string, team *types.Team) error {
	return a.db.UpdateTeam(ctx, id, team)
}
