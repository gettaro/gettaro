package member

// AddOrganizationMemberRequest represents the request to add a member to an organization
type AddOrganizationMemberRequest struct {
	Email                  string  `json:"email" binding:"required"`
	Username               string  `json:"username" binding:"required"`
	TitleID                string  `json:"title_id" binding:"required"`
	SourceControlAccountID string  `json:"source_control_account_id" binding:"required"`
	ManagerID              *string `json:"manager_id,omitempty"`
}

// UpdateOrganizationMemberRequest represents the request to update a member in an organization
type UpdateOrganizationMemberRequest struct {
	Username               string  `json:"username" binding:"required"`
	TitleID                string  `json:"title_id" binding:"required"`
	SourceControlAccountID string  `json:"source_control_account_id" binding:"required"`
	ManagerID              *string `json:"manager_id,omitempty"`
}

// RemoveOrganizationMemberRequest represents the request to remove a member from an organization
type RemoveOrganizationMemberRequest struct {
	UserID string `json:"user_id" binding:"required"`
}
