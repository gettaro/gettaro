package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"ems.dev/backend/services/title/types"
	"github.com/stretchr/testify/assert"
)

func TestUpdateTitle(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		request       types.Title
		mockTitle     *types.Title
		getTitleError error
		updateError   error
		expectedTitle *types.Title
		expectedError error
	}{
		{
			name: "successful update",
			id:   "title-1",
			request: types.Title{
				Name:      "Senior Software Engineer",
				IsManager: true,
			},
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
				Name:           "Senior Software Engineer",
				OrganizationID: "org-1",
				IsManager:      true,
			},
		},
		{
			name: "title not found",
			id:   "non-existent",
			request: types.Title{
				Name:      "Senior Software Engineer",
				IsManager: true,
			},
			mockTitle:     nil,
			expectedTitle: nil,
		},
		{
			name: "get title database error",
			id:   "title-1",
			request: types.Title{
				Name:      "Senior Software Engineer",
				IsManager: true,
			},
			getTitleError: errors.New("database error"),
			expectedError: errors.New("database error"),
		},
		{
			name: "update database error",
			id:   "title-1",
			request: types.Title{
				Name:      "Senior Software Engineer",
				IsManager: true,
			},
			mockTitle: &types.Title{
				ID:             "title-1",
				Name:           "Software Engineer",
				OrganizationID: "org-1",
				IsManager:      false,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			updateError:   errors.New("update failed"),
			expectedError: errors.New("update failed"),
		},
		{
			name: "invalid id",
			id:   "",
			request: types.Title{
				Name:      "Senior Software Engineer",
				IsManager: true,
			},
			mockTitle:     nil,
			expectedTitle: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockTitleDB)
			api := &Api{db: mockDB}

			mockDB.On("GetTitle", tt.id).Return(tt.mockTitle, tt.getTitleError)
			if tt.mockTitle != nil && tt.getTitleError == nil {
				updatedTitle := *tt.mockTitle
				updatedTitle.Name = tt.request.Name
				updatedTitle.IsManager = tt.request.IsManager
				mockDB.On("UpdateTitle", updatedTitle).Return(tt.updateError)
			}

			result, err := api.UpdateTitle(context.Background(), tt.id, tt.request)

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
