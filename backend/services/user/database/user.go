package database

import (
	"ems.dev/backend/database"
	"ems.dev/backend/services/user/types"
	"gorm.io/gorm"
)

type UserDB struct {
	db *gorm.DB
}

func NewUserDB() *UserDB {
	return &UserDB{
		db: database.DB,
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

// GetUserByEmail finds a user by their email address
func (d *UserDB) GetUserByEmail(email string) (*types.User, error) {
	var user types.User
	err := d.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
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
