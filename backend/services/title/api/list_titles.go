package api

import (
	"context"

	"ems.dev/backend/services/title/types"
)

// ListTitles retrieves all titles for an organization
func (s *Api) ListTitles(ctx context.Context, orgID string) ([]types.Title, error) {
	return s.db.ListTitles(orgID)
}
