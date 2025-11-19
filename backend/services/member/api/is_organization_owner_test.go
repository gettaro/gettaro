package api

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsOrganizationOwner(t *testing.T) {
	ctx := context.Background()
	orgID := "org-1"
	userID := "user-1"

	tests := []struct {
		name          string
		orgID         string
		userID        string
		mockResult    bool
		mockError     error
		expectedResult bool
		expectedError error
	}{
		{
			name:          "success - is owner",
			orgID:         orgID,
			userID:        userID,
			mockResult:    true,
			expectedResult: true,
		},
		{
			name:          "success - is not owner",
			orgID:         orgID,
			userID:        userID,
			mockResult:    false,
			expectedResult: false,
		},
		{
			name:          "database error",
			orgID:         orgID,
			userID:        userID,
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockMemberDB)
			mockUserAPI := new(MockUserAPI)
			mockDirectsAPI := new(MockDirectReportsAPI)
			mockSourceControlAPI := new(MockSourceControlAPI)
			mockTitleAPI := new(MockTitleAPI)

			api := NewApi(mockDB, mockUserAPI, mockSourceControlAPI, mockTitleAPI, mockDirectsAPI)

			mockDB.On("IsOrganizationOwner", tt.orgID, tt.userID).Return(tt.mockResult, tt.mockError)

			result, err := api.IsOrganizationOwner(ctx, tt.orgID, tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			mockDB.AssertExpectations(t)
		})
	}
}
