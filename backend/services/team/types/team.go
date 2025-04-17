package types

import (
	"time"

	orgtypes "ems.dev/backend/services/organization/types"
	usertypes "ems.dev/backend/services/user/types"
)

// Team represents a team in the system
type Team struct {
	ID             string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	OrganizationID string    `json:"organizationId" gorm:"type:uuid"`
	CreatedAt      time.Time `json:"createdAt" gorm:"default:now()"`
	UpdatedAt      time.Time `json:"updatedAt"`

	// Relationships
	Organization orgtypes.Organization `json:"organization" gorm:"foreignKey:OrganizationID"`
	Members      []TeamMember          `json:"members" gorm:"foreignKey:TeamID"`
}

// TeamMember represents a user's membership in a team
type TeamMember struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TeamID    string    `json:"teamId" gorm:"type:uuid"`
	UserID    string    `json:"userId" gorm:"type:uuid"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updatedAt"`

	// Relationships
	Team Team           `json:"team" gorm:"foreignKey:TeamID"`
	User usertypes.User `json:"user" gorm:"foreignKey:UserID"`
}

// TeamSearchParams represents parameters for searching teams
type TeamSearchParams struct {
	ID             *string `json:"id"`
	Name           *string `json:"name"`
	OrganizationID *string `json:"organizationId"`
}

// CreateTeamRequest represents the request body for creating a team
type CreateTeamRequest struct {
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description"`
	OrganizationID string `json:"organizationId" binding:"required"`
}

// UpdateTeamRequest represents the request body for updating a team
type UpdateTeamRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

// AddTeamMemberRequest represents the request body for adding a member to a team
type AddTeamMemberRequest struct {
	UserID string `json:"userId" binding:"required"`
	Role   string `json:"role" binding:"required"`
}

// RemoveTeamMemberRequest represents the request body for removing a member from a team
type RemoveTeamMemberRequest struct {
	UserID string `json:"userId" binding:"required"`
}
