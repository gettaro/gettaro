package types

import (
	"time"

	membertypes "ems.dev/backend/services/member/types"
	orgtypes "ems.dev/backend/services/organization/types"
)

// TeamType represents the type of team
type TeamType string

const (
	TeamTypeSquad   TeamType = "squad"
	TeamTypeChapter TeamType = "chapter"
	TeamTypeTribe   TeamType = "tribe"
	TeamTypeGuild   TeamType = "guild"
)

// Team represents a team in the system
type Team struct {
	ID             string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Type           *TeamType `json:"type" gorm:"type:varchar(50);check:type IN ('squad', 'chapter', 'tribe', 'guild')"`
	PRPrefix       *string   `json:"pr_prefix" gorm:"type:varchar(255)"`
	OrganizationID string    `json:"organization_id" gorm:"type:uuid"`
	CreatedAt      time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Relationships
	Organization orgtypes.Organization `json:"organization" gorm:"foreignKey:OrganizationID"`
	Members      []TeamMember          `json:"members" gorm:"foreignKey:TeamID"`
}

// TeamMember represents a member's membership in a team
type TeamMember struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TeamID    string    `json:"team_id" gorm:"type:uuid"`
	MemberID  string    `json:"member_id" gorm:"type:uuid"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Team   Team                           `json:"team" gorm:"foreignKey:TeamID"`
	Member membertypes.OrganizationMember `json:"member" gorm:"foreignKey:MemberID"`
}

// TeamSearchParams represents parameters for searching teams
type TeamSearchParams struct {
	ID             *string `json:"id"`
	Name           *string `json:"name"`
	OrganizationID *string `json:"organization_id"`
}

// CreateTeamRequest represents the request body for creating a team
type CreateTeamRequest struct {
	Name           string    `json:"name" binding:"required"`
	Description    string    `json:"description"`
	Type           *TeamType  `json:"type" binding:"omitempty,oneof=squad chapter tribe guild"`
	PRPrefix       *string   `json:"pr_prefix"`
	OrganizationID string    `json:"organization_id" binding:"required"`
}

// UpdateTeamRequest represents the request body for updating a team
type UpdateTeamRequest struct {
	Name        *string    `json:"name"`
	Description *string   `json:"description"`
	Type        *TeamType  `json:"type" binding:"omitempty,oneof=squad chapter tribe guild"`
	PRPrefix    *string   `json:"pr_prefix"`
}

// AddTeamMemberRequest represents the request body for adding a member to a team
type AddTeamMemberRequest struct {
	MemberID string `json:"member_id" binding:"required"`
}

// RemoveTeamMemberRequest represents the request body for removing a member from a team
type RemoveTeamMemberRequest struct {
	MemberID string `json:"member_id" binding:"required"`
}
