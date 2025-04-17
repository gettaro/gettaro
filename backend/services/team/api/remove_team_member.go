package api

import (
	"context"
)

func (a *Api) RemoveTeamMember(ctx context.Context, teamID, userID string) error {
	return a.db.RemoveTeamMember(ctx, teamID, userID)
}
