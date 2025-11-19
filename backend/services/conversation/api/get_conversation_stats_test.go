package api

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConversationStats(t *testing.T) {
	ctx := context.Background()
	organizationID := "org-123"
	managerMemberID := "manager-123"

	tests := []struct {
		name          string
		organizationID string
		managerMemberID *string
		mockStats     map[string]int
		mockError     error
		expectedError error
	}{
		{
			name:           "success with all stats",
			organizationID: organizationID,
			managerMemberID: nil,
			mockStats: map[string]int{
				"draft":     5,
				"completed": 10,
				"total":     15,
			},
		},
		{
			name:           "success with manager filter",
			organizationID: organizationID,
			managerMemberID: &managerMemberID,
			mockStats: map[string]int{
				"draft":     2,
				"completed": 3,
				"total":     5,
			},
		},
		{
			name:           "success with empty stats",
			organizationID: organizationID,
			managerMemberID: nil,
			mockStats: map[string]int{
				"total": 0,
			},
		},
		{
			name:           "success with only draft",
			organizationID: organizationID,
			managerMemberID: nil,
			mockStats: map[string]int{
				"draft": 5,
				"total": 5,
			},
		},
		{
			name:           "success with only completed",
			organizationID: organizationID,
			managerMemberID: nil,
			mockStats: map[string]int{
				"completed": 10,
				"total":     10,
			},
		},
		{
			name:           "database error",
			organizationID: organizationID,
			managerMemberID: nil,
			mockError:      errors.New("database error"),
			expectedError:  errors.New("database error"),
		},
		{
			name:           "empty organization ID",
			organizationID: "",
			managerMemberID: nil,
			mockStats: map[string]int{
				"total": 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockConversationDB)
			mockTemplateAPI := new(MockConversationTemplateAPI)
			api := NewConversationAPI(mockDB, mockTemplateAPI)

			mockDB.On("GetConversationStats", ctx, tt.organizationID, tt.managerMemberID).Return(tt.mockStats, tt.mockError)

			result, err := api.GetConversationStats(ctx, tt.organizationID, tt.managerMemberID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockStats, result)
				// Verify total is always present
				if len(tt.mockStats) > 0 {
					_, hasTotal := result["total"]
					assert.True(t, hasTotal || len(result) == 0)
				}
			}

			mockDB.AssertExpectations(t)
		})
	}
}
