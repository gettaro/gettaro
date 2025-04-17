package api

import (
	"context"

	"ems.dev/backend/services/team/types"
)

func (a *Api) ListTeams(ctx context.Context, params types.TeamSearchParams) ([]types.Team, error) {
	return a.db.ListTeams(ctx, params)
}
