package database

import (
	"context"

	"ems.dev/backend/services/sourcecontrol/types"
)

// SourceControlDB defines the interface for source control database operations
type SourceControlDB interface {
	// CreatePullRequest creates a new pull request
	CreatePullRequest(ctx context.Context, pr *types.PullRequest) error
	// CreatePRComment creates a new pull request comment
	CreatePRComment(ctx context.Context, comment *types.PRComment) error
}
