package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"ems.dev/backend/services/sourcecontrol/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetPullRequests(t *testing.T) {
	tests := []struct {
		name           string
		params         *types.PullRequestParams
		mockPRs        []*types.PullRequest
		mockError      error
		expectedPRs    []*types.PullRequest
		expectedError  error
	}{
		{
			name: "success - returns pull requests",
			params: &types.PullRequestParams{
				RepositoryName: "test-repo",
			},
			mockPRs: []*types.PullRequest{
				{
					ID:             "pr-1",
					RepositoryName: "test-repo",
					Title:          "Test PR",
					Status:         "open",
				},
			},
			expectedPRs: []*types.PullRequest{
				{
					ID:             "pr-1",
					RepositoryName: "test-repo",
					Title:          "Test PR",
					Status:         "open",
				},
			},
		},
		{
			name: "success - empty result",
			params: &types.PullRequestParams{
				RepositoryName: "empty-repo",
			},
			mockPRs:     []*types.PullRequest{},
			expectedPRs: []*types.PullRequest{},
		},
		{
			name: "error - database error",
			params: &types.PullRequestParams{
				RepositoryName: "test-repo",
			},
			mockError:     errors.New("database connection failed"),
			expectedError: errors.New("database connection failed"),
		},
		{
			name: "success - with date range filter",
			params: &types.PullRequestParams{
				RepositoryName: "test-repo",
				StartDate:      timePtr(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				EndDate:        timePtr(time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)),
			},
			mockPRs: []*types.PullRequest{
				{
					ID:             "pr-1",
					RepositoryName: "test-repo",
					Title:          "Test PR",
					Status:         "open",
				},
			},
			expectedPRs: []*types.PullRequest{
				{
					ID:             "pr-1",
					RepositoryName: "test-repo",
					Title:          "Test PR",
					Status:         "open",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			api := &Api{
				db:            mockDB,
				metricsEngine: nil,
			}

			mockDB.On("GetPullRequests", mock.Anything, tt.params).Return(tt.mockPRs, tt.mockError)

			prs, err := api.GetPullRequests(context.Background(), tt.params)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, prs)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPRs, prs)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

// Helper function to create time pointers
func timePtr(t time.Time) *time.Time {
	return &t
}
