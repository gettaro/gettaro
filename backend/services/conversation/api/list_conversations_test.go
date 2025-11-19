package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"ems.dev/backend/services/conversation/types"
	"github.com/stretchr/testify/assert"
)

func TestListConversations(t *testing.T) {
	ctx := context.Background()
	organizationID := "org-123"

	tests := []struct {
		name            string
		organizationID  string
		query           *types.ListConversationsQuery
		mockConversations []*types.Conversation
		mockError       error
		expectedError   error
	}{
		{
			name:           "success with empty list",
			organizationID: organizationID,
			query:          &types.ListConversationsQuery{},
			mockConversations: []*types.Conversation{},
		},
		{
			name:           "success with conversations",
			organizationID: organizationID,
			query:          &types.ListConversationsQuery{},
			mockConversations: []*types.Conversation{
				{
					ID:              "conv-1",
					OrganizationID:  organizationID,
					Title:           "Conversation 1",
					ManagerMemberID: "manager-123",
					DirectMemberID:  "direct-123",
					Status:          types.ConversationStatusDraft,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				},
				{
					ID:              "conv-2",
					OrganizationID:  organizationID,
					Title:           "Conversation 2",
					ManagerMemberID: "manager-123",
					DirectMemberID:  "direct-456",
					Status:          types.ConversationStatusCompleted,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				},
			},
		},
		{
			name:           "success with filters",
			organizationID: organizationID,
			query: &types.ListConversationsQuery{
				ManagerMemberID: stringPtr("manager-123"),
				Status:          stringPtr("draft"),
				Limit:           intPtr(10),
				Offset:          intPtr(0),
			},
			mockConversations: []*types.Conversation{
				{
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
			mockConversations: []*types.Conversation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockConversationDB)
			mockTemplateAPI := new(MockConversationTemplateAPI)
			api := NewConversationAPI(mockDB, mockTemplateAPI)

			mockDB.On("ListConversations", ctx, tt.organizationID, tt.query).Return(tt.mockConversations, tt.mockError)

			result, err := api.ListConversations(ctx, tt.organizationID, tt.query)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, len(tt.mockConversations), len(result))
				for i, conv := range result {
					assert.Equal(t, tt.mockConversations[i].ID, conv.ID)
					assert.Equal(t, tt.mockConversations[i].Title, conv.Title)
				}
			}

			mockDB.AssertExpectations(t)
		})
	}
}
