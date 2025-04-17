package orgtypes

// AddOrganizationMemberRequest represents a request to add a user as a member to an organization
type AddOrganizationMemberRequest struct {
	UserID string `json:"user_id" binding:"required"`
}
