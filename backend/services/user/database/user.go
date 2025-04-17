package database

import (
	orgtypes "ems.dev/backend/services/organization/types"
	"ems.dev/backend/services/user/types"
	"gorm.io/gorm"
)

type UserDB struct {
	db *gorm.DB
}

func NewUserDB(db *gorm.DB) *UserDB {
	return &UserDB{
		db: db,
	}
}

// GetOrCreateUserFromAuthProvider checks if a user exists for the given auth provider and ID,
// and if not, creates a new user and auth provider entry
func (d *UserDB) GetOrCreateUserFromAuthProvider(provider string, providerID string, email string, name string) (*types.User, error) {
	var authProvider types.AuthProvider
	err := d.db.Where("provider = ? AND provider_id = ?", provider, providerID).First(&authProvider).Error
	if err == nil {
		// Auth provider exists, return the associated user
		var user types.User
		err = d.db.First(&user, "id = ?", authProvider.UserID).Error
		if err != nil {
			return nil, err
		}
		return &user, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Create new user
	user := types.User{
		Email:    email,
		Name:     name,
		IsActive: true,
		Status:   "active",
	}

	err = d.db.Create(&user).Error
	if err != nil {
		return nil, err
	}

	// Create auth provider
	authProvider = types.AuthProvider{
		UserID:     user.ID,
		Provider:   provider,
		ProviderID: providerID,
	}

	err = d.db.Create(&authProvider).Error
	if err != nil {
		// If auth provider creation fails, delete the user
		d.db.Delete(&user)
		return nil, err
	}

	return &user, nil
}

// FindUser searches for a user by ID or email
func (d *UserDB) FindUser(params types.UserSearchParams) (*types.User, error) {
	var user types.User
	query := d.db.Model(&types.User{})

	if params.ID != nil {
		query = query.Where("id = ?", *params.ID)
	}
	if params.Email != nil {
		query = query.Where("email = ?", *params.Email)
	}

	err := query.First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// CreateOrganizationWithOwner creates a new organization and sets the specified user as its owner
func (d *UserDB) CreateOrganizationWithOwner(org *orgtypes.Organization, userID string) error {
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
func (d *UserDB) GetUserOrganizations(userID string) ([]orgtypes.Organization, error) {
	var orgs []orgtypes.Organization
	err := d.db.Raw(`
		SELECT o.*, uo.is_owner
		FROM organizations o
		JOIN user_organizations uo ON o.id = uo.organization_id
		WHERE uo.user_id = ?
	`, userID).Scan(&orgs).Error

	if err != nil {
		return nil, err
	}

	return orgs, nil
}
