package github

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"ems.dev/backend/libraries/github"
	githubtypes "ems.dev/backend/libraries/github/types"
	"ems.dev/backend/services/integration/api"
	"ems.dev/backend/services/integration/types"
	sourcecontrolapi "ems.dev/backend/services/sourcecontrol/api"
	internaltypes "ems.dev/backend/services/sourcecontrol/types"
	"gorm.io/datatypes"
)

type GitHubProvider struct {
	githubClient     *github.Client
	integrationAPI   api.IntegrationAPI
	sourceControlAPI sourcecontrolapi.SourceControlAPI
}

func NewProvider(
	githubClient *github.Client,
	integrationAPI api.IntegrationAPI,
	sourceControlAPI sourcecontrolapi.SourceControlAPI,
) *GitHubProvider {
	return &GitHubProvider{
		githubClient:     githubClient,
		integrationAPI:   integrationAPI,
		sourceControlAPI: sourceControlAPI,
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

		// Fetch imported pull requests
		importedPRs, err := p.sourceControlAPI.GetPullRequests(ctx, &internaltypes.PullRequestParams{
			ProviderID:     fmt.Sprintf("%s/%s", owner, repoName),
			OrganizationID: &config.OrganizationID,
			ProviderName:   "github",
			RepositoryName: repoName,
		})

		if err != nil {
			return fmt.Errorf("failed to fetch imported pull requests for %s: %w", repo, err)
		}

		// Create a map of imported pull requests
		importedPRsMap := make(map[string]internaltypes.PullRequest)
		for _, pr := range importedPRs {
			importedPRsMap[pr.ProviderID] = *pr
		}

		// Fetch pull requests
		prs, err := p.githubClient.GetPullRequests(ctx, owner, repoName, token)
		if err != nil {
			return fmt.Errorf("failed to fetch pull requests for %s: %w", repo, err)
		}

		githubUsers := make(map[string]githubtypes.User)
		prsToSave := make([]*internaltypes.PullRequest, 0)
		commentsToSave := make([]*internaltypes.PRComment, 0)
		commentAuthors := make(map[string]string)

		// Map pull requests
		for _, pr := range prs {
			// If the pull request has already been imported and is not open, skip it
			if importedPR, ok := importedPRsMap[fmt.Sprintf("%d", pr.ID)]; ok {
				if importedPR.Status != "open" {
					continue
				}
			}

			// Get the pull request details
			prDetails, err := p.githubClient.GetPullRequest(ctx, owner, repoName, token, pr.Number)
			if err != nil {
				return fmt.Errorf("failed to fetch pull request details for %s: %w", repo, err)
			}

			githubUsers[prDetails.User.Login] = prDetails.User

			// Map GitHub PR to our PullRequest type
			prDetailsBytes, _ := json.Marshal(prDetails)
			sourceControlPR := &internaltypes.PullRequest{
				ProviderID:     fmt.Sprintf("%d", prDetails.ID),
				URL:            prDetails.URL,
				RepositoryName: repoName,
				OrganizationID: config.OrganizationID,
				Title:          prDetails.Title,
				Description:    prDetails.Body,
				Status:         prDetails.State,
				CreatedAt:      prDetails.CreatedAt,
				MergedAt:       prDetails.MergedAt,
				LastUpdatedAt:  prDetails.UpdatedAt,
				Comments:       prDetails.Comments,
				ReviewComments: prDetails.ReviewComments,
				Additions:      prDetails.Additions,
				Deletions:      prDetails.Deletions,
				ChangedFiles:   prDetails.ChangedFiles,
				Metadata:       datatypes.JSON(prDetailsBytes),
			}

			prsToSave = append(prsToSave, sourceControlPR)

			// Fetch review reviewComments for this PR
			reviewComments, err := p.githubClient.GetPullRequestReviewComments(ctx, owner, repoName, token, pr.Number)
			if err != nil {
				return fmt.Errorf("failed to fetch review comments for PR %d: %w", pr.Number, err)
			}

			// Fetch regular comments for this PR
			comments, err := p.githubClient.GetPullRequestComments(ctx, owner, repoName, token, pr.Number)
			if err != nil {
				return fmt.Errorf("failed to fetch regular comments for PR %d: %w", pr.Number, err)
			}

			reviewComments = append(reviewComments, comments...)

			// Map review comments
			for _, comment := range reviewComments {
				sourceControlComment := &internaltypes.PRComment{
					PRID:       sourceControlPR.ID,
					ProviderID: fmt.Sprintf("%d", comment.ID),
					Body:       comment.Body,
					CreatedAt:  comment.CreatedAt,
					UpdatedAt:  &comment.UpdatedAt,
				}

				githubUsers[comment.User.Login] = comment.User
				commentAuthors[sourceControlComment.ID] = comment.User.Login
				commentsToSave = append(commentsToSave, sourceControlComment)
			}
		}

		// Get unique usernames
		usernames := make([]string, 0, len(githubUsers))
		for username := range githubUsers {
			usernames = append(usernames, username)
		}

		// Fetch existing source control accounts
		existingAccounts, err := p.sourceControlAPI.GetSourceControlAccountsByUsernames(ctx, usernames)
		if err != nil {
			return fmt.Errorf("failed to fetch source control accounts: %w", err)
		}

		// Create new source control accounts for missing users
		newAccounts := make([]*internaltypes.SourceControlAccount, 0)
		for username := range githubUsers {
			if _, exists := existingAccounts[username]; !exists {
				metadata, _ := json.Marshal(githubUsers[username])
				newAccounts = append(newAccounts, &internaltypes.SourceControlAccount{
					ProviderName:   "github",
					OrganizationID: &config.OrganizationID,
					Username:       username,
					Metadata:       datatypes.JSON(metadata),
				})
			}
		}

		if len(newAccounts) > 0 {
			if err := p.sourceControlAPI.CreateSourceControlAccounts(ctx, newAccounts); err != nil {
				return fmt.Errorf("failed to create source control accounts: %w", err)
			}
		}

		// Update all accounts map with newly created accounts
		for _, account := range newAccounts {
			existingAccounts[account.Username] = account
		}

		// Update pull requests with source control account IDs
		for _, pr := range prsToSave {
			var metadata map[string]interface{}
			if err := json.Unmarshal(pr.Metadata, &metadata); err != nil {
				return fmt.Errorf("failed to unmarshal PR metadata: %w", err)
			}
			user := metadata["user"].(map[string]interface{})
			author := user["login"].(string)
			if account, exists := existingAccounts[author]; exists {
				pr.SourceControlAccountID = account.ID
			}
		}

		// Update comments with source control account IDs
		for _, comment := range commentsToSave {
			authorLogin := commentAuthors[comment.ID]
			if account, exists := existingAccounts[authorLogin]; exists {
				comment.AuthorID = account.ID
			}
		}

		// Save pull requests
		if err := p.sourceControlAPI.CreatePullRequests(ctx, prsToSave); err != nil {
			return fmt.Errorf("failed to save pull requests: %w", err)
		}

		// Save comments
		if err := p.sourceControlAPI.CreatePRComments(ctx, commentsToSave); err != nil {
			return fmt.Errorf("failed to save comments: %w", err)
		}
	}

	return nil
}
