package api

import (
	"context"

	"ems.dev/backend/services/team/types"
)

func (a *Api) CreateTeam(ctx context.Context, team *types.Team) error {
	return a.db.CreateTeam(ctx, team)
}
