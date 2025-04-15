package database

import (
	"ems.dev/backend/database"
	"ems.dev/backend/services/organization/types"
	"gorm.io/gorm"
)

type OrganizationDB struct {
	db *gorm.DB
}

func NewOrganizationDB() *OrganizationDB {
	return &OrganizationDB{
		db: database.DB,
	}
}

// CreateOrganization creates a new organization and sets the specified user as its owner
func (d *OrganizationDB) CreateOrganization(org *types.Organization, userID string) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		// Create organization
		if err := tx.Create(org).Error; err != nil {
			return err
		}

		// Create user-organization relationship with owner flag
		if err := tx.Exec(
			"INSERT INTO user_organizations (user_id, organization_id, is_owner) VALUES (?, ?, true)",
			userID,
			org.ID,
		).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetUserOrganizations returns all organizations a user is part of, with ownership information
func (d *OrganizationDB) GetUserOrganizations(userID string) ([]types.Organization, error) {
	type result struct {
		types.Organization
		IsOwner bool `gorm:"column:is_owner"`
	}

	var results []result
	err := d.db.Raw(`
		SELECT o.*, uo.is_owner
		FROM organizations o
		JOIN user_organizations uo ON o.id = uo.organization_id
		WHERE uo.user_id = ?
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
		if err := tx.Exec("DELETE FROM user_organizations WHERE organization_id = ?", id).Error; err != nil {
			return err
		}

		// Delete the organization
		if err := tx.Delete(&types.Organization{}, "id = ?", id).Error; err != nil {
			return err
		}

		return nil
	})
}
