package api

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"ems.dev/backend/services/conversation/types"
	templatetypes "ems.dev/backend/services/conversationtemplate/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/datatypes"
)

func TestCreateConversation(t *testing.T) {
	ctx := context.Background()
	organizationID := "org-123"
	managerMemberID := "manager-123"
	conversationDate := time.Now()

	tests := []struct {
		name           string
		req            *types.CreateConversationRequest
		templateID     *uuid.UUID
		template       *templatetypes.ConversationTemplate
		templateError  error
		dbError        error
		expectedError  error
		expectedStatus types.ConversationStatus
		validateResult func(t *testing.T, result *types.Conversation)
	}{
		{
			name: "success without template",
			req: &types.CreateConversationRequest{
				Title:            "Test Conversation",
				DirectMemberID:   "direct-123",
				ConversationDate: &conversationDate,
				Content:          datatypes.JSON(`{"key": "value"}`),
			},
			expectedStatus: types.ConversationStatusDraft,
			validateResult: func(t *testing.T, result *types.Conversation) {
				assert.Equal(t, "Test Conversation", result.Title)
				assert.Equal(t, organizationID, result.OrganizationID)
				assert.Equal(t, managerMemberID, result.ManagerMemberID)
				assert.Equal(t, "direct-123", result.DirectMemberID)
				assert.Equal(t, types.ConversationStatusDraft, result.Status)
			},
		},
		{
			name: "success with template",
			req: &types.CreateConversationRequest{
				Title:            "",
				DirectMemberID:   "direct-123",
				ConversationDate: &conversationDate,
			},
			templateID: uuidPtr(uuid.New()),
			template: &templatetypes.ConversationTemplate{
				ID:             uuid.New(),
				OrganizationID: uuid.New(),
				Name:           "Template Name",
				TemplateFields: []templatetypes.TemplateField{
					{ID: "field-1", Label: "Field 1", Type: "text", Required: true, Order: 1},
				},
			},
			expectedStatus: types.ConversationStatusDraft,
			validateResult: func(t *testing.T, result *types.Conversation) {
				assert.Equal(t, "Template Name", result.Title)
				var contentMap map[string]interface{}
				json.Unmarshal(result.Content, &contentMap)
				assert.NotNil(t, contentMap["_template_fields"])
			},
		},
		{
			name: "success with template and provided title",
			req: &types.CreateConversationRequest{
				Title:            "Custom Title",
				DirectMemberID:   "direct-123",
				ConversationDate: &conversationDate,
			},
			templateID: uuidPtr(uuid.New()),
			template: &templatetypes.ConversationTemplate{
				ID:             uuid.New(),
				OrganizationID: uuid.New(),
				Name:           "Template Name",
				TemplateFields: []templatetypes.TemplateField{
					{ID: "field-1", Label: "Field 1", Type: "text", Required: true, Order: 1},
				},
			},
			expectedStatus: types.ConversationStatusDraft,
			validateResult: func(t *testing.T, result *types.Conversation) {
				assert.Equal(t, "Custom Title", result.Title)
			},
		},
		{
			name: "success with empty content",
			req: &types.CreateConversationRequest{
				Title:          "Test Conversation",
				DirectMemberID: "direct-123",
			},
			expectedStatus: types.ConversationStatusDraft,
			validateResult: func(t *testing.T, result *types.Conversation) {
				var contentMap map[string]interface{}
				json.Unmarshal(result.Content, &contentMap)
				assert.NotNil(t, contentMap)
			},
		},
		{
			name: "database error",
			req: &types.CreateConversationRequest{
				Title:          "Test Conversation",
				DirectMemberID: "direct-123",
			},
			dbError:       errors.New("database error"),
			expectedError: errors.New("database error"),
		},
		{
			name: "invalid template ID",
			req: &types.CreateConversationRequest{
				Title:          "Test Conversation",
				DirectMemberID: "direct-123",
				TemplateID:     stringPtr("invalid-uuid"),
			},
			expectedStatus: types.ConversationStatusDraft,
			validateResult: func(t *testing.T, result *types.Conversation) {
				assert.Equal(t, "Test Conversation", result.Title)
			},
		},
		{
			name: "template not found",
			req: &types.CreateConversationRequest{
				Title:          "Test Conversation",
				DirectMemberID: "direct-123",
			},
			templateID:    uuidPtr(uuid.New()),
			templateError: errors.New("template not found"),
			expectedStatus: types.ConversationStatusDraft,
			validateResult: func(t *testing.T, result *types.Conversation) {
				assert.Equal(t, "Test Conversation", result.Title)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockConversationDB)
			mockTemplateAPI := new(MockConversationTemplateAPI)
			api := NewConversationAPI(mockDB, mockTemplateAPI)

			if tt.templateID != nil {
				templateIDStr := tt.templateID.String()
				tt.req.TemplateID = &templateIDStr
				mockTemplateAPI.On("GetConversationTemplate", *tt.templateID).Return(tt.template, tt.templateError)
			}

			mockDB.On("CreateConversation", ctx, mock.AnythingOfType("*types.Conversation")).Return(tt.dbError).Run(func(args mock.Arguments) {
				conv := args.Get(1).(*types.Conversation)
				conv.ID = "conv-123"
				if tt.validateResult != nil {
					tt.validateResult(t, conv)
				}
			})

			result, err := api.CreateConversation(ctx, organizationID, managerMemberID, tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedStatus, result.Status)
			}

			mockDB.AssertExpectations(t)
			if tt.templateID != nil {
				mockTemplateAPI.AssertExpectations(t)
			}
		})
	}
}
