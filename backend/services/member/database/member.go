package database

import (
	"ems.dev/backend/services/member/types"
	"gorm.io/gorm"
)

// DB defines the interface for member database operations
type DB interface {
	AddOrganizationMember(member *types.OrganizationMember) error
	RemoveOrganizationMember(orgID string, userID string) error
	GetOrganizationMembers(orgID string) ([]types.OrganizationMember, error)
	GetOrganizationMember(orgID string, userID string) (*types.OrganizationMember, error)
	IsOrganizationOwner(orgID string, userID string) (bool, error)
	UpdateOrganizationMember(orgID string, userID string, username string) error
}

type MemberDB struct {
	db *gorm.DB
}

func NewMemberDB(db *gorm.DB) *MemberDB {
	return &MemberDB{
		db: db,
	}
}

// AddOrganizationMember adds a user as a member to an organization
func (d *MemberDB) AddOrganizationMember(member *types.OrganizationMember) error {
	return d.db.Exec(
		"INSERT INTO organization_members (user_id, organization_id, email, username, is_owner) VALUES (?, ?, ?, ?, false)",
		member.UserID,
		member.OrganizationID,
		member.Email,
		member.Username,
	).Error
}

// RemoveOrganizationMember removes a user from an organization
func (d *MemberDB) RemoveOrganizationMember(orgID string, userID string) error {
	return d.db.Exec(
		"DELETE FROM organization_members WHERE organization_id = ? AND user_id = ? AND is_owner = false",
		orgID,
		userID,
	).Error
}

// GetOrganizationMembers returns all members of an organization
func (d *MemberDB) GetOrganizationMembers(orgID string) ([]types.OrganizationMember, error) {
	var members []types.OrganizationMember
	err := d.db.Raw(`
		SELECT om.*
		FROM organization_members om
		WHERE om.organization_id = ?
	`, orgID).Scan(&members).Error
	return members, err
}

// GetOrganizationMember returns a specific member of an organization
func (d *MemberDB) GetOrganizationMember(orgID string, userID string) (*types.OrganizationMember, error) {
	var member types.OrganizationMember
	err := d.db.First(&member, "organization_id = ? AND user_id = ?", orgID, userID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &member, nil
}

// IsOrganizationOwner checks if a user is the owner of an organization
func (d *MemberDB) IsOrganizationOwner(orgID string, userID string) (bool, error) {
	var count int64
	err := d.db.Raw(`
		SELECT COUNT(*)
		FROM organization_members
		WHERE organization_id = ? AND user_id = ? AND is_owner = true
	`, orgID, userID).Scan(&count).Error
	return count > 0, err
}

// UpdateOrganizationMember updates a member's details in an organization
func (d *MemberDB) UpdateOrganizationMember(orgID string, userID string, username string) error {
	return d.db.Exec(
		"UPDATE organization_members SET username = ?, updated_at = NOW() WHERE organization_id = ? AND user_id = ? AND is_owner = false",
		username,
		orgID,
		userID,
	).Error
}
