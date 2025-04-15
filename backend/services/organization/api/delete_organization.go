package api

import "context"

// DeleteOrganization deletes an organization and its relationships
func (a *Api) DeleteOrganization(ctx context.Context, id string) error {
	return a.db.DeleteOrganization(id)
}
