package api

import (
	"context"

	"ems.dev/backend/services/title/types"
)

// GetTitle retrieves a title by ID
func (s *Api) GetTitle(ctx context.Context, id string) (*types.Title, error) {
	return s.db.GetTitle(id)
}
