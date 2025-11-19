package api

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDeleteConversationTemplate(t *testing.T) {
	templateID := uuid.New()

	tests := []struct {
		name          string
		id            uuid.UUID
		mockError     error
		expectedError error
	}{
		{
			name: "success - delete existing template",
			id:   templateID,
		},
		{
			name:          "error - template not found",
			id:            uuid.New(),
			mockError:     gorm.ErrRecordNotFound,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:          "error - invalid UUID",
			id:            uuid.Nil,
			mockError:     errors.New("invalid UUID"),
			expectedError: errors.New("invalid UUID"),
		},
		{
			name:          "error - database error",
			id:            templateID,
			mockError:     errors.New("database connection failed"),
			expectedError: errors.New("database connection failed"),
		},
		{
			name:          "error - foreign key constraint",
			id:            templateID,
			mockError:     errors.New("foreign key constraint violation"),
			expectedError: errors.New("foreign key constraint violation"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockConversationTemplateDatabase)
			api := &ConversationTemplateAPI{db: mockDB}

			mockDB.On("DeleteConversationTemplate", tt.id).Return(tt.mockError)

			err := api.DeleteConversationTemplate(tt.id)

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
