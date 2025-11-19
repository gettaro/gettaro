package api

import (
	"errors"
	"testing"
	"time"

	"ems.dev/backend/services/conversationtemplate/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockConversationTemplateDatabase is a mock implementation of ConversationTemplateDatabase
type MockConversationTemplateDatabase struct {
	mock.Mock
}

func (m *MockConversationTemplateDatabase) CreateConversationTemplate(params types.CreateConversationTemplateParams) (*types.ConversationTemplate, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.ConversationTemplate), args.Error(1)
}

func (m *MockConversationTemplateDatabase) GetConversationTemplate(id uuid.UUID) (*types.ConversationTemplate, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.ConversationTemplate), args.Error(1)
}

func (m *MockConversationTemplateDatabase) ListConversationTemplates(params types.ConversationTemplateSearchParams) ([]*types.ConversationTemplate, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*types.ConversationTemplate), args.Error(1)
}

func (m *MockConversationTemplateDatabase) UpdateConversationTemplate(params types.UpdateConversationTemplateParams) (*types.ConversationTemplate, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.ConversationTemplate), args.Error(1)
}

func (m *MockConversationTemplateDatabase) DeleteConversationTemplate(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateConversationTemplate(t *testing.T) {
	orgID := uuid.New()
	templateID := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		params         types.CreateConversationTemplateParams
		mockTemplate   *types.ConversationTemplate
		mockError      error
		expectedResult *types.ConversationTemplate
		expectedError  error
	}{
		{
			name: "success - create template with all fields",
			params: types.CreateConversationTemplateParams{
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
				IsActive: boolPtr(true),
			},
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
			name: "success - create template with minimal fields",
			params: types.CreateConversationTemplateParams{
				OrganizationID: orgID,
				Name:           "Minimal Template",
				TemplateFields: []types.TemplateField{},
			},
			mockTemplate: &types.ConversationTemplate{
				ID:             templateID,
				OrganizationID: orgID,
				Name:           "Minimal Template",
				Description:    nil,
				TemplateFields: []types.TemplateField{},
				IsActive:       true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			expectedResult: &types.ConversationTemplate{
				ID:             templateID,
				OrganizationID: orgID,
				Name:           "Minimal Template",
				Description:    nil,
				TemplateFields: []types.TemplateField{},
				IsActive:       true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},
		{
			name: "error - database error",
			params: types.CreateConversationTemplateParams{
				OrganizationID: orgID,
				Name:           "Test Template",
				TemplateFields: []types.TemplateField{},
			},
			mockError:     errors.New("database connection failed"),
			expectedError: errors.New("database connection failed"),
		},
		{
			name: "error - duplicate name",
			params: types.CreateConversationTemplateParams{
				OrganizationID: orgID,
				Name:           "Duplicate Template",
				TemplateFields: []types.TemplateField{},
			},
			mockError:     errors.New("duplicate key value violates unique constraint"),
			expectedError: errors.New("duplicate key value violates unique constraint"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockConversationTemplateDatabase)
			api := &ConversationTemplateAPI{db: mockDB}

			mockDB.On("CreateConversationTemplate", tt.params).Return(tt.mockTemplate, tt.mockError)

			result, err := api.CreateConversationTemplate(tt.params)

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

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
