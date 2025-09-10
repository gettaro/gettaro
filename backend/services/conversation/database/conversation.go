package database

import (
	"context"

	"ems.dev/backend/services/conversation/types"
	templatetypes "ems.dev/backend/services/conversationtemplate/types"
	membertypes "ems.dev/backend/services/member/types"
	"gorm.io/gorm"
)

type ConversationDB struct {
	db *gorm.DB
}

func NewConversationDB(db *gorm.DB) *ConversationDB {
	return &ConversationDB{db: db}
}

// CreateConversation creates a new conversation
func (d *ConversationDB) CreateConversation(ctx context.Context, conversation *types.Conversation) error {
	return d.db.WithContext(ctx).Create(conversation).Error
}

// GetConversation retrieves a conversation by ID
func (d *ConversationDB) GetConversation(ctx context.Context, id string) (*types.Conversation, error) {
	var conversation types.Conversation
	err := d.db.WithContext(ctx).Where("id = ?", id).First(&conversation).Error
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

// GetConversationWithDetails retrieves a conversation with related data
func (d *ConversationDB) GetConversationWithDetails(ctx context.Context, id string) (*types.ConversationWithDetails, error) {
	var conversation types.Conversation
	err := d.db.WithContext(ctx).Where("id = ?", id).First(&conversation).Error
	if err != nil {
		return nil, err
	}

	details := &types.ConversationWithDetails{
		Conversation: conversation,
	}

	// Load template if template_id exists
	if conversation.TemplateID != nil {
		var template templatetypes.ConversationTemplate
		err = d.db.WithContext(ctx).Where("id = ?", *conversation.TemplateID).First(&template).Error
		if err == nil {
			// Convert template to our types
			conversationTemplate := &types.ConversationTemplate{
				ID:             template.ID.String(),
				OrganizationID: template.OrganizationID.String(),
				Name:           template.Name,
				Description:    template.Description,
				IsActive:       template.IsActive,
				CreatedAt:      template.CreatedAt,
				UpdatedAt:      template.UpdatedAt,
			}

			// Convert template fields
			for _, field := range template.TemplateFields {
				var placeholder *string
				if field.Placeholder != "" {
					placeholder = &field.Placeholder
				}
				conversationTemplate.TemplateFields = append(conversationTemplate.TemplateFields, types.TemplateField{
					ID:          field.ID,
					Label:       field.Label,
					Type:        field.Type,
					Required:    field.Required,
					Placeholder: placeholder,
					Order:       field.Order,
				})
			}

			details.Template = conversationTemplate
		}
	}

	// Load manager details
	var manager membertypes.OrganizationMember
	err = d.db.WithContext(ctx).Where("id = ?", conversation.ManagerMemberID).First(&manager).Error
	if err == nil {
		details.Manager = &types.OrganizationMember{
			ID:             manager.ID,
			UserID:         manager.UserID,
			Email:          manager.Email,
			Username:       manager.Username,
			OrganizationID: manager.OrganizationID,
			IsOwner:        manager.IsOwner,
			TitleID:        manager.TitleID,
			ManagerID:      manager.ManagerID,
			CreatedAt:      manager.CreatedAt,
			UpdatedAt:      manager.UpdatedAt,
		}
	}

	// Load direct report details
	var directReport membertypes.OrganizationMember
	err = d.db.WithContext(ctx).Where("id = ?", conversation.DirectMemberID).First(&directReport).Error
	if err == nil {
		details.DirectReport = &types.OrganizationMember{
			ID:             directReport.ID,
			UserID:         directReport.UserID,
			Email:          directReport.Email,
			Username:       directReport.Username,
			OrganizationID: directReport.OrganizationID,
			IsOwner:        directReport.IsOwner,
			TitleID:        directReport.TitleID,
			ManagerID:      directReport.ManagerID,
			CreatedAt:      directReport.CreatedAt,
			UpdatedAt:      directReport.UpdatedAt,
		}
	}

	return details, nil
}

// ListConversations retrieves conversations with optional filters
func (d *ConversationDB) ListConversations(ctx context.Context, organizationID string, query *types.ListConversationsQuery) ([]*types.Conversation, error) {
	var conversations []*types.Conversation

	db := d.db.WithContext(ctx).Where("organization_id = ?", organizationID)

	if query.ManagerMemberID != nil {
		db = db.Where("manager_member_id = ?", *query.ManagerMemberID)
	}
	if query.DirectMemberID != nil {
		db = db.Where("direct_member_id = ?", *query.DirectMemberID)
	}
	if query.TemplateID != nil {
		db = db.Where("template_id = ?", *query.TemplateID)
	}
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}

	// Apply pagination
	if query.Limit != nil {
		db = db.Limit(*query.Limit)
	}
	if query.Offset != nil {
		db = db.Offset(*query.Offset)
	}

	err := db.Order("created_at DESC").Find(&conversations).Error
	return conversations, err
}

// ListConversationsWithDetails retrieves conversations with related data
func (d *ConversationDB) ListConversationsWithDetails(ctx context.Context, organizationID string, query *types.ListConversationsQuery) ([]*types.ConversationWithDetails, error) {
	conversations, err := d.ListConversations(ctx, organizationID, query)
	if err != nil {
		return nil, err
	}

	var details []*types.ConversationWithDetails
	for _, conv := range conversations {
		convDetails, err := d.GetConversationWithDetails(ctx, conv.ID)
		if err != nil {
			continue // Skip conversations with errors
		}
		details = append(details, convDetails)
	}

	return details, nil
}

// UpdateConversation updates a conversation
func (d *ConversationDB) UpdateConversation(ctx context.Context, id string, updates *types.UpdateConversationRequest) error {
	updatesMap := make(map[string]interface{})

	if updates.ConversationDate != nil {
		updatesMap["conversation_date"] = updates.ConversationDate
	}
	if updates.Status != nil {
		updatesMap["status"] = *updates.Status
	}
	if updates.Content != nil {
		updatesMap["content"] = updates.Content
	}

	return d.db.WithContext(ctx).Model(&types.Conversation{}).Where("id = ?", id).Updates(updatesMap).Error
}

// DeleteConversation deletes a conversation
func (d *ConversationDB) DeleteConversation(ctx context.Context, id string) error {
	return d.db.WithContext(ctx).Where("id = ?", id).Delete(&types.Conversation{}).Error
}

// GetConversationStats returns statistics for conversations
func (d *ConversationDB) GetConversationStats(ctx context.Context, organizationID string, managerMemberID *string) (map[string]int, error) {
	stats := make(map[string]int)

	db := d.db.WithContext(ctx).Model(&types.Conversation{}).Where("organization_id = ?", organizationID)
	if managerMemberID != nil {
		db = db.Where("manager_member_id = ?", *managerMemberID)
	}

	// Count by status
	var results []struct {
		Status string
		Count  int
	}

	err := db.Select("status, COUNT(*) as count").Group("status").Find(&results).Error
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		stats[result.Status] = result.Count
	}

	// Count total
	var totalCount int64
	err = db.Count(&totalCount).Error
	if err != nil {
		return nil, err
	}
	stats["total"] = int(totalCount)

	return stats, nil
}
