package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"ems.dev/backend/services/title/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTitle(t *testing.T) {
	tests := []struct {
		name          string
		title         types.Title
		mockError     error
		expectedTitle *types.Title
		expectedError error
	}{
		{
			name: "successful creation",
			title: types.Title{
				Name:           "Software Engineer",
				OrganizationID: "org-1",
				IsManager:      false,
			},
			expectedTitle: &types.Title{
				Name:           "Software Engineer",
				OrganizationID: "org-1",
				IsManager:      false,
			},
		},
		{
			name: "database error",
			title: types.Title{
				Name:           "Software Engineer",
				OrganizationID: "org-1",
				IsManager:      false,
			},
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
		{
			name: "duplicate title error",
			title: types.Title{
				Name:           "Software Engineer",
				OrganizationID: "org-1",
				IsManager:      false,
			},
			mockError:     errors.New("duplicate title"),
			expectedError: errors.New("duplicate title"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockTitleDB)
			api := &Api{db: mockDB}

			mockDB.On("CreateTitle", mock.AnythingOfType("*types.Title")).Return(tt.mockError).Run(func(args mock.Arguments) {
				title := args.Get(0).(*types.Title)
				title.ID = "title-1"
				title.CreatedAt = time.Now()
				title.UpdatedAt = time.Now()
			})

			result, err := api.CreateTitle(context.Background(), tt.title)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedTitle.Name, result.Name)
				assert.Equal(t, tt.expectedTitle.OrganizationID, result.OrganizationID)
				assert.Equal(t, tt.expectedTitle.IsManager, result.IsManager)
			}

			mockDB.AssertExpectations(t)
		})
	}
}
