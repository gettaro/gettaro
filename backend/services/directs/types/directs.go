package types

import (
	"time"
	
	membertypes "ems.dev/backend/services/member/types"
)

// DirectReport represents a manager-direct report relationship
type DirectReport struct {
	ID                string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ManagerMemberID   string    `json:"managerMemberId" gorm:"column:manager_member_id;not null"`
	ReportMemberID    string    `json:"reportMemberId" gorm:"column:report_member_id;not null"`
	OrganizationID    string    `json:"organizationId" gorm:"not null"`
	Depth             int       `json:"depth" gorm:"not null"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`

	// Relations
	Manager membertypes.OrganizationMember `json:"manager,omitempty" gorm:"foreignKey:ManagerMemberID;references:ID"`
	Report  membertypes.OrganizationMember `json:"report,omitempty" gorm:"foreignKey:ReportMemberID;references:ID"`
}

// DirectReportSearchParams represents search parameters for direct reports
type DirectReportSearchParams struct {
	ID                *string `json:"id,omitempty"`
	ManagerMemberID   *string `json:"managerMemberId,omitempty"`
	ReportMemberID    *string `json:"reportMemberId,omitempty"`
	OrganizationID    *string `json:"organizationId,omitempty"`
	Depth             *int    `json:"depth,omitempty"`
}

// CreateDirectReportParams represents parameters for creating a direct report
type CreateDirectReportParams struct {
	ManagerMemberID   string `json:"managerMemberId" binding:"required"`
	ReportMemberID    string `json:"reportMemberId" binding:"required"`
	OrganizationID    string `json:"organizationId" binding:"required"`
	Depth             int    `json:"depth" binding:"required"`
}

// UpdateDirectReportParams represents parameters for updating a direct report
type UpdateDirectReportParams struct {
	Depth *int `json:"depth,omitempty"`
}

// OrgChartNode represents a node in the organizational chart
type OrgChartNode struct {
	Member        membertypes.OrganizationMember `json:"member"`
	DirectReports []OrgChartNode                 `json:"directReports,omitempty"`
	Depth         int                            `json:"depth"`
}

// ManagementChain represents the management chain for a member
type ManagementChain struct {
	Member  membertypes.OrganizationMember  `json:"member"`
	Manager *membertypes.OrganizationMember `json:"manager,omitempty"`
	Depth   int                             `json:"depth"`
}

