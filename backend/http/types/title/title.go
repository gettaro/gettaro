package title

// CreateTitleRequest represents the request to create a new title
type CreateTitleRequest struct {
	Name      string `json:"name" binding:"required"`
	IsManager bool   `json:"is_manager"`
}

// UpdateTitleRequest represents the request to update a title
type UpdateTitleRequest struct {
	Name      string `json:"name" binding:"required"`
	IsManager bool   `json:"is_manager"`
}

// AssignMemberTitleRequest represents the request to assign a title to a member
type AssignMemberTitleRequest struct {
	MemberID string `json:"member_id" binding:"required"`
	TitleID  string `json:"title_id" binding:"required"`
}
