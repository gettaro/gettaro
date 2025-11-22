package api

import (
	"context"
	"errors"
	"testing"

	"ems.dev/backend/services/sourcecontrol/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdatePullRequest(t *testing.T) {
	tests := []struct {
		name          string
		pr            *types.PullRequest
		mockError     error
		expectedError error
	}{
		{
			name: "success - updates pull request",
			pr: &types.PullRequest{
				ID:     "pr-1",
				Status: "closed",
			},
		},
		{
			name: "error - database error",
			pr: &types.PullRequest{
				ID:     "pr-1",
				Status: "closed",
			},
			mockError:     errors.New("database update failed"),
			expectedError: errors.New("database update failed"),
		},
		{
			name: "error - pull request not found",
			pr: &types.PullRequest{
				ID:     "non-existent",
				Status: "closed",
			},
			mockError:     errors.New("pull request not found"),
			expectedError: errors.New("pull request not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			api := &Api{
				db:            mockDB,
				metricsEngine: nil,
			}

			mockDB.On("UpdatePullRequest", mock.Anything, tt.pr).Return(tt.mockError)

			err := api.UpdatePullRequest(context.Background(), tt.pr)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
		})
	}
}
