package api

import (
	"errors"
	"testing"
	"time"

	"ems.dev/backend/services/conversationtemplate/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestListConversationTemplates(t *testing.T) {
	orgID := uuid.New()
	templateID1 := uuid.New()
	templateID2 := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		params         types.ConversationTemplateSearchParams
		mockTemplates  []*types.ConversationTemplate
		mockError      error
		expectedResult []*types.ConversationTemplate
		expectedError  error
	}{
		{
			name: "success - list all templates",
			params: types.ConversationTemplateSearchParams{},
			mockTemplates: []*types.ConversationTemplate{
				{
					ID:             templateID1,
					OrganizationID: orgID,
					Name:           "Template 1",
					Description:    stringPtr("Description 1"),
					TemplateFields: []types.TemplateField{},
					IsActive:       true,
					CreatedAt:      now,
					UpdatedAt:      now,
				},
				{
					ID:             templateID2,
					OrganizationID: orgID,
					Name:           "Template 2",
					Description:    nil,
					TemplateFields: []types.TemplateField{},
					IsActive:       false,
					CreatedAt:      now,
					UpdatedAt:      now,
				},
			},
			expectedResult: []*types.ConversationTemplate{
				{
					ID:             templateID1,
					OrganizationID: orgID,
					Name:           "Template 1",
					Description:    stringPtr("Description 1"),
					TemplateFields: []types.TemplateField{},
					IsActive:       true,
					CreatedAt:      now,
					UpdatedAt:      now,
				},
				{
					ID:             templateID2,
					OrganizationID: orgID,
					Name:           "Template 2",
					Description:    nil,
					TemplateFields: []types.TemplateField{},
					IsActive:       false,
					CreatedAt:      now,
					UpdatedAt:      now,
				},
			},
		},
		{
			name: "success - list templates by organization",
			params: types.ConversationTemplateSearchParams{
				OrganizationID: &orgID,
			},
			mockTemplates: []*types.ConversationTemplate{
				{
					ID:             templateID1,
					OrganizationID: orgID,
					Name:           "Template 1",
					TemplateFields: []types.TemplateField{},
					IsActive:       true,
					CreatedAt:      now,
					UpdatedAt:      now,
				},
			},
			expectedResult: []*types.ConversationTemplate{
				{
					ID:             templateID1,
					OrganizationID: orgID,
					Name:           "Template 1",
					TemplateFields: []types.TemplateField{},
					IsActive:       true,
					CreatedAt:      now,
					UpdatedAt:      now,
				},
			},
		},
		{
			name: "success - list active templates only",
			params: types.ConversationTemplateSearchParams{
				IsActive: boolPtr(true),
			},
			mockTemplates: []*types.ConversationTemplate{
				{
					ID:             templateID1,
					OrganizationID: orgID,
					Name:           "Active Template",
					TemplateFields: []types.TemplateField{},
					IsActive:       true,
					CreatedAt:      now,
					UpdatedAt:      now,
				},
			},
			expectedResult: []*types.ConversationTemplate{
				{
					ID:             templateID1,
					OrganizationID: orgID,
					Name:           "Active Template",
					TemplateFields: []types.TemplateField{},
					IsActive:       true,
					CreatedAt:      now,
					UpdatedAt:      now,
				},
			},
		},
		{
			name: "success - empty list",
			params: types.ConversationTemplateSearchParams{
				OrganizationID: &orgID,
			},
			mockTemplates:  []*types.ConversationTemplate{},
			expectedResult: []*types.ConversationTemplate{},
		},
		{
			name: "error - database error",
			params: types.ConversationTemplateSearchParams{
				OrganizationID: &orgID,
			},
			mockError:     errors.New("database connection failed"),
			expectedError: errors.New("database connection failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockConversationTemplateDatabase)
			api := &ConversationTemplateAPI{db: mockDB}

			mockDB.On("ListConversationTemplates", tt.params).Return(tt.mockTemplates, tt.mockError)

			result, err := api.ListConversationTemplates(tt.params)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
				// Ensure empty arrays are returned, not nil
				if len(tt.expectedResult) == 0 {
					assert.NotNil(t, result)
					assert.Len(t, result, 0)
				}
			}

			mockDB.AssertExpectations(t)
		})
	}
}
