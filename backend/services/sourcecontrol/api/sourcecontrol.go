package api

import (
	"context"

	"ems.dev/backend/services/sourcecontrol/database"
	"ems.dev/backend/services/sourcecontrol/types"
)

type SourceControlAPI interface {
	// Source Control Accounts
	GetSourceControlAccountsByUsernames(ctx context.Context, usernames []string) (map[string]*types.SourceControlAccount, error)
	CreateSourceControlAccounts(ctx context.Context, accounts []*types.SourceControlAccount) error

	// Pull Requests
	GetPullRequests(ctx context.Context, params *types.PullRequestParams) ([]*types.PullRequest, error)
	CreatePullRequests(ctx context.Context, prs []*types.PullRequest) error

	// Comments
	CreatePRComments(ctx context.Context, comments []*types.PRComment) error
}

type Api struct {
	db database.DB
}

func NewAPI(db database.DB) SourceControlAPI {
	return &Api{
		db: db,
	}
}

func (a *Api) GetSourceControlAccountsByUsernames(ctx context.Context, usernames []string) (map[string]*types.SourceControlAccount, error) {
	return a.db.GetSourceControlAccountsByUsernames(ctx, usernames)
}

func (a *Api) CreateSourceControlAccounts(ctx context.Context, accounts []*types.SourceControlAccount) error {
	return a.db.CreateSourceControlAccounts(ctx, accounts)
}

func (a *Api) GetPullRequests(ctx context.Context, params *types.PullRequestParams) ([]*types.PullRequest, error) {
	return a.db.GetPullRequests(ctx, params)
}

func (a *Api) CreatePullRequests(ctx context.Context, prs []*types.PullRequest) error {
	return a.db.CreatePullRequests(ctx, prs)
}

func (a *Api) CreatePRComments(ctx context.Context, comments []*types.PRComment) error {
	return a.db.CreatePRComments(ctx, comments)
}
