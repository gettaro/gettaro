package api

import (
	"context"
)

func (a *Api) DeleteTeam(ctx context.Context, id string) error {
	return a.db.DeleteTeam(ctx, id)
}
