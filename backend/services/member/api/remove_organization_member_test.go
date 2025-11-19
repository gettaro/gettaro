package api

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveOrganizationMember(t *testing.T) {
	ctx := context.Background()
	orgID := "org-1"
	userID := "user-1"

	tests := []struct {
		name          string
		orgID         string
		userID        string
		mockError     error
		expectedError error
	}{
		{
			name:          "success",
			orgID:         orgID,
			userID:        userID,
			mockError:     nil,
			expectedError: nil,
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

			mockDB.On("RemoveOrganizationMember", tt.orgID, tt.userID).Return(tt.mockError)

			err := api.RemoveOrganizationMember(ctx, tt.orgID, tt.userID)

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
