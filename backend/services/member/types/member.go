package types

import (
	"time"
)

// OrganizationMember represents a user's membership in an organization
// This is stored in the organization_members table
type OrganizationMember struct {
	ID             string    `json:"id"`
	UserID         string    `json:"userId"`
	Email          string    `json:"email"`
	Username       string    `json:"username"`
	OrganizationID string    `json:"organizationId"`
	IsOwner        bool      `json:"isOwner"`
	TitleID        *string   `json:"titleId,omitempty"`   // New field for direct title reference
	ManagerID      *string   `json:"managerId,omitempty"` // Manager's member ID
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type OrganizationMemberParams struct {
	TitleIDs []string `json:"titleIds"`
	IDs      []string `json:"ids"`
}

// AddMemberRequest represents the request to add a member to an organization
type AddMemberRequest struct {
	Email                  string  `json:"email"`
	Username               string  `json:"username"`
	TitleID                string  `json:"titleId"`
	SourceControlAccountID string  `json:"sourceControlAccountId"`
	ManagerID              *string `json:"managerId,omitempty"`
}

// UpdateMemberRequest represents the request to update a member in an organization
type UpdateMemberRequest struct {
	Username               string  `json:"username"`
	TitleID                string  `json:"titleId"`
	SourceControlAccountID string  `json:"sourceControlAccountId"`
	ManagerID              *string `json:"managerId,omitempty"`
}

// RemoveMemberRequest represents the request to remove a member from an organization
type RemoveMemberRequest struct {
	UserID string `json:"userId"`
}
