package api

import (
	"errors"
	"testing"
	"time"

	"ems.dev/backend/services/conversationtemplate/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetConversationTemplate(t *testing.T) {
	orgID := uuid.New()
	templateID := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		id             uuid.UUID
		mockTemplate   *types.ConversationTemplate
		mockError      error
		expectedResult *types.ConversationTemplate
		expectedError  error
	}{
		{
			name: "success - get existing template",
			id:   templateID,
			mockTemplate: &types.ConversationTemplate{
				ID:             templateID,
				OrganizationID: orgID,
				Name:           "Test Template",
				Description:    stringPtr("Test Description"),
				TemplateFields: []types.TemplateField{
					{
						ID:       "field-1",
						Label:    "Name",
						Type:     "text",
						Required: true,
						Order:    1,
					},
				},
				IsActive:  true,
				CreatedAt: now,
				UpdatedAt: now,
			},
			expectedResult: &types.ConversationTemplate{
				ID:             templateID,
				OrganizationID: orgID,
				Name:           "Test Template",
				Description:    stringPtr("Test Description"),
				TemplateFields: []types.TemplateField{
					{
						ID:       "field-1",
						Label:    "Name",
						Type:     "text",
						Required: true,
						Order:    1,
					},
				},
				IsActive:  true,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name:           "error - template not found",
			id:             uuid.New(),
			mockError:      gorm.ErrRecordNotFound,
			expectedError:  gorm.ErrRecordNotFound,
		},
		{
			name:           "error - invalid UUID format",
			id:             uuid.Nil,
			mockError:      errors.New("invalid UUID"),
			expectedError:  errors.New("invalid UUID"),
		},
		{
			name:           "error - database error",
			id:             templateID,
			mockError:      errors.New("database connection failed"),
			expectedError:  errors.New("database connection failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockConversationTemplateDatabase)
			api := &ConversationTemplateAPI{db: mockDB}

			mockDB.On("GetConversationTemplate", tt.id).Return(tt.mockTemplate, tt.mockError)

			result, err := api.GetConversationTemplate(tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			mockDB.AssertExpectations(t)
		})
	}
}
