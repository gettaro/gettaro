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
	AssignMemberTitle(memberTitle types.MemberTitle) error
	GetMemberTitle(memberID string, orgID string) (*types.MemberTitle, error)
	RemoveMemberTitle(memberID string, orgID string) error
}

type DB struct {
	db *gorm.DB
}

func NewTitleDB(db *gorm.DB) *DB {
	return &DB{
		db: db,
	}
}

// CreateTitle creates a new title
func (d *DB) CreateTitle(title *types.Title) error {
	return d.db.Create(title).Error
}

// GetTitle retrieves a title by ID
func (d *DB) GetTitle(id string) (*types.Title, error) {
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
func (d *DB) ListTitles(orgID string) ([]types.Title, error) {
	var titles []types.Title
	err := d.db.Where("organization_id = ?", orgID).Find(&titles).Error
	return titles, err
}

// UpdateTitle updates an existing title
func (d *DB) UpdateTitle(title types.Title) error {
	return d.db.Save(title).Error
}

// DeleteTitle deletes a title
func (d *DB) DeleteTitle(id string) error {
	return d.db.Delete(&types.Title{}, "id = ?", id).Error
}

// AssignMemberTitle assigns a title to a member
func (d *DB) AssignMemberTitle(memberTitle types.MemberTitle) error {
	return d.db.Create(&memberTitle).Error
}

// GetMemberTitle retrieves a member's title assignment
func (d *DB) GetMemberTitle(memberID string, orgID string) (*types.MemberTitle, error) {
	var memberTitle types.MemberTitle
	err := d.db.First(&memberTitle, "member_id = ? AND organization_id = ?", memberID, orgID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &memberTitle, nil
}

// RemoveMemberTitle removes a member's title assignment
func (d *DB) RemoveMemberTitle(memberID string, orgID string) error {
	return d.db.Delete(&types.MemberTitle{}, "member_id = ? AND organization_id = ?", memberID, orgID).Error
}
