package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"ems.dev/backend/services/title/types"
	"github.com/stretchr/testify/assert"
)

func TestListTitles(t *testing.T) {
	tests := []struct {
		name           string
		orgID          string
		mockTitles     []types.Title
		mockError      error
		expectedTitles []types.Title
		expectedError  error
	}{
		{
			name:  "successful retrieval",
			orgID: "org-1",
			mockTitles: []types.Title{
				{
					ID:             "title-1",
					Name:           "Software Engineer",
					OrganizationID: "org-1",
					IsManager:      false,
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				},
				{
					ID:             "title-2",
					Name:           "Senior Engineer",
					OrganizationID: "org-1",
					IsManager:      true,
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				},
			},
			expectedTitles: []types.Title{
				{
					ID:             "title-1",
					Name:           "Software Engineer",
					OrganizationID: "org-1",
					IsManager:      false,
				},
				{
					ID:             "title-2",
					Name:           "Senior Engineer",
					OrganizationID: "org-1",
					IsManager:      true,
				},
			},
		},
		{
			name:           "empty list",
			orgID:          "org-1",
			mockTitles:     []types.Title{},
			expectedTitles: []types.Title{},
		},
		{
			name:        "database error",
			orgID:       "org-1",
			mockError:   errors.New("database error"),
			expectedError: errors.New("database error"),
		},
		{
			name:           "invalid org id",
			orgID:          "",
			mockTitles:     []types.Title{},
			expectedTitles: []types.Title{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockTitleDB)
			api := &Api{db: mockDB}

			mockDB.On("ListTitles", tt.orgID).Return(tt.mockTitles, tt.mockError)

			result, err := api.ListTitles(context.Background(), tt.orgID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, len(tt.expectedTitles), len(result))
				for i, expected := range tt.expectedTitles {
					assert.Equal(t, expected.ID, result[i].ID)
					assert.Equal(t, expected.Name, result[i].Name)
					assert.Equal(t, expected.OrganizationID, result[i].OrganizationID)
					assert.Equal(t, expected.IsManager, result[i].IsManager)
				}
			}

			mockDB.AssertExpectations(t)
		})
	}
}
