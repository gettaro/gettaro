package api

import (
	"context"

	"ems.dev/backend/services/title/types"
)

// UpdateTitle updates an existing title
func (s *Api) UpdateTitle(ctx context.Context, id string, request types.Title) (*types.Title, error) {
	title, err := s.db.GetTitle(id)
	if err != nil {
		return nil, err
	}
	if title == nil {
		return nil, nil
	}

	title.Name = request.Name
	title.IsManager = request.IsManager
	err = s.db.UpdateTitle(*title)
	if err != nil {
		return nil, err
	}

	return title, nil
}
