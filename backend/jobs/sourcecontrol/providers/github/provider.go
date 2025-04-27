package github

import (
	"context"
	"fmt"
	"strings"

	"ems.dev/backend/libraries/github"
	"ems.dev/backend/services/integration/api"
	"ems.dev/backend/services/integration/types"
	internaltypes "ems.dev/backend/services/sourcecontrol/types"
)

type GitHubProvider struct {
	githubClient   *github.Client
	integrationAPI api.IntegrationAPI
}

func NewProvider(githubClient *github.Client, integrationAPI api.IntegrationAPI) *GitHubProvider {
	return &GitHubProvider{
		githubClient:   githubClient,
		integrationAPI: integrationAPI,
	}
}

func (p *GitHubProvider) Name() string {
	return "github"
}

func (p *GitHubProvider) SyncRepositories(ctx context.Context, config *types.IntegrationConfig, repositories []string) error {
	// Decrypt the token
	token, err := p.integrationAPI.DecryptToken(config.EncryptedToken)
	if err != nil {
		return fmt.Errorf("failed to decrypt token: %w", err)
	}

	// 1. For each repository:
	for _, repo := range repositories {
		parts := strings.Split(repo, "/")
		if len(parts) != 2 {
			return fmt.Errorf("invalid repository format: %s. Expected format: owner/repo", repo)
		}
		owner, repoName := parts[0], parts[1]

		// Fetch pull requests
		prs, err := p.githubClient.GetPullRequests(ctx, owner, repoName, token)
		if err != nil {
			return fmt.Errorf("failed to fetch pull requests for %s: %w", repo, err)
		}

		// Map pull requests
		for _, pr := range prs {

			// Get the pull request details
			prDetails, err := p.githubClient.GetPullRequest(ctx, owner, repoName, token, pr.Number)
			if err != nil {
				return fmt.Errorf("failed to fetch pull request details for %s: %w", repo, err)
			}

			fmt.Printf("prDetails: %+v\n", prDetails)
			// Map GitHub PR to our PullRequest type
			sourceControlPR := &internaltypes.PullRequest{
				SourceControlAccountID: config.ID,
				ProviderID:             fmt.Sprintf("%d", pr.ID),
				URL:                    pr.URL,
				Title:                  pr.Title,
				Status:                 pr.State,
				CreatedAt:              pr.CreatedAt,
				MergedAt:               pr.MergedAt,
				LastUpdatedAt:          pr.UpdatedAt,
			}

			// Fetch review comments for this PR
			comments, err := p.githubClient.GetPullRequestReviewComments(ctx, owner, repoName, token, pr.Number)
			if err != nil {
				return fmt.Errorf("failed to fetch review comments for PR %d: %w", pr.Number, err)
			}

			// Map review comments
			for _, comment := range comments {
				_ = &internaltypes.PRComment{
					PRID:      sourceControlPR.ID,
					AuthorID:  &comment.User.Login,
					Body:      comment.Body,
					CreatedAt: comment.CreatedAt,
					UpdatedAt: &comment.UpdatedAt,
				}

				// Here you would typically save the comment to the database
				// For now, we just log it
				fmt.Printf("Mapped comment %d for PR %d\n", comment.ID, pr.Number)
			}

			// Here you would typically save the PR to the database
			// For now, we just log it
			fmt.Printf("Mapped PR %d: %s\n", pr.Number, pr.Title)
		}
	}

	return nil
}
