package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"ems.dev/backend/services/conversation/types"
	"github.com/stretchr/testify/assert"
)

func TestGetConversationWithDetails(t *testing.T) {
	ctx := context.Background()
	conversationID := "conv-123"

	tests := []struct {
		name          string
		id            string
		mockDetails   *types.ConversationWithDetails
		mockError     error
		expectedError error
	}{
		{
			name: "success with all details",
			id:   conversationID,
			mockDetails: &types.ConversationWithDetails{
				Conversation: types.Conversation{
					ID:              conversationID,
					OrganizationID:  "org-123",
					Title:           "Test Conversation",
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
					OrganizationID: "org-123",
				},
				DirectReport: &types.OrganizationMember{
					ID:             "direct-123",
					Email:          "direct@example.com",
					Username:       "direct",
					OrganizationID: "org-123",
				},
			},
		},
		{
			name: "success with template",
			id:   conversationID,
			mockDetails: &types.ConversationWithDetails{
				Conversation: types.Conversation{
					ID:              conversationID,
					OrganizationID:  "org-123",
					Title:           "Test Conversation",
					ManagerMemberID: "manager-123",
					DirectMemberID:  "direct-123",
					Status:          types.ConversationStatusDraft,
					TemplateID:      stringPtr("template-123"),
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				},
				Template: &types.ConversationTemplate{
					ID:             "template-123",
					OrganizationID: "org-123",
					Name:           "Template Name",
					IsActive:       true,
				},
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

			mockDB.On("GetConversationWithDetails", ctx, tt.id).Return(tt.mockDetails, tt.mockError)

			result, err := api.GetConversationWithDetails(ctx, tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockDetails.ID, result.ID)
				assert.Equal(t, tt.mockDetails.Title, result.Title)
				if tt.mockDetails.Template != nil {
					assert.NotNil(t, result.Template)
					assert.Equal(t, tt.mockDetails.Template.Name, result.Template.Name)
				}
			}

			mockDB.AssertExpectations(t)
		})
	}
}
