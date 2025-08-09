package types

import (
	"time"
)

// OrganizationMember represents a user's membership in an organization (stored in organization_members table)
type OrganizationMember struct {
	ID             string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID         string    `json:"userId"`
	OrganizationID string    `json:"organizationId"`
	IsOwner        bool      `json:"isOwner"`
	Email          string    `json:"email"`
	Username       string    `json:"username"`
	CreatedAt      time.Time `json:"createdAt" gorm:"default:now()"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// AddMemberRequest represents the request to add a member to an organization
type AddMemberRequest struct {
	Email                  string `json:"email" binding:"required"`
	Username               string `json:"username" binding:"required"`
	TitleID                string `json:"titleId" binding:"required"`
	SourceControlAccountID string `json:"sourceControlAccountId" binding:"required"`
}

// RemoveMemberRequest represents the request to remove a member from an organization
type RemoveMemberRequest struct {
	UserID string `json:"userId" binding:"required"`
}
