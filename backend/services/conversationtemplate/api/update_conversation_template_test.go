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

func TestUpdateConversationTemplate(t *testing.T) {
	orgID := uuid.New()
	templateID := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		params         types.UpdateConversationTemplateParams
		mockTemplate   *types.ConversationTemplate
		mockError      error
		expectedResult *types.ConversationTemplate
		expectedError  error
	}{
		{
			name: "success - update all fields",
			params: types.UpdateConversationTemplateParams{
				ID:             templateID,
				Name:           stringPtr("Updated Name"),
				Description:    stringPtr("Updated Description"),
				TemplateFields: &[]types.TemplateField{
					{
						ID:       "field-1",
						Label:    "Updated Field",
						Type:     "textarea",
						Required: false,
						Order:    1,
					},
				},
				IsActive: boolPtr(false),
			},
			mockTemplate: &types.ConversationTemplate{
				ID:             templateID,
				OrganizationID: orgID,
				Name:           "Updated Name",
				Description:    stringPtr("Updated Description"),
				TemplateFields: []types.TemplateField{
					{
						ID:       "field-1",
						Label:    "Updated Field",
						Type:     "textarea",
						Required: false,
						Order:    1,
					},
				},
				IsActive:  false,
				CreatedAt: now,
				UpdatedAt: now,
			},
			expectedResult: &types.ConversationTemplate{
				ID:             templateID,
				OrganizationID: orgID,
				Name:           "Updated Name",
				Description:    stringPtr("Updated Description"),
				TemplateFields: []types.TemplateField{
					{
						ID:       "field-1",
						Label:    "Updated Field",
						Type:     "textarea",
						Required: false,
						Order:    1,
					},
				},
				IsActive:  false,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name: "success - update only name",
			params: types.UpdateConversationTemplateParams{
				ID:   templateID,
				Name: stringPtr("New Name"),
			},
			mockTemplate: &types.ConversationTemplate{
				ID:             templateID,
				OrganizationID: orgID,
				Name:           "New Name",
				Description:    stringPtr("Original Description"),
				TemplateFields: []types.TemplateField{},
				IsActive:       true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			expectedResult: &types.ConversationTemplate{
				ID:             templateID,
				OrganizationID: orgID,
				Name:           "New Name",
				Description:    stringPtr("Original Description"),
				TemplateFields: []types.TemplateField{},
				IsActive:       true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},
		{
			name: "success - update is_active only",
			params: types.UpdateConversationTemplateParams{
				ID:       templateID,
				IsActive: boolPtr(false),
			},
			mockTemplate: &types.ConversationTemplate{
				ID:             templateID,
				OrganizationID: orgID,
				Name:           "Original Name",
				TemplateFields: []types.TemplateField{},
				IsActive:       false,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			expectedResult: &types.ConversationTemplate{
				ID:             templateID,
				OrganizationID: orgID,
				Name:           "Original Name",
				TemplateFields: []types.TemplateField{},
				IsActive:       false,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},
		{
			name: "error - template not found",
			params: types.UpdateConversationTemplateParams{
				ID:   uuid.New(),
				Name: stringPtr("Updated Name"),
			},
			mockError:     gorm.ErrRecordNotFound,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "error - invalid UUID",
			params: types.UpdateConversationTemplateParams{
				ID:   uuid.Nil,
				Name: stringPtr("Updated Name"),
			},
			mockError:     errors.New("invalid UUID"),
			expectedError: errors.New("invalid UUID"),
		},
		{
			name: "error - database error",
			params: types.UpdateConversationTemplateParams{
				ID:   templateID,
				Name: stringPtr("Updated Name"),
			},
			mockError:     errors.New("database connection failed"),
			expectedError: errors.New("database connection failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockConversationTemplateDatabase)
			api := &ConversationTemplateAPI{db: mockDB}

			mockDB.On("UpdateConversationTemplate", tt.params).Return(tt.mockTemplate, tt.mockError)

			result, err := api.UpdateConversationTemplate(tt.params)

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
