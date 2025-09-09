package database

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"ems.dev/backend/services/conversationtemplate/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TemplateFieldArray is a custom type for handling JSONB array of template fields
type TemplateFieldArray []types.TemplateField

// Scan implements the sql.Scanner interface
func (tfa *TemplateFieldArray) Scan(value interface{}) error {
	if value == nil {
		*tfa = TemplateFieldArray{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into TemplateFieldArray", value)
	}

	return json.Unmarshal(bytes, tfa)
}

// Value implements the driver.Valuer interface
func (tfa TemplateFieldArray) Value() (driver.Value, error) {
	if len(tfa) == 0 {
		return "[]", nil
	}
	return json.Marshal(tfa)
}

// ConversationTemplateDB represents the database model for conversation templates
type ConversationTemplateDB struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null"`
	Name           string    `gorm:"not null"`
	Description    *string
	TemplateFields TemplateFieldArray `gorm:"type:jsonb;not null;default:'[]'::jsonb"`
	IsActive       bool               `gorm:"not null;default:true"`
	CreatedAt      time.Time          `gorm:"not null;default:now()"`
	UpdatedAt      time.Time          `gorm:"not null;default:now()"`
}

// TableName returns the table name for ConversationTemplateDB
func (ConversationTemplateDB) TableName() string {
	return "conversation_templates"
}

// ConversationTemplateDatabase handles database operations for conversation templates
type ConversationTemplateDatabase struct {
	db *gorm.DB
}

// NewConversationTemplateDatabase creates a new ConversationTemplateDatabase instance
func NewConversationTemplateDatabase(db *gorm.DB) *ConversationTemplateDatabase {
	return &ConversationTemplateDatabase{db: db}
}

// CreateConversationTemplate creates a new conversation template
func (d *ConversationTemplateDatabase) CreateConversationTemplate(params types.CreateConversationTemplateParams) (*types.ConversationTemplate, error) {
	template := &ConversationTemplateDB{
		OrganizationID: params.OrganizationID,
		Name:           params.Name,
		Description:    params.Description,
		TemplateFields: TemplateFieldArray(params.TemplateFields),
		IsActive:       true,
	}

	if params.IsActive != nil {
		template.IsActive = *params.IsActive
	}

	if err := d.db.Create(template).Error; err != nil {
		return nil, err
	}

	return d.mapToConversationTemplate(template), nil
}

// GetConversationTemplate retrieves a conversation template by ID
func (d *ConversationTemplateDatabase) GetConversationTemplate(id uuid.UUID) (*types.ConversationTemplate, error) {
	var template ConversationTemplateDB
	if err := d.db.Where("id = ?", id).First(&template).Error; err != nil {
		return nil, err
	}

	return d.mapToConversationTemplate(&template), nil
}

// ListConversationTemplates retrieves conversation templates based on search parameters
func (d *ConversationTemplateDatabase) ListConversationTemplates(params types.ConversationTemplateSearchParams) ([]*types.ConversationTemplate, error) {
	var templates []ConversationTemplateDB
	query := d.db

	if params.OrganizationID != nil {
		query = query.Where("organization_id = ?", *params.OrganizationID)
	}

	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	if err := query.Order("created_at DESC").Find(&templates).Error; err != nil {
		return nil, err
	}

	result := make([]*types.ConversationTemplate, len(templates))
	for i, template := range templates {
		result[i] = d.mapToConversationTemplate(&template)
	}

	return result, nil
}

// UpdateConversationTemplate updates an existing conversation template
func (d *ConversationTemplateDatabase) UpdateConversationTemplate(params types.UpdateConversationTemplateParams) (*types.ConversationTemplate, error) {
	var template ConversationTemplateDB
	if err := d.db.Where("id = ?", params.ID).First(&template).Error; err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	if params.Name != nil {
		updates["name"] = *params.Name
	}

	if params.Description != nil {
		updates["description"] = *params.Description
	}

	if params.TemplateFields != nil {
		updates["template_fields"] = TemplateFieldArray(*params.TemplateFields)
	}

	if params.IsActive != nil {
		updates["is_active"] = *params.IsActive
	}

	if len(updates) > 0 {
		updates["updated_at"] = time.Now()
		if err := d.db.Model(&template).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	return d.mapToConversationTemplate(&template), nil
}

// DeleteConversationTemplate deletes a conversation template
func (d *ConversationTemplateDatabase) DeleteConversationTemplate(id uuid.UUID) error {
	return d.db.Where("id = ?", id).Delete(&ConversationTemplateDB{}).Error
}

// mapToConversationTemplate converts ConversationTemplateDB to types.ConversationTemplate
func (d *ConversationTemplateDatabase) mapToConversationTemplate(template *ConversationTemplateDB) *types.ConversationTemplate {
	return &types.ConversationTemplate{
		ID:             template.ID,
		OrganizationID: template.OrganizationID,
		Name:           template.Name,
		Description:    template.Description,
		TemplateFields: []types.TemplateField(template.TemplateFields),
		IsActive:       template.IsActive,
		CreatedAt:      template.CreatedAt,
		UpdatedAt:      template.UpdatedAt,
	}
}
