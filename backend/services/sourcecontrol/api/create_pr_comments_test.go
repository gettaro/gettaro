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

func TestCreatePRComments(t *testing.T) {
	tests := []struct {
		name          string
		comments      []*types.PRComment
		mockError     error
		expectedError error
	}{
		{
			name: "success - creates comments",
			comments: []*types.PRComment{
				{
					PRID:     "pr-1",
					Body:     "Great work!",
					Type:     "COMMENT",
					CreatedAt: time.Now(),
				},
				{
					PRID:     "pr-1",
					Body:     "Looks good to me",
					Type:     "REVIEW",
					CreatedAt: time.Now(),
				},
			},
		},
		{
			name: "success - empty comments array",
			comments: []*types.PRComment{},
		},
		{
			name: "error - database error",
			comments: []*types.PRComment{
				{
					PRID: "pr-1",
					Body: "Comment",
				},
			},
			mockError:     errors.New("database insert failed"),
			expectedError: errors.New("database insert failed"),
		},
		{
			name: "error - invalid PR ID",
			comments: []*types.PRComment{
				{
					PRID: "",
					Body: "Comment",
				},
			},
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

			mockDB.On("CreatePRComments", mock.Anything, tt.comments).Return(tt.mockError)

			err := api.CreatePRComments(context.Background(), tt.comments)

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
