package types

import (
	"time"
)

// OrganizationMember represents a user's membership in an organization
// This is stored in the organization_members table
type OrganizationMember struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Email          string    `json:"email"`
	Username       string    `json:"username"`
	OrganizationID string    `json:"organization_id"`
	IsOwner        bool      `json:"is_owner"`
	TitleID        *string   `json:"title_id,omitempty"`   // New field for direct title reference
	ManagerID      *string   `json:"manager_id,omitempty"` // Manager's member ID
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type OrganizationMemberParams struct {
	TitleIDs []string `json:"title_ids"`
	IDs      []string `json:"ids"`
}

// AddMemberRequest represents the request to add a member to an organization
type AddMemberRequest struct {
	Email              string  `json:"email"`
	Username           string  `json:"username"`
	TitleID            string  `json:"title_id"`
	ExternalAccountID string  `json:"external_account_id"` // Renamed from SourceControlAccountID
	ManagerID          *string `json:"manager_id,omitempty"`
}

// UpdateMemberRequest represents the request to update a member in an organization
type UpdateMemberRequest struct {
	Username           string  `json:"username"`
	TitleID            string  `json:"title_id"`
	ExternalAccountID  string  `json:"external_account_id"` // Renamed from SourceControlAccountID
	ManagerID          *string `json:"manager_id,omitempty"`
}

// RemoveMemberRequest represents the request to remove a member from an organization
type RemoveMemberRequest struct {
	UserID string `json:"user_id"`
}
