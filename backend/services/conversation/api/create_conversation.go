package api

import (
	"context"
	"encoding/json"

	"ems.dev/backend/services/conversation/types"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// CreateConversation creates a new conversation
func (a *ConversationAPI) CreateConversation(ctx context.Context, organizationID string, managerMemberID string, req *types.CreateConversationRequest) (*types.Conversation, error) {
	// Start with the provided content
	var contentJSON datatypes.JSON
	if req.Content != nil {
		contentJSON = req.Content
	} else {
		contentJSON = datatypes.JSON("{}")
	}

	// Parse content to add template metadata if needed
	var contentMap map[string]interface{}
	if err := json.Unmarshal(contentJSON, &contentMap); err != nil {
		contentMap = make(map[string]interface{})
	}

	// Determine the title - use provided title or default to template name
	title := req.Title
	if req.TemplateID != nil {
		templateID, err := uuid.Parse(*req.TemplateID)
		if err == nil {
			template, err := a.templateApi.GetConversationTemplate(templateID)
			if err == nil && template != nil {
				// Add template fields metadata
				contentMap["_template_fields"] = template.TemplateFields

				// If no title provided, use template name
				if title == "" {
					title = template.Name
				}
			}
		}
	}

	// Convert contentMap back to datatypes.JSON
	var err error
	contentJSON, err = json.Marshal(contentMap)
	if err != nil {
		return nil, err
	}

	conversation := &types.Conversation{
		OrganizationID:   organizationID,
		TemplateID:       req.TemplateID,
		Title:            title,
		ManagerMemberID:  managerMemberID,
		DirectMemberID:   req.DirectMemberID,
		ConversationDate: req.ConversationDate,
		Status:           types.ConversationStatusDraft,
		Content:          contentJSON,
	}

	err = a.db.CreateConversation(ctx, conversation)
	if err != nil {
		return nil, err
	}

	return conversation, nil
}
