package api

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveMemberTitle(t *testing.T) {
	tests := []struct {
		name          string
		memberID      string
		orgID         string
		mockError     error
		expectedError error
	}{
		{
			name:     "successful removal",
			memberID: "member-1",
			orgID:    "org-1",
		},
		{
			name:          "database error",
			memberID:      "member-1",
			orgID:         "org-1",
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
		{
			name:          "member title not found",
			memberID:      "member-1",
			orgID:         "org-1",
			mockError:     errors.New("record not found"),
			expectedError: errors.New("record not found"),
		},
		{
			name:          "invalid member id",
			memberID:      "",
			orgID:         "org-1",
			mockError:     errors.New("invalid member id"),
			expectedError: errors.New("invalid member id"),
		},
		{
			name:          "invalid org id",
			memberID:      "member-1",
			orgID:         "",
			mockError:     errors.New("invalid org id"),
			expectedError: errors.New("invalid org id"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockTitleDB)
			api := &Api{db: mockDB}

			mockDB.On("RemoveMemberTitle", tt.memberID, tt.orgID).Return(tt.mockError)

			err := api.RemoveMemberTitle(context.Background(), tt.memberID, tt.orgID)

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
