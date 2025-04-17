package api

import (
	"context"

	"ems.dev/backend/services/team/types"
)

func (a *Api) GetTeam(ctx context.Context, id string) (*types.Team, error) {
	return a.db.GetTeam(ctx, id)
}
