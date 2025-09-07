package types

import "time"

// DirectReport represents a manager-direct report relationship
type DirectReport struct {
	ID             string    `json:"id" gorm:"primaryKey"`
	ManagerID      string    `json:"managerId" gorm:"not null"`
	ReportID       string    `json:"reportId" gorm:"not null"`
	OrganizationID string    `json:"organizationId" gorm:"not null"`
	Depth          int       `json:"depth" gorm:"not null"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`

	// Relations
	Manager User `json:"manager,omitempty" gorm:"foreignKey:ManagerID;references:ID"`
	Report  User `json:"report,omitempty" gorm:"foreignKey:ReportID;references:ID"`
}

// DirectReportSearchParams represents search parameters for direct reports
type DirectReportSearchParams struct {
	ID             *string `json:"id,omitempty"`
	ManagerID      *string `json:"managerId,omitempty"`
	ReportID       *string `json:"reportId,omitempty"`
	OrganizationID *string `json:"organizationId,omitempty"`
	Depth          *int    `json:"depth,omitempty"`
}

// CreateDirectReportParams represents parameters for creating a direct report
type CreateDirectReportParams struct {
	ManagerID      string `json:"managerId" binding:"required"`
	ReportID       string `json:"reportId" binding:"required"`
	OrganizationID string `json:"organizationId" binding:"required"`
	Depth          int    `json:"depth" binding:"required"`
}

// UpdateDirectReportParams represents parameters for updating a direct report
type UpdateDirectReportParams struct {
	Depth *int `json:"depth,omitempty"`
}

// OrgChartNode represents a node in the organizational chart
type OrgChartNode struct {
	User          User           `json:"user"`
	DirectReports []OrgChartNode `json:"directReports,omitempty"`
	Depth         int            `json:"depth"`
}

// ManagementChain represents the management chain for a user
type ManagementChain struct {
	User    User  `json:"user"`
	Manager *User `json:"manager,omitempty"`
	Depth   int   `json:"depth"`
}

// User represents a user (imported from user service)
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"isActive"`
	Status    *string   `json:"status,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
