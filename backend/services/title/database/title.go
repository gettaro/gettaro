package database

import (
	"ems.dev/backend/services/title/types"
	"gorm.io/gorm"
)

// TitleDB defines the interface for title database operations
type TitleDB interface {
	CreateTitle(title *types.Title) error
	GetTitle(id string) (*types.Title, error)
	ListTitles(orgID string) ([]types.Title, error)
	UpdateTitle(title types.Title) error
	DeleteTitle(id string) error
	AssignUserTitle(userTitle types.UserTitle) error
	GetUserTitle(userID string, orgID string) (*types.UserTitle, error)
	RemoveUserTitle(userID string, orgID string) error
}

type TitleDBImpl struct {
	db *gorm.DB
}

func NewTitleDB(db *gorm.DB) *TitleDBImpl {
	return &TitleDBImpl{
		db: db,
	}
}

// CreateTitle creates a new title
func (d *TitleDBImpl) CreateTitle(title *types.Title) error {
	return d.db.Create(title).Error
}

// GetTitle retrieves a title by ID
func (d *TitleDBImpl) GetTitle(id string) (*types.Title, error) {
	var title types.Title
	err := d.db.First(&title, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &title, nil
}

// ListTitles retrieves all titles for an organization
func (d *TitleDBImpl) ListTitles(orgID string) ([]types.Title, error) {
	var titles []types.Title
	err := d.db.Where("organization_id = ?", orgID).Find(&titles).Error
	return titles, err
}

// UpdateTitle updates an existing title
func (d *TitleDBImpl) UpdateTitle(title types.Title) error {
	return d.db.Save(title).Error
}

// DeleteTitle deletes a title
func (d *TitleDBImpl) DeleteTitle(id string) error {
	return d.db.Delete(&types.Title{}, "id = ?", id).Error
}

// AssignUserTitle assigns a title to a user
func (d *TitleDBImpl) AssignUserTitle(userTitle types.UserTitle) error {
	return d.db.Create(&userTitle).Error
}

// GetUserTitle retrieves a user's title assignment
func (d *TitleDBImpl) GetUserTitle(userID string, orgID string) (*types.UserTitle, error) {
	var userTitle types.UserTitle
	err := d.db.First(&userTitle, "user_id = ? AND organization_id = ?", userID, orgID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &userTitle, nil
}

// RemoveUserTitle removes a user's title assignment
func (d *TitleDBImpl) RemoveUserTitle(userID string, orgID string) error {
	return d.db.Delete(&types.UserTitle{}, "user_id = ? AND organization_id = ?", userID, orgID).Error
}
