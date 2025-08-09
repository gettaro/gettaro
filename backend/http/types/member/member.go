package member

// AddOrganizationMemberRequest represents the request to add a member to an organization
type AddOrganizationMemberRequest struct {
	Email                  string `json:"email" binding:"required"`
	Username               string `json:"username" binding:"required"`
	TitleID                string `json:"titleId" binding:"required"`
	SourceControlAccountID string `json:"sourceControlAccountId" binding:"required"`
}

// UpdateOrganizationMemberRequest represents the request to update a member in an organization
type UpdateOrganizationMemberRequest struct {
	Username               string `json:"username" binding:"required"`
	TitleID                string `json:"titleId" binding:"required"`
	SourceControlAccountID string `json:"sourceControlAccountId" binding:"required"`
}

// RemoveOrganizationMemberRequest represents the request to remove a member from an organization
type RemoveOrganizationMemberRequest struct {
	UserID string `json:"userId" binding:"required"`
}
