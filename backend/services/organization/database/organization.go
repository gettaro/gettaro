package database

import (
	"strings"

	"ems.dev/backend/services/errors"
	"ems.dev/backend/services/organization/types"
	"gorm.io/gorm"
)

// DB defines the interface for organization database operations
type DB interface {
	CreateOrganization(org *types.Organization, ownerID string) error
	GetMemberOrganizations(userID string) ([]types.Organization, error)
	GetOrganizations() ([]types.Organization, error)
	GetOrganizationByID(id string) (*types.Organization, error)
	UpdateOrganization(org *types.Organization) error
	DeleteOrganization(id string) error
}

type OrganizationDB struct {
	db *gorm.DB
}

func NewOrganizationDB(db *gorm.DB) *OrganizationDB {
	return &OrganizationDB{
		db: db,
	}
}

// CreateOrganization creates a new organization and sets the specified user as its owner
func (d *OrganizationDB) CreateOrganization(org *types.Organization, userID string) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		// Create organization
		err := tx.Create(org).Error
		if err != nil {
			// Check for unique constraint violation on slug
			if strings.Contains(err.Error(), "duplicate key value") {
				return &errors.ErrDuplicateConflict{
					Resource: "organization",
					Field:    "slug",
					Value:    org.Slug,
				}
			}
			return err
		}

		// Create user-organization relationship with owner flag
		if err := tx.Exec(
			"INSERT INTO organization_members (user_id, organization_id, is_owner) VALUES (?, ?, true)",
			userID,
			org.ID,
		).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetMemberOrganizations returns all organizations a user is part of, with ownership information
func (d *OrganizationDB) GetMemberOrganizations(userID string) ([]types.Organization, error) {
	type result struct {
		types.Organization
		IsOwner bool `gorm:"column:is_owner"`
	}

	var results []result
	err := d.db.Raw(`
		SELECT o.*, om.is_owner
		FROM organizations o
		JOIN organization_members om ON o.id = om.organization_id
		WHERE om.user_id = ?
	`, userID).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Convert results to []types.Organization
	orgs := make([]types.Organization, len(results))
	for i, r := range results {
		r.Organization.IsOwner = r.IsOwner
		orgs[i] = r.Organization
	}

	return orgs, nil
}

// GetOrganizations returns all organizations in the system
func (d *OrganizationDB) GetOrganizations() ([]types.Organization, error) {
	var orgs []types.Organization
	err := d.db.Find(&orgs).Error
	return orgs, err
}

// GetOrganizationByID returns an organization by its ID
func (d *OrganizationDB) GetOrganizationByID(id string) (*types.Organization, error) {
	var org types.Organization
	err := d.db.First(&org, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &org, nil
}

// UpdateOrganization updates an existing organization
func (d *OrganizationDB) UpdateOrganization(org *types.Organization) error {
	return d.db.Save(org).Error
}

// DeleteOrganization deletes an organization and its relationships
func (d *OrganizationDB) DeleteOrganization(id string) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		// Delete user-organization relationships
		if err := tx.Exec("DELETE FROM organization_members WHERE organization_id = ?", id).Error; err != nil {
			return err
		}

		// Delete the organization
		if err := tx.Delete(&types.Organization{}, "id = ?", id).Error; err != nil {
			return err
		}

		return nil
	})
}
