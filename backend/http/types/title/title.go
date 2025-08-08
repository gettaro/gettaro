package title

// CreateTitleRequest represents the request to create a new title
type CreateTitleRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateTitleRequest represents the request to update a title
type UpdateTitleRequest struct {
	Name string `json:"name" binding:"required"`
}

// AssignUserTitleRequest represents the request to assign a title to a user
type AssignUserTitleRequest struct {
	UserID  string `json:"userId" binding:"required"`
	TitleID string `json:"titleId" binding:"required"`
}
