package directs

import "time"

// CreateDirectReportRequest represents the request body for creating a direct report
type CreateDirectReportRequest struct {
	ManagerID string `json:"manager_id" binding:"required"`
	ReportID  string `json:"report_id" binding:"required"`
	Depth     int    `json:"depth" binding:"required"`
}

// CreateDirectReportResponse represents the response for creating a direct report
type CreateDirectReportResponse struct {
	DirectReport DirectReportResponse `json:"direct_report"`
}

// DirectReportResponse represents a direct report in API responses
type DirectReportResponse struct {
	ID             string    `json:"id"`
	ManagerID      string    `json:"manager_id"`
	ReportID       string    `json:"report_id"`
	OrganizationID string    `json:"organization_id"`
	Depth          int       `json:"depth"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Populated fields
	Manager *MemberResponse `json:"manager,omitempty"`
	Report  *MemberResponse `json:"report,omitempty"`
}

// MemberResponse represents a user in API responses
type MemberResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	TitleID   string    `json:"title_id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
