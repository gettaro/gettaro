package github

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"ems.dev/backend/libraries/github"
	"ems.dev/backend/services/integration/api"
	"ems.dev/backend/services/integration/types"
	internaltypes "ems.dev/backend/services/sourcecontrol/types"
	"gorm.io/datatypes"
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

		githubUsers := make(map[string]string)
		prsToSave := make([]*internaltypes.PullRequest, 0)
		commentsToSave := make([]*internaltypes.PRComment, 0)

		// Map pull requests
		for _, pr := range prs {
			// Get the pull request details
			prDetails, err := p.githubClient.GetPullRequest(ctx, owner, repoName, token, pr.Number)
			if err != nil {
				return fmt.Errorf("failed to fetch pull request details for %s: %w", repo, err)
			}

			githubUsers[prDetails.User.Login] = prDetails.User.Login

			// Map GitHub PR to our PullRequest type
			prDetailsBytes, _ := json.Marshal(prDetails)
			sourceControlPR := &internaltypes.PullRequest{
				SourceControlAccountID: config.ID,
				ProviderID:             fmt.Sprintf("%d", prDetails.ID),
				URL:                    prDetails.URL,
				Title:                  prDetails.Title,
				Description:            prDetails.Body,
				Status:                 prDetails.State,
				CreatedAt:              prDetails.CreatedAt,
				MergedAt:               prDetails.MergedAt,
				LastUpdatedAt:          prDetails.UpdatedAt,
				Comments:               prDetails.Comments,
				ReviewComments:         prDetails.ReviewComments,
				Additions:              prDetails.Additions,
				Deletions:              prDetails.Deletions,
				ChangedFiles:           prDetails.ChangedFiles,
				Metadata:               datatypes.JSON(prDetailsBytes),
			}

			prsToSave = append(prsToSave, sourceControlPR)

			// Fetch review comments for this PR
			comments, err := p.githubClient.GetPullRequestReviewComments(ctx, owner, repoName, token, pr.Number)
			if err != nil {
				return fmt.Errorf("failed to fetch review comments for PR %d: %w", pr.Number, err)
			}

			// Map review comments
			for _, comment := range comments {
				sourceControlComment := &internaltypes.PRComment{
					PRID:      sourceControlPR.ID,
					Body:      comment.Body,
					CreatedAt: comment.CreatedAt,
					UpdatedAt: &comment.UpdatedAt,
				}

				githubUsers[comment.User.Login] = comment.User.Login

				commentsToSave = append(commentsToSave, sourceControlComment)
			}

			// Here you would typically save the PR to the database
			// For now, we just log it
			fmt.Printf("Mapped PR %d: %s\n", pr.Number, pr.Title)
		}

		// Fetch source control accounts for the github users using the username
		// Identify the source control accounts that need to be created
		// Create the source control accounts
		// Map the pull requests to the source control accounts
		// Save the pull requests
		// Map the comments to the source control accounts
		// Save the comments
	}

	return nil
}
