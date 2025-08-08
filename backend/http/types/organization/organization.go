package organization

type CreateOrganizationRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

type AddOrganizationMemberRequest struct {
	Email    string `json:"email" binding:"required"`
	Username string `json:"username" binding:"required"`
}

type RemoveOrganizationMemberRequest struct {
	UserID string `json:"userId" binding:"required"`
}
