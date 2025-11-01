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
	teamapi "ems.dev/backend/services/team/api"
	teamtypes "ems.dev/backend/services/team/types"
	"github.com/samber/lo"
	"gorm.io/datatypes"
)

type GitHubProvider struct {
	githubClient     github.GithubClient
	integrationAPI   api.IntegrationAPI
	sourceControlAPI sourcecontrolapi.SourceControlAPI
	teamAPI          teamapi.TeamAPI
}

func NewProvider(
	githubClient github.GithubClient,
	integrationAPI api.IntegrationAPI,
	sourceControlAPI sourcecontrolapi.SourceControlAPI,
	teamAPI teamapi.TeamAPI,
) *GitHubProvider {
	return &GitHubProvider{
		githubClient:     githubClient,
		integrationAPI:   integrationAPI,
		sourceControlAPI: sourceControlAPI,
		teamAPI:          teamAPI,
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

	// Get all teams for the organization and extract their PR prefixes
	teams, err := p.teamAPI.ListTeams(ctx, teamtypes.TeamSearchParams{
		OrganizationID: &config.OrganizationID,
	})
	if err != nil {
		// Log error but don't fail - prefix matching is optional
		fmt.Printf("Warning: failed to fetch teams for prefix matching: %v\n", err)
	}

	// Extract prefixes from teams (filter out empty/null prefixes)
	prefixes := []string{}
	for _, team := range teams {
		if team.PRPrefix != nil && *team.PRPrefix != "" {
			prefixes = append(prefixes, *team.PRPrefix)
		}
	}

	// Helper function to find matching prefix for a PR title
	// Matches if title starts with prefix followed by a hyphen (e.g., "WL-123: Fix bug")
	findPrefixForTitle := func(title string) *string {
		upperTitle := strings.ToUpper(title)
		for _, prefix := range prefixes {
			prefixWithHyphen := strings.ToUpper(prefix) + "-"
			// Check if title starts with the prefix followed by a hyphen (case-insensitive)
			if strings.HasPrefix(upperTitle, prefixWithHyphen) {
				return &prefix
			}
		}
		return nil
	}

	// Helper function to find matching prefix in commit messages
	// Matches if commit message starts with prefix followed by a hyphen (e.g., "WL-123: Fix bug")
	findPrefixInCommits := func(commits []*githubtypes.Commit) *string {
		for _, commit := range commits {
			commitMessage := strings.TrimSpace(commit.Commit.Message)
			// Check first line of commit message (commit messages can be multi-line)
			firstLine := strings.Split(commitMessage, "\n")[0]
			upperFirstLine := strings.ToUpper(firstLine)

			for _, prefix := range prefixes {
				prefixWithHyphen := strings.ToUpper(prefix) + "-"
				// Check if commit message starts with the prefix followed by a hyphen (case-insensitive)
				if strings.HasPrefix(upperFirstLine, prefixWithHyphen) {
					return &prefix
				}
			}
		}
		return nil
	}

	// Process each repository
	for _, repo := range repositories {
		parts := strings.Split(repo, "/")
		if len(parts) != 2 {
			return fmt.Errorf("invalid repository format: %s. Expected format: owner/repo", repo)
		}

		owner, repoName := parts[0], parts[1]

		// 1. Fetch all PRs from GitHub
		// TODO: This is really inefficient, better to filter by date and try to backfill little by little
		prs, err := p.githubClient.GetPullRequests(ctx, owner, repoName, token, 20)
		if err != nil {
			return fmt.Errorf("failed to fetch pull requests for %s: %w", repo, err)
		}

		prIds := []string{}
		for _, pr := range prs {
			prIds = append(prIds, fmt.Sprintf("%d", pr.ID))
		}

		// 2. Get all PRs from the database
		// TODO: Add a filter by created at
		importedPRs, err := p.sourceControlAPI.GetPullRequests(ctx, &internaltypes.PullRequestParams{
			ProviderIDs: prIds,
		})
		if err != nil {
			return fmt.Errorf("failed to fetch imported pull requests for %s: %w", repo, err)
		}

		// Create a map of imported PRs for quick lookup
		importedPRsMap := make(map[string]internaltypes.PullRequest)
		for _, pr := range importedPRs {
			importedPRsMap[pr.ProviderID] = *pr
		}

		// 3. Process each PR
		for _, pr := range prs {
			// TODO: This should come from the integrations config. It should give the user the possibility to add settings to each repo he wants to import. (ex. Exclude PRs to prod branch)
			// Ignore PRs who's base ref is prod
			if pr.Base.Ref == "prod" {
				continue
			}

			// Check if PR exists in DB and skip if not open
			existingPR, exists := importedPRsMap[fmt.Sprintf("%d", pr.ID)]

			if exists {
				if existingPR.Status != "open" && existingPR.Status == pr.State {
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

			// Match PR title with team prefixes
			matchedPrefix := findPrefixForTitle(prDetails.Title)

			// If no prefix found in title, try to derive it from commits
			if matchedPrefix == nil {
				commits, err := p.githubClient.GetPullRequestCommits(ctx, owner, repoName, token, pr.Number)
				if err != nil {
					// Log error but don't fail - commit fetching is optional
					fmt.Printf("Warning: failed to fetch commits for PR %d: %v\n", pr.Number, err)
				} else {
					matchedPrefix = findPrefixInCommits(commits)
				}
			}

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
				Prefix:                 matchedPrefix,
				Metadata:               datatypes.JSON(prDetailsBytes),
			}

			if exists {
				sourceControlPR.ID = existingPR.ID
				err := p.sourceControlAPI.UpdatePullRequest(ctx, sourceControlPR)
				if err != nil {
					return fmt.Errorf("failed to update pull request %d: %w", pr.Number, err)
				}
			} else {
				createdPR, err := p.sourceControlAPI.CreatePullRequest(ctx, sourceControlPR)
				if err != nil {
					return fmt.Errorf("failed to create pull request %d: %w", pr.Number, err)
				}
				// Update the sourceControlPR with the created PR's ID if needed
				sourceControlPR.ID = createdPR.ID
			}

			// 6. Get all reviews, comments and review comments
			existingComments, err := p.sourceControlAPI.GetPullRequestComments(ctx, sourceControlPR.ID)
			if err != nil {
				return fmt.Errorf("failed to fetch comments for PR %d: %w", pr.Number, err)
			}

			existingCommentsMap := make(map[string]internaltypes.PRComment)
			for _, comment := range existingComments {
				existingCommentsMap[comment.ProviderID] = *comment
			}

			reviews, err := p.githubClient.GetPullRequestReviews(ctx, owner, repoName, token, pr.Number)
			if err != nil {
				return fmt.Errorf("failed to fetch reviews for PR %d: %w", pr.Number, err)
			}

			allComments := []*githubtypes.ReviewComment{}
			for _, review := range reviews {
				if _, exists := existingCommentsMap[fmt.Sprintf("%d", review.ID)]; exists {
					continue
				}

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
				if _, exists := existingCommentsMap[fmt.Sprintf("%d", reviewComment.ID)]; exists {
					continue
				}

				reviewComment.Type = string(githubtypes.CommentTypeReviewComment)
				allComments = append(allComments, reviewComment)
			}

			comments, err := p.githubClient.GetPullRequestComments(ctx, owner, repoName, token, pr.Number)
			if err != nil {
				return fmt.Errorf("failed to fetch regular comments for PR %d: %w", pr.Number, err)
			}

			for _, comment := range comments {
				if _, exists := existingCommentsMap[fmt.Sprintf("%d", comment.ID)]; exists {
					continue
				}

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
	accounts, err := p.sourceControlAPI.GetSourceControlAccounts(ctx, &internaltypes.SourceControlAccountParams{
		Usernames: []string{user.Login},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch source control account: %w", err)
	}

	accountsMap := make(map[string]internaltypes.SourceControlAccount)
	for _, account := range accounts {
		accountsMap[account.Username] = account
	}

	if account, exists := accountsMap[user.Login]; exists {
		if *account.OrganizationID == organizationID {
			return &account, nil
		}
	}

	// Get the member ID for this user in this org username
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
