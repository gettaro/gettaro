package organization

type CreateOrganizationRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

type AddOrganizationMemberRequest struct {
	Email            string `json:"email" binding:"required"`
	Username         string `json:"username" binding:"required"`
	TitleID          string `json:"title_id" binding:"required"`
	ExternalAccountID string `json:"external_account_id" binding:"required"` // Renamed from SourceControlAccountID
}

type RemoveOrganizationMemberRequest struct {
	UserID string `json:"user_id" binding:"required"`
}
