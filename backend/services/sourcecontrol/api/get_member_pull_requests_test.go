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

func TestGetMemberPullRequests(t *testing.T) {
	now := time.Now()
	startDate := now.AddDate(0, -1, 0)
	endDate := now

	tests := []struct {
		name           string
		params         *types.MemberPullRequestParams
		mockPRs        []*types.PullRequestWithComments
		mockError      error
		expectedPRs    []*types.PullRequestWithComments
		expectedError  error
	}{
		{
			name: "success - returns member pull requests",
			params: &types.MemberPullRequestParams{
				MemberID: "member-1",
			},
			mockPRs: []*types.PullRequestWithComments{
				{
					PullRequest: &types.PullRequest{
						ID:             "pr-1",
						RepositoryName: "test-repo",
						Title:          "Test PR",
						Status:         "open",
					},
					Comments: []*types.PRComment{},
				},
			},
			expectedPRs: []*types.PullRequestWithComments{
				{
					PullRequest: &types.PullRequest{
						ID:             "pr-1",
						RepositoryName: "test-repo",
						Title:          "Test PR",
						Status:         "open",
					},
					Comments: []*types.PRComment{},
				},
			},
		},
		{
			name: "success - with date range filter",
			params: &types.MemberPullRequestParams{
				MemberID:  "member-1",
				StartDate: &startDate,
				EndDate:   &endDate,
			},
			mockPRs: []*types.PullRequestWithComments{
				{
					PullRequest: &types.PullRequest{
						ID:             "pr-1",
						RepositoryName: "test-repo",
						Title:          "Test PR",
						Status:         "open",
					},
					Comments: []*types.PRComment{},
				},
			},
			expectedPRs: []*types.PullRequestWithComments{
				{
					PullRequest: &types.PullRequest{
						ID:             "pr-1",
						RepositoryName: "test-repo",
						Title:          "Test PR",
						Status:         "open",
					},
					Comments: []*types.PRComment{},
				},
			},
		},
		{
			name: "success - with comments included",
			params: &types.MemberPullRequestParams{
				MemberID:        "member-1",
				IncludeComments: boolPtr(true),
			},
			mockPRs: []*types.PullRequestWithComments{
				{
					PullRequest: &types.PullRequest{
						ID:             "pr-1",
						RepositoryName: "test-repo",
						Title:          "Test PR",
						Status:         "open",
					},
					Comments: []*types.PRComment{
						{
							ID:   "comment-1",
							PRID: "pr-1",
							Body: "Great work!",
						},
					},
				},
			},
			expectedPRs: []*types.PullRequestWithComments{
				{
					PullRequest: &types.PullRequest{
						ID:             "pr-1",
						RepositoryName: "test-repo",
						Title:          "Test PR",
						Status:         "open",
					},
					Comments: []*types.PRComment{
						{
							ID:   "comment-1",
							PRID: "pr-1",
							Body: "Great work!",
						},
					},
				},
			},
		},
		{
			name: "success - empty result",
			params: &types.MemberPullRequestParams{
				MemberID: "member-1",
			},
			mockPRs:     []*types.PullRequestWithComments{},
			expectedPRs: []*types.PullRequestWithComments{},
		},
		{
			name: "error - database error",
			params: &types.MemberPullRequestParams{
				MemberID: "member-1",
			},
			mockError:     errors.New("database query failed"),
			expectedError: errors.New("database query failed"),
		},
		{
			name: "error - invalid member ID",
			params: &types.MemberPullRequestParams{
				MemberID: "",
			},
			mockError:     errors.New("member ID is required"),
			expectedError: errors.New("member ID is required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			api := &Api{
				db:            mockDB,
				metricsEngine: nil,
			}

			mockDB.On("GetMemberPullRequests", mock.Anything, tt.params).Return(tt.mockPRs, tt.mockError)

			prs, err := api.GetMemberPullRequests(context.Background(), tt.params)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, prs)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedPRs), len(prs))
				if len(prs) > 0 {
					assert.Equal(t, tt.expectedPRs[0].PullRequest.ID, prs[0].PullRequest.ID)
				}
			}

			mockDB.AssertExpectations(t)
		})
	}
}

// Helper function to create bool pointers
func boolPtr(b bool) *bool {
	return &b
}
