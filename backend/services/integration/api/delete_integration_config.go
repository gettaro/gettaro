package api

import (
	"context"
)

// DeleteIntegrationConfig deletes an integration config
func (a *Api) DeleteIntegrationConfig(ctx context.Context, id string) error {
	return a.db.DeleteIntegrationConfig(id)
}
