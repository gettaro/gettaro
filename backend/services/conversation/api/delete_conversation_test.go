package api

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteConversation(t *testing.T) {
	ctx := context.Background()
	conversationID := "conv-123"

	tests := []struct {
		name          string
		id            string
		mockError     error
		expectedError error
	}{
		{
			name: "success",
			id:   conversationID,
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

			mockDB.On("DeleteConversation", ctx, tt.id).Return(tt.mockError)

			err := api.DeleteConversation(ctx, tt.id)

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
