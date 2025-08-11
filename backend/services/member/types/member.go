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
	TitleID        *string   `json:"titleId,omitempty"` // New field for direct title reference
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// AddMemberRequest represents the request to add a member to an organization
type AddMemberRequest struct {
	Email                  string `json:"email"`
	Username               string `json:"username"`
	TitleID                string `json:"titleId"`
	SourceControlAccountID string `json:"sourceControlAccountId"`
}

// RemoveMemberRequest represents the request to remove a member from an organization
type RemoveMemberRequest struct {
	UserID string `json:"userId"`
}
