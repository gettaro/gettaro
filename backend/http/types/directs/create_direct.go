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
	Manager *UserResponse `json:"manager,omitempty"`
	Report  *UserResponse `json:"report,omitempty"`
}

// UserResponse represents a user in API responses
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"isActive"`
	Status    *string   `json:"status,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
