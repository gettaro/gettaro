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
