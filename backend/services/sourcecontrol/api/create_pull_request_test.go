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

func TestCreatePullRequest(t *testing.T) {
	tests := []struct {
		name          string
		pr            *types.PullRequest
		mockPR        *types.PullRequest
		mockError     error
		expectedPR    *types.PullRequest
		expectedError error
	}{
		{
			name: "success - creates pull request",
			pr: &types.PullRequest{
				RepositoryName: "test-repo",
				Title:          "New PR",
				Status:         "open",
			},
			mockPR: &types.PullRequest{
				ID:             "pr-1",
				RepositoryName: "test-repo",
				Title:          "New PR",
				Status:         "open",
				CreatedAt:      time.Now(),
			},
			expectedPR: &types.PullRequest{
				ID:             "pr-1",
				RepositoryName: "test-repo",
				Title:          "New PR",
				Status:         "open",
				CreatedAt:      time.Now(),
			},
		},
		{
			name: "error - database error",
			pr: &types.PullRequest{
				RepositoryName: "test-repo",
				Title:          "New PR",
			},
			mockError:     errors.New("database insert failed"),
			expectedError: errors.New("database insert failed"),
		},
		{
			name: "error - validation error",
			pr: &types.PullRequest{
				RepositoryName: "",
				Title:          "New PR",
			},
			mockError:     errors.New("repository name is required"),
			expectedError: errors.New("repository name is required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			api := &Api{
				db:            mockDB,
				metricsEngine: nil,
			}

			mockDB.On("CreatePullRequest", mock.Anything, tt.pr).Return(tt.mockPR, tt.mockError)

			pr, err := api.CreatePullRequest(context.Background(), tt.pr)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, pr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPR.ID, pr.ID)
				assert.Equal(t, tt.expectedPR.RepositoryName, pr.RepositoryName)
				assert.Equal(t, tt.expectedPR.Title, pr.Title)
			}

			mockDB.AssertExpectations(t)
		})
	}
}
