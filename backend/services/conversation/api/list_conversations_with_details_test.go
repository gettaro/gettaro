package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"ems.dev/backend/services/conversation/types"
	"github.com/stretchr/testify/assert"
)

func TestListConversationsWithDetails(t *testing.T) {
	ctx := context.Background()
	organizationID := "org-123"

	tests := []struct {
		name            string
		organizationID  string
		query           *types.ListConversationsQuery
		mockDetails     []*types.ConversationWithDetails
		mockError       error
		expectedError   error
	}{
		{
			name:           "success with empty list",
			organizationID: organizationID,
			query:          &types.ListConversationsQuery{},
			mockDetails:    []*types.ConversationWithDetails{},
		},
		{
			name:           "success with conversations",
			organizationID: organizationID,
			query:          &types.ListConversationsQuery{},
			mockDetails: []*types.ConversationWithDetails{
				{
					Conversation: types.Conversation{
						ID:              "conv-1",
						OrganizationID:  organizationID,
						Title:           "Conversation 1",
						ManagerMemberID: "manager-123",
						DirectMemberID:  "direct-123",
						Status:          types.ConversationStatusDraft,
						CreatedAt:       time.Now(),
						UpdatedAt:       time.Now(),
					},
					Manager: &types.OrganizationMember{
						ID:             "manager-123",
						Email:          "manager@example.com",
						Username:       "manager",
						OrganizationID: organizationID,
					},
					DirectReport: &types.OrganizationMember{
						ID:             "direct-123",
						Email:          "direct@example.com",
						Username:       "direct",
						OrganizationID: organizationID,
					},
				},
				{
					Conversation: types.Conversation{
						ID:              "conv-2",
						OrganizationID:  organizationID,
						Title:           "Conversation 2",
						ManagerMemberID: "manager-123",
						DirectMemberID:  "direct-456",
						Status:          types.ConversationStatusCompleted,
						CreatedAt:       time.Now(),
						UpdatedAt:       time.Now(),
					},
					Manager: &types.OrganizationMember{
						ID:             "manager-123",
						Email:          "manager@example.com",
						Username:       "manager",
						OrganizationID: organizationID,
					},
					DirectReport: &types.OrganizationMember{
						ID:             "direct-456",
						Email:          "direct2@example.com",
						Username:       "direct2",
						OrganizationID: organizationID,
					},
				},
			},
		},
		{
			name:           "success with filters",
			organizationID: organizationID,
			query: &types.ListConversationsQuery{
				ManagerMemberID: stringPtr("manager-123"),
				Status:          stringPtr("draft"),
			},
			mockDetails: []*types.ConversationWithDetails{
				{
					Conversation: types.Conversation{
						ID:              "conv-1",
						OrganizationID:  organizationID,
						Title:           "Conversation 1",
						ManagerMemberID: "manager-123",
						DirectMemberID:  "direct-123",
						Status:          types.ConversationStatusDraft,
						CreatedAt:       time.Now(),
						UpdatedAt:       time.Now(),
					},
				},
			},
		},
		{
			name:           "database error",
			organizationID: organizationID,
			query:          &types.ListConversationsQuery{},
			mockError:      errors.New("database error"),
			expectedError:  errors.New("database error"),
		},
		{
			name:           "empty organization ID",
			organizationID: "",
			query:          &types.ListConversationsQuery{},
			mockDetails:    []*types.ConversationWithDetails{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockConversationDB)
			mockTemplateAPI := new(MockConversationTemplateAPI)
			api := NewConversationAPI(mockDB, mockTemplateAPI)

			mockDB.On("ListConversationsWithDetails", ctx, tt.organizationID, tt.query).Return(tt.mockDetails, tt.mockError)

			result, err := api.ListConversationsWithDetails(ctx, tt.organizationID, tt.query)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, len(tt.mockDetails), len(result))
				for i, detail := range result {
					assert.Equal(t, tt.mockDetails[i].ID, detail.ID)
					assert.Equal(t, tt.mockDetails[i].Title, detail.Title)
				}
			}

			mockDB.AssertExpectations(t)
		})
	}
}
