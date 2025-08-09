package github

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"ems.dev/backend/libraries/github"
	githubtypes "ems.dev/backend/libraries/github/types"
	"ems.dev/backend/services/integration/api"
	"ems.dev/backend/services/integration/types"
	sourcecontrolapi "ems.dev/backend/services/sourcecontrol/api"
	internaltypes "ems.dev/backend/services/sourcecontrol/types"
	"github.com/samber/lo"
	"gorm.io/datatypes"
)

type GitHubProvider struct {
	githubClient     github.GithubClient
	integrationAPI   api.IntegrationAPI
	sourceControlAPI sourcecontrolapi.SourceControlAPI
}

func NewProvider(
	githubClient github.GithubClient,
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

// TOOD: Need to extract this logic to sync.go. This should be provide agnostic and instead is using the github client directly.
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
		// TODO: Add a filter by created at
		importedPRs, err := p.sourceControlAPI.GetPullRequests(ctx, &internaltypes.PullRequestParams{
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
		prs, err := p.githubClient.GetPullRequests(ctx, owner, repoName, token, 20)
		if err != nil {
			return fmt.Errorf("failed to fetch pull requests for %s: %w", repo, err)
		}

		// 3. Process each PR
		for _, pr := range prs {
			// TODO: This should come from the integrations config. It should give the user the possibility to add settings to each repo he wants to import. (ex. Exclude PRs to prod branch)
			// Ignore PRs who's base ref is prod
			if pr.Base.Ref == "prod" {
				continue
			}

			// Check if PR exists in DB and skip if not open
			// TODO: If PR is closed, we should flag it as such in our db.
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

			// 6. Get all reviews, comments and review comments
			reviews, err := p.githubClient.GetPullRequestReviews(ctx, owner, repoName, token, pr.Number)
			if err != nil {
				return fmt.Errorf("failed to fetch reviews for PR %d: %w", pr.Number, err)
			}

			allComments := []*githubtypes.ReviewComment{}
			for _, review := range reviews {
				allComments = append(allComments, &githubtypes.ReviewComment{
					ID:        review.ID,
					User:      review.User,
					Body:      review.Body,
					CreatedAt: review.SubmittedAt,
					Type:      string(githubtypes.CommentTypeReview),
				})
			}

			reviewComments, err := p.githubClient.GetPullRequestReviewComments(ctx, owner, repoName, token, pr.Number)
			if err != nil {
				return fmt.Errorf("failed to fetch review comments for PR %d: %w", pr.Number, err)
			}

			for _, reviewComment := range reviewComments {
				reviewComment.Type = string(githubtypes.CommentTypeReviewComment)
				allComments = append(allComments, reviewComment)
			}

			comments, err := p.githubClient.GetPullRequestComments(ctx, owner, repoName, token, pr.Number)
			if err != nil {
				return fmt.Errorf("failed to fetch regular comments for PR %d: %w", pr.Number, err)
			}

			for _, comment := range comments {
				comment.Type = string(githubtypes.CommentTypeComment)
				allComments = append(allComments, comment)
			}

			// 7. Process each comment
			for _, comment := range allComments {
				// Insert/Update comment author
				commentAuthor, err := p.upsertAuthor(ctx, config.OrganizationID, comment.User)
				if err != nil {
					return fmt.Errorf("failed to upsert comment author for PR %d: %w", pr.Number, err)
				}

				// Insert/Update comment
				sourceControlComment := &internaltypes.PRComment{
					PRID:                   sourceControlPR.ID,
					SourceControlAccountID: commentAuthor.ID,
					ProviderID:             fmt.Sprintf("%d", comment.ID),
					Body:                   comment.Body,
					Type:                   comment.Type,
					CreatedAt:              comment.CreatedAt,
					UpdatedAt:              &comment.UpdatedAt,
				}

				if err := p.sourceControlAPI.CreatePRComments(ctx, []*internaltypes.PRComment{sourceControlComment}); err != nil {
					return fmt.Errorf("failed to save comment for PR %d: %w", pr.Number, err)
				}
			}

			// 9. Calculate PR metrics
			metrics := make(map[string]interface{})

			// Filter out comments that are from bots
			nonBotComments := lo.Filter(allComments, func(comment *githubtypes.ReviewComment, _ int) bool {
				return comment.User.Type != "Bot"
			})

			metrics["number_of_non_bot_comments"] = len(nonBotComments)

			// Calculate time to merge if the PR is merged
			if prDetails.MergedAt != nil {
				timeToMerge := prDetails.MergedAt.Sub(prDetails.CreatedAt)
				metrics["time_to_merge_seconds"] = int64(timeToMerge.Seconds())
			}

			// Calculate time to first review
			var firstReviewedAt time.Time
			if len(reviews) > 0 {
				firstReviewedAt = reviews[0].SubmittedAt
				metrics["time_to_first_non_bot_review_seconds"] = int64(firstReviewedAt.Sub(prDetails.CreatedAt).Seconds())
			}

			if len(nonBotComments) > 0 {
				// Sort comments by creation time to find the first one
				sort.Slice(nonBotComments, func(i, j int) bool {
					return nonBotComments[i].CreatedAt.Before(nonBotComments[j].CreatedAt)
				})

				firstComment := nonBotComments[0]
				if firstComment.CreatedAt.Before(firstReviewedAt) {
					firstReviewedAt = firstComment.CreatedAt
				}

				timeToFirstReview := firstReviewedAt.Sub(prDetails.CreatedAt)
				metrics["time_to_first_non_bot_review_seconds"] = int64(timeToFirstReview.Seconds())
			}

			// Update PR with metrics
			metricsBytes, _ := json.Marshal(metrics)
			sourceControlPR.Metrics = datatypes.JSON(metricsBytes)

			// Save updated PR with metrics
			if err := p.sourceControlAPI.UpdatePullRequest(ctx, sourceControlPR); err != nil {
				return fmt.Errorf("failed to save pull request metrics for PR %d: %w", pr.Number, err)
			}
		}
	}

	return nil
}

// upsertAuthor handles the creation or update of a source control account
func (p *GitHubProvider) upsertAuthor(ctx context.Context, organizationID string, user githubtypes.User) (*internaltypes.SourceControlAccount, error) {
	// Check if account exists
	// TODO: Get account by username
	accounts, err := p.sourceControlAPI.GetSourceControlAccountsByUsernames(ctx, []string{user.Login})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch source control account: %w", err)
	}

	if account, exists := accounts[user.Login]; exists {
		if *account.OrganizationID == organizationID {
			return account, nil
		}
	}

	// For now, we'll create the account without member_id and let the application handle the association later
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
