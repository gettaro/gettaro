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

func TestGetPullRequestComments(t *testing.T) {
	tests := []struct {
		name            string
		prID            string
		mockComments    []*types.PRComment
		mockError       error
		expectedComments []*types.PRComment
		expectedError   error
	}{
		{
			name: "success - returns comments",
			prID: "pr-1",
			mockComments: []*types.PRComment{
				{
					ID:        "comment-1",
					PRID:      "pr-1",
					Body:      "Great work!",
					Type:      "COMMENT",
					CreatedAt: time.Now(),
				},
				{
					ID:        "comment-2",
					PRID:      "pr-1",
					Body:      "Looks good",
					Type:      "REVIEW",
					CreatedAt: time.Now(),
				},
			},
			expectedComments: []*types.PRComment{
				{
					ID:        "comment-1",
					PRID:      "pr-1",
					Body:      "Great work!",
					Type:      "COMMENT",
					CreatedAt: time.Now(),
				},
				{
					ID:        "comment-2",
					PRID:      "pr-1",
					Body:      "Looks good",
					Type:      "REVIEW",
					CreatedAt: time.Now(),
				},
			},
		},
		{
			name:            "success - empty result",
			prID:            "pr-1",
			mockComments:    []*types.PRComment{},
			expectedComments: []*types.PRComment{},
		},
		{
			name:          "error - database error",
			prID:          "pr-1",
			mockError:     errors.New("database query failed"),
			expectedError: errors.New("database query failed"),
		},
		{
			name:          "error - invalid PR ID",
			prID:          "",
			mockError:     errors.New("PR ID is required"),
			expectedError: errors.New("PR ID is required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			api := &Api{
				db:            mockDB,
				metricsEngine: nil,
			}

			mockDB.On("GetPullRequestComments", mock.Anything, tt.prID).Return(tt.mockComments, tt.mockError)

			comments, err := api.GetPullRequestComments(context.Background(), tt.prID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, comments)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedComments), len(comments))
				if len(comments) > 0 {
					assert.Equal(t, tt.expectedComments[0].ID, comments[0].ID)
					assert.Equal(t, tt.expectedComments[0].PRID, comments[0].PRID)
				}
			}

			mockDB.AssertExpectations(t)
		})
	}
}
