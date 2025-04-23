package database

import (
	"ems.dev/backend/services/organization/types"
	"gorm.io/gorm"
)

// IntegrationDB defines the interface for integration database operations
type IntegrationDB interface {
	CreateIntegrationConfig(config *types.IntegrationConfig) error
	GetIntegrationConfig(id string) (*types.IntegrationConfig, error)
	GetOrganizationIntegrationConfigs(orgID string) ([]types.IntegrationConfig, error)
	UpdateIntegrationConfig(config *types.IntegrationConfig) error
	DeleteIntegrationConfig(id string) error
}

type IntegrationDBImpl struct {
	db *gorm.DB
}

func NewIntegrationDB(db *gorm.DB) *IntegrationDBImpl {
	return &IntegrationDBImpl{
		db: db,
	}
}

// CreateIntegrationConfig creates a new integration config
func (d *IntegrationDBImpl) CreateIntegrationConfig(config *types.IntegrationConfig) error {
	return d.db.Create(config).Error
}

// GetIntegrationConfig retrieves an integration config by ID
func (d *IntegrationDBImpl) GetIntegrationConfig(id string) (*types.IntegrationConfig, error) {
	var config types.IntegrationConfig
	err := d.db.First(&config, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// GetOrganizationIntegrationConfigs retrieves all integration configs for an organization
func (d *IntegrationDBImpl) GetOrganizationIntegrationConfigs(orgID string) ([]types.IntegrationConfig, error) {
	var configs []types.IntegrationConfig
	err := d.db.Where("organization_id = ?", orgID).Find(&configs).Error
	if err != nil {
		return nil, err
	}
	return configs, nil
}

// UpdateIntegrationConfig updates an existing integration config
func (d *IntegrationDBImpl) UpdateIntegrationConfig(config *types.IntegrationConfig) error {
	return d.db.Save(config).Error
}

// DeleteIntegrationConfig deletes an integration config
func (d *IntegrationDBImpl) DeleteIntegrationConfig(id string) error {
	return d.db.Delete(&types.IntegrationConfig{}, "id = ?", id).Error
}
