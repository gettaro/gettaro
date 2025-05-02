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

	// Process each repository
	for _, repo := range repositories {
		parts := strings.Split(repo, "/")
		if len(parts) != 2 {
			return fmt.Errorf("invalid repository format: %s. Expected format: owner/repo", repo)
		}
		owner, repoName := parts[0], parts[1]

		// 1. Get all PRs from the database
		importedPRs, err := p.sourceControlAPI.GetPullRequests(ctx, &internaltypes.PullRequestParams{
			ProviderID:     fmt.Sprintf("%s/%s", owner, repoName),
			OrganizationID: &config.OrganizationID,
			RepositoryName: repoName,
		})
		if err != nil {
			return fmt.Errorf("failed to fetch imported pull requests for %s: %w", repo, err)
		}

		// Create a map of imported PRs for quick lookup
		importedPRsMap := make(map[string]internaltypes.PullRequest)
		for _, pr := range importedPRs {
			importedPRsMap[pr.ProviderID] = *pr
		}

		// 2. Fetch all PRs from GitHub
		prs, err := p.githubClient.GetPullRequests(ctx, owner, repoName, token)
		if err != nil {
			return fmt.Errorf("failed to fetch pull requests for %s: %w", repo, err)
		}

		// Process each PR
		for _, pr := range prs {
			// 3. Check if PR exists in DB and skip if not open
			if existingPR, exists := importedPRsMap[fmt.Sprintf("%d", pr.ID)]; exists {
				if existingPR.Status != "open" {
					continue
				}
			}

			// Get detailed PR information
			prDetails, err := p.githubClient.GetPullRequest(ctx, owner, repoName, token, pr.Number)
			if err != nil {
				return fmt.Errorf("failed to fetch pull request details for PR %d: %w", pr.Number, err)
			}

			// 4. Insert/Update author
			authorAccount, err := p.upsertAuthor(ctx, config.OrganizationID, prDetails.User)
			if err != nil {
				return fmt.Errorf("failed to upsert author for PR %d: %w", pr.Number, err)
			}

			// 5. Insert/Update PR
			prDetailsBytes, _ := json.Marshal(prDetails)
			sourceControlPR := &internaltypes.PullRequest{
				SourceControlAccountID: authorAccount.ID,
				ProviderID:             fmt.Sprintf("%d", prDetails.ID),
				RepositoryName:         repoName,
				OrganizationID:         config.OrganizationID,
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
				URL:                    prDetails.URL,
				Metadata:               datatypes.JSON(prDetailsBytes),
			}

			if err := p.sourceControlAPI.CreatePullRequests(ctx, []*internaltypes.PullRequest{sourceControlPR}); err != nil {
				return fmt.Errorf("failed to save pull request %d: %w", pr.Number, err)
			}

			// 6. Get all comments and review comments
			reviewComments, err := p.githubClient.GetPullRequestReviewComments(ctx, owner, repoName, token, pr.Number)
			if err != nil {
				return fmt.Errorf("failed to fetch review comments for PR %d: %w", pr.Number, err)
			}

			comments, err := p.githubClient.GetPullRequestComments(ctx, owner, repoName, token, pr.Number)
			if err != nil {
				return fmt.Errorf("failed to fetch regular comments for PR %d: %w", pr.Number, err)
			}

			allComments := append(reviewComments, comments...)

			// Process each comment
			for _, comment := range allComments {
				// 7. Insert/Update comment author
				commentAuthor, err := p.upsertAuthor(ctx, config.OrganizationID, comment.User)
				if err != nil {
					return fmt.Errorf("failed to upsert comment author for PR %d: %w", pr.Number, err)
				}

				// 8. Insert/Update comment
				sourceControlComment := &internaltypes.PRComment{
					PRID:       sourceControlPR.ID,
					AuthorID:   commentAuthor.ID,
					ProviderID: fmt.Sprintf("%d", comment.ID),
					Body:       comment.Body,
					CreatedAt:  comment.CreatedAt,
					UpdatedAt:  &comment.UpdatedAt,
				}

				if err := p.sourceControlAPI.CreatePRComments(ctx, []*internaltypes.PRComment{sourceControlComment}); err != nil {
					return fmt.Errorf("failed to save comment for PR %d: %w", pr.Number, err)
				}
			}
		}
	}

	return nil
}

// upsertAuthor handles the creation or update of a source control account
func (p *GitHubProvider) upsertAuthor(ctx context.Context, organizationID string, user githubtypes.User) (*internaltypes.SourceControlAccount, error) {
	// Check if account exists
	accounts, err := p.sourceControlAPI.GetSourceControlAccountsByUsernames(ctx, []string{user.Login})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch source control account: %w", err)
	}

	if account, exists := accounts[user.Login]; exists {
		return account, nil
	}

	// Create new account
	metadata, _ := json.Marshal(user)
	newAccount := &internaltypes.SourceControlAccount{
		ProviderName:   "github",
		OrganizationID: &organizationID,
		Username:       user.Login,
		ProviderID:     fmt.Sprintf("%d", user.ID),
		Metadata:       datatypes.JSON(metadata),
	}

	if err := p.sourceControlAPI.CreateSourceControlAccounts(ctx, []*internaltypes.SourceControlAccount{newAccount}); err != nil {
		return nil, fmt.Errorf("failed to create source control account: %w", err)
	}

	return newAccount, nil
}
