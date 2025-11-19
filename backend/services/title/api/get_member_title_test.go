package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"ems.dev/backend/services/title/types"
	"github.com/stretchr/testify/assert"
)

func TestGetMemberTitle(t *testing.T) {
	tests := []struct {
		name            string
		memberID        string
		orgID           string
		mockMemberTitle *types.MemberTitle
		mockError       error
		expectedTitle   *types.MemberTitle
		expectedError   error
	}{
		{
			name:     "successful retrieval",
			memberID: "member-1",
			orgID:    "org-1",
			mockMemberTitle: &types.MemberTitle{
				ID:             "member-title-1",
				MemberID:       "member-1",
				TitleID:        "title-1",
				OrganizationID: "org-1",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			expectedTitle: &types.MemberTitle{
				ID:             "member-title-1",
				MemberID:       "member-1",
				TitleID:        "title-1",
				OrganizationID: "org-1",
			},
		},
		{
			name:            "member title not found",
			memberID:        "member-1",
			orgID:           "org-1",
			mockMemberTitle: nil,
			expectedTitle:   nil,
		},
		{
			name:          "database error",
			memberID:      "member-1",
			orgID:         "org-1",
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
		{
			name:            "invalid member id",
			memberID:        "",
			orgID:           "org-1",
			mockMemberTitle: nil,
			expectedTitle:   nil,
		},
		{
			name:            "invalid org id",
			memberID:        "member-1",
			orgID:           "",
			mockMemberTitle: nil,
			expectedTitle:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockTitleDB)
			api := &Api{db: mockDB}

			mockDB.On("GetMemberTitle", tt.memberID, tt.orgID).Return(tt.mockMemberTitle, tt.mockError)

			result, err := api.GetMemberTitle(context.Background(), tt.memberID, tt.orgID)

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
					assert.Equal(t, tt.expectedTitle.MemberID, result.MemberID)
					assert.Equal(t, tt.expectedTitle.TitleID, result.TitleID)
					assert.Equal(t, tt.expectedTitle.OrganizationID, result.OrganizationID)
				}
			}

			mockDB.AssertExpectations(t)
		})
	}
}
