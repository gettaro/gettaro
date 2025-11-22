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

func TestGetMemberPullRequestReviews(t *testing.T) {
	now := time.Now()
	startDate := now.AddDate(0, -1, 0)
	endDate := now

	tests := []struct {
		name            string
		params          *types.MemberPullRequestReviewsParams
		mockActivities  []*types.MemberActivity
		mockError       error
		expectedActivities []*types.MemberActivity
		expectedError   error
	}{
		{
			name: "success - returns member pull request reviews",
			params: &types.MemberPullRequestReviewsParams{
				MemberID: "member-1",
			},
			mockActivities: []*types.MemberActivity{
				{
					ID:        "activity-1",
					Type:      "pr_review",
					Title:     "Reviewed PR",
					CreatedAt: now,
				},
			},
			expectedActivities: []*types.MemberActivity{
				{
					ID:        "activity-1",
					Type:      "pr_review",
					Title:     "Reviewed PR",
					CreatedAt: now,
				},
			},
		},
		{
			name: "success - with date range filter",
			params: &types.MemberPullRequestReviewsParams{
				MemberID:  "member-1",
				StartDate: &startDate,
				EndDate:   &endDate,
			},
			mockActivities: []*types.MemberActivity{
				{
					ID:        "activity-1",
					Type:      "pr_review",
					Title:     "Reviewed PR",
					CreatedAt: now,
				},
			},
			expectedActivities: []*types.MemberActivity{
				{
					ID:        "activity-1",
					Type:      "pr_review",
					Title:     "Reviewed PR",
					CreatedAt: now,
				},
			},
		},
		{
			name: "success - with has body filter",
			params: &types.MemberPullRequestReviewsParams{
				MemberID: "member-1",
				HasBody:  boolPtr(true),
			},
			mockActivities: []*types.MemberActivity{
				{
					ID:        "activity-1",
					Type:      "pr_review",
					Title:     "Reviewed PR",
					Description: "This is a detailed review",
					CreatedAt: now,
				},
			},
			expectedActivities: []*types.MemberActivity{
				{
					ID:        "activity-1",
					Type:      "pr_review",
					Title:     "Reviewed PR",
					Description: "This is a detailed review",
					CreatedAt: now,
				},
			},
		},
		{
			name: "success - empty result",
			params: &types.MemberPullRequestReviewsParams{
				MemberID: "member-1",
			},
			mockActivities:     []*types.MemberActivity{},
			expectedActivities: []*types.MemberActivity{},
		},
		{
			name: "error - database error",
			params: &types.MemberPullRequestReviewsParams{
				MemberID: "member-1",
			},
			mockError:     errors.New("database query failed"),
			expectedError: errors.New("database query failed"),
		},
		{
			name: "error - invalid member ID",
			params: &types.MemberPullRequestReviewsParams{
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

			mockDB.On("GetMemberPullRequestReviews", mock.Anything, tt.params).Return(tt.mockActivities, tt.mockError)

			activities, err := api.GetMemberPullRequestReviews(context.Background(), tt.params)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, activities)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedActivities), len(activities))
				if len(activities) > 0 {
					assert.Equal(t, tt.expectedActivities[0].ID, activities[0].ID)
					assert.Equal(t, tt.expectedActivities[0].Type, activities[0].Type)
				}
			}

			mockDB.AssertExpectations(t)
		})
	}
}
