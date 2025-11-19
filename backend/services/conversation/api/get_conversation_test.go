package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"ems.dev/backend/services/conversation/types"
	"github.com/stretchr/testify/assert"
)

func TestGetConversation(t *testing.T) {
	ctx := context.Background()
	conversationID := "conv-123"

	tests := []struct {
		name          string
		id            string
		mockConversation *types.Conversation
		mockError     error
		expectedError error
	}{
		{
			name: "success",
			id:   conversationID,
			mockConversation: &types.Conversation{
				ID:              conversationID,
				OrganizationID:  "org-123",
				Title:           "Test Conversation",
				ManagerMemberID: "manager-123",
				DirectMemberID:  "direct-123",
				Status:          types.ConversationStatusDraft,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
		},
		{
			name:          "conversation not found",
			id:            conversationID,
			mockError:     errors.New("record not found"),
			expectedError: errors.New("record not found"),
		},
		{
			name:          "database error",
			id:            conversationID,
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
		{
			name:          "empty ID",
			id:            "",
			mockError:     errors.New("invalid ID"),
			expectedError: errors.New("invalid ID"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockConversationDB)
			mockTemplateAPI := new(MockConversationTemplateAPI)
			api := NewConversationAPI(mockDB, mockTemplateAPI)

			mockDB.On("GetConversation", ctx, tt.id).Return(tt.mockConversation, tt.mockError)

			result, err := api.GetConversation(ctx, tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockConversation.ID, result.ID)
				assert.Equal(t, tt.mockConversation.Title, result.Title)
			}

			mockDB.AssertExpectations(t)
		})
	}
}
