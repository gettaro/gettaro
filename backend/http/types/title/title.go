package title

// CreateTitleRequest represents the request to create a new title
type CreateTitleRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateTitleRequest represents the request to update a title
type UpdateTitleRequest struct {
	Name string `json:"name" binding:"required"`
}

// AssignMemberTitleRequest represents the request to assign a title to a member
type AssignMemberTitleRequest struct {
	MemberID string `json:"memberId" binding:"required"`
	TitleID  string `json:"titleId" binding:"required"`
}
