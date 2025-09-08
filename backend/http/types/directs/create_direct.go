package directs

import "time"

// CreateDirectReportRequest represents the request body for creating a direct report
type CreateDirectReportRequest struct {
	ManagerID string `json:"managerId" binding:"required"`
	ReportID  string `json:"reportId" binding:"required"`
	Depth     int    `json:"depth" binding:"required"`
}

// CreateDirectReportResponse represents the response for creating a direct report
type CreateDirectReportResponse struct {
	DirectReport DirectReportResponse `json:"directReport"`
}

// DirectReportResponse represents a direct report in API responses
type DirectReportResponse struct {
	ID             string    `json:"id"`
	ManagerID      string    `json:"managerId"`
	ReportID       string    `json:"reportId"`
	OrganizationID string    `json:"organizationId"`
	Depth          int       `json:"depth"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`

	// Populated fields
	Manager *MemberResponse `json:"manager,omitempty"`
	Report  *MemberResponse `json:"report,omitempty"`
}

// MemberResponse represents a user in API responses
type MemberResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	TitleID   string    `json:"titleId"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
