package database

import (
	"strings"

	orgtypes "ems.dev/backend/services/organization/types"
	"ems.dev/backend/services/user/types"
	"gorm.io/gorm"
)

// DB defines the interface for user database operations
type DB interface {
	FindUser(params types.UserSearchParams) (*types.User, error)
	CreateOrganizationWithOwner(org *orgtypes.Organization, userID string) error
	GetUserOrganizations(userID string) ([]orgtypes.Organization, error)
	CreateUser(user *types.User) (*types.User, error)
	UpdateUser(user *types.User) error
	DeleteUser(userID string) error
	GetUserByID(userID string) (*types.User, error)
	GetUserByEmail(email string) (*types.User, error)
	ListUsers() ([]types.User, error)
}

type UserDB struct {
	db *gorm.DB
}

func NewUserDB(db *gorm.DB) *UserDB {
	return &UserDB{
		db: db,
	}
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
			"INSERT INTO organization_members (user_id, organization_id, is_owner) VALUES (?, ?, true)",
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
		SELECT o.*, om.is_owner
		FROM organizations o
		JOIN organization_members om ON o.id = om.organization_id
		WHERE om.user_id = ?
	`, userID).Scan(&orgs).Error

	if err != nil {
		return nil, err
	}

	return orgs, nil
}

// CreateUser creates a new user in the database and returns the created user
func (d *UserDB) CreateUser(user *types.User) (*types.User, error) {
	err := d.db.Create(user).Error
	// If error is a duplicate key error, return nil
	if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		existingUser, err := d.GetUserByEmail(user.Email)
		if err != nil {
			return nil, err
		}
		return existingUser, nil
	}

	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser updates an existing user in the database
func (d *UserDB) UpdateUser(user *types.User) error {
	return d.db.Save(user).Error
}

// DeleteUser deletes a user from the database
func (d *UserDB) DeleteUser(userID string) error {
	return d.db.Delete(&types.User{}, "id = ?", userID).Error
}

// GetUserByID retrieves a user by their ID
func (d *UserDB) GetUserByID(userID string) (*types.User, error) {
	var user types.User
	err := d.db.First(&user, "id = ?", userID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by their email
func (d *UserDB) GetUserByEmail(email string) (*types.User, error) {
	var user types.User
	err := d.db.First(&user, "email = ?", email).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// ListUsers retrieves all users from the database
func (d *UserDB) ListUsers() ([]types.User, error) {
	var users []types.User
	err := d.db.Find(&users).Error
	return users, err
}
