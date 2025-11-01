package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"ems.dev/backend/libraries/github/types"
)

type GithubClient interface {
	GetPullRequests(ctx context.Context, owner, repo, token string, maxPages int) ([]*types.PullRequest, error)
	GetPullRequestReviewComments(ctx context.Context, owner, repo, token string, prNumber int) ([]*types.ReviewComment, error)
	GetPullRequest(ctx context.Context, owner, repo, token string, prNumber int) (*types.PullRequest, error)
	GetPullRequestComments(ctx context.Context, owner, repo, token string, prNumber int) ([]*types.ReviewComment, error)
	GetPullRequestReviews(ctx context.Context, owner, repo, token string, prNumber int) ([]*types.Review, error)
	GetPullRequestCommits(ctx context.Context, owner, repo, token string, prNumber int) ([]*types.Commit, error)
}

// Client represents a GitHub API client
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new GitHub client
func NewClient() *Client {
	return &Client{
		baseURL:    "https://api.github.com",
		httpClient: &http.Client{},
	}
}

// GetPullRequests fetches pull requests for a repository with pagination support
func (c *Client) GetPullRequests(ctx context.Context, owner, repo, token string, maxPages int) ([]*types.PullRequest, error) {
	var allPRs []*types.PullRequest
	page := 1
	url := fmt.Sprintf("%s/repos/%s/%s/pulls?state=all&per_page=100", c.baseURL, owner, repo)

	for page <= maxPages {
		pageURL := fmt.Sprintf("%s&page=%d", url, page)
		req, err := http.NewRequestWithContext(ctx, "GET", pageURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Accept", "application/vnd.github.v3+json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		var prs []*types.PullRequest
		if err := json.NewDecoder(resp.Body).Decode(&prs); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		// If no PRs were returned, we've reached the end
		if len(prs) == 0 {
			break
		}

		allPRs = append(allPRs, prs...)
		page++
	}

	return allPRs, nil
}

// GetPullRequestReviewComments fetches review comments for a specific pull request
func (c *Client) GetPullRequestReviewComments(ctx context.Context, owner, repo, token string, prNumber int) ([]*types.ReviewComment, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%d/comments", c.baseURL, owner, repo, prNumber)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var comments []*types.ReviewComment
	if err := json.NewDecoder(resp.Body).Decode(&comments); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return comments, nil
}

// GetPullRequest fetches details of a single pull request
func (c *Client) GetPullRequest(ctx context.Context, owner, repo, token string, prNumber int) (*types.PullRequest, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%d", c.baseURL, owner, repo, prNumber)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var pr types.PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &pr, nil
}

// GetPullRequestComments fetches regular comments for a specific pull request
func (c *Client) GetPullRequestComments(ctx context.Context, owner, repo, token string, prNumber int) ([]*types.ReviewComment, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/issues/%d/comments", c.baseURL, owner, repo, prNumber)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var comments []*types.ReviewComment
	if err := json.NewDecoder(resp.Body).Decode(&comments); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return comments, nil
}

// GetPullRequestReviewComments fetches review comments for a specific pull request
func (c *Client) GetPullRequestReviews(ctx context.Context, owner, repo, token string, prNumber int) ([]*types.Review, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%d/reviews", c.baseURL, owner, repo, prNumber)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var reviews []*types.Review
	if err := json.NewDecoder(resp.Body).Decode(&reviews); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return reviews, nil
}

// GetPullRequestCommits fetches commits for a specific pull request
func (c *Client) GetPullRequestCommits(ctx context.Context, owner, repo, token string, prNumber int) ([]*types.Commit, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%d/commits", c.baseURL, owner, repo, prNumber)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var commits []*types.Commit
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return commits, nil
}
