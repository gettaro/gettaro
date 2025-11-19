package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"ems.dev/backend/services/title/types"
	"github.com/stretchr/testify/assert"
)

func TestGetTitle(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockTitle     *types.Title
		mockError     error
		expectedTitle *types.Title
		expectedError error
	}{
		{
			name: "successful retrieval",
			id:   "title-1",
			mockTitle: &types.Title{
				ID:             "title-1",
				Name:           "Software Engineer",
				OrganizationID: "org-1",
				IsManager:      false,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			expectedTitle: &types.Title{
				ID:             "title-1",
				Name:           "Software Engineer",
				OrganizationID: "org-1",
				IsManager:      false,
			},
		},
		{
			name:          "title not found",
			id:            "non-existent",
			mockTitle:     nil,
			expectedTitle: nil,
		},
		{
			name:        "database error",
			id:          "title-1",
			mockError:   errors.New("database error"),
			expectedError: errors.New("database error"),
		},
		{
			name: "invalid id",
			id:   "",
			mockTitle: nil,
			expectedTitle: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockTitleDB)
			api := &Api{db: mockDB}

			mockDB.On("GetTitle", tt.id).Return(tt.mockTitle, tt.mockError)

			result, err := api.GetTitle(context.Background(), tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if tt.expectedTitle == nil {
					assert.Nil(t, result)
				} else {
					assert.NotNil(t, result)
					assert.Equal(t, tt.expectedTitle.ID, result.ID)
					assert.Equal(t, tt.expectedTitle.Name, result.Name)
					assert.Equal(t, tt.expectedTitle.OrganizationID, result.OrganizationID)
					assert.Equal(t, tt.expectedTitle.IsManager, result.IsManager)
				}
			}

			mockDB.AssertExpectations(t)
		})
	}
}
