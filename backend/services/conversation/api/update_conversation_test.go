package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"ems.dev/backend/services/conversation/types"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestUpdateConversation(t *testing.T) {
	ctx := context.Background()
	conversationID := "conv-123"
	conversationDate := time.Now()

	tests := []struct {
		name          string
		id            string
		req           *types.UpdateConversationRequest
		mockError     error
		expectedError error
	}{
		{
			name: "success update status",
			id:   conversationID,
			req: &types.UpdateConversationRequest{
				Status: statusPtr(types.ConversationStatusCompleted),
			},
		},
		{
			name: "success update date",
			id:   conversationID,
			req: &types.UpdateConversationRequest{
				ConversationDate: &conversationDate,
			},
		},
		{
			name: "success update content",
			id:   conversationID,
			req: &types.UpdateConversationRequest{
				Content: datatypes.JSON(`{"key": "value"}`),
			},
		},
		{
			name: "success update all fields",
			id:   conversationID,
			req: &types.UpdateConversationRequest{
				ConversationDate: &conversationDate,
				Status:           statusPtr(types.ConversationStatusCompleted),
				Content:          datatypes.JSON(`{"key": "value"}`),
			},
		},
		{
			name:          "conversation not found",
			id:            conversationID,
			req:           &types.UpdateConversationRequest{Status: statusPtr(types.ConversationStatusCompleted)},
			mockError:     errors.New("record not found"),
			expectedError: errors.New("record not found"),
		},
		{
			name:          "database error",
			id:            conversationID,
			req:           &types.UpdateConversationRequest{Status: statusPtr(types.ConversationStatusCompleted)},
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
		{
			name:          "empty ID",
			id:            "",
			req:           &types.UpdateConversationRequest{Status: statusPtr(types.ConversationStatusCompleted)},
			mockError:     errors.New("invalid ID"),
			expectedError: errors.New("invalid ID"),
		},
		{
			name:          "empty request",
			id:            conversationID,
			req:           &types.UpdateConversationRequest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockConversationDB)
			mockTemplateAPI := new(MockConversationTemplateAPI)
			api := NewConversationAPI(mockDB, mockTemplateAPI)

			mockDB.On("UpdateConversation", ctx, tt.id, tt.req).Return(tt.mockError)

			err := api.UpdateConversation(ctx, tt.id, tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
		})
	}
}
