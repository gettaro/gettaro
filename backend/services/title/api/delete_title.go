package api

import (
	"context"
)

// DeleteTitle deletes a title
func (s *Api) DeleteTitle(ctx context.Context, id string) error {
	return s.db.DeleteTitle(id)
}
