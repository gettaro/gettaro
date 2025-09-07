package directs

// GetManagerDirectReportsResponse represents the response for getting manager's direct reports
type GetManagerDirectReportsResponse struct {
	DirectReports []DirectReportResponse `json:"directReports"`
}

// GetManagerTreeResponse represents the response for getting manager's tree
type GetManagerTreeResponse struct {
	OrgChart []OrgChartNodeResponse `json:"orgChart"`
}

// GetMemberManagerResponse represents the response for getting member's manager
type GetMemberManagerResponse struct {
	Manager *DirectReportResponse `json:"manager,omitempty"`
}

// GetMemberManagementChainResponse represents the response for getting member's management chain
type GetMemberManagementChainResponse struct {
	ManagementChain []ManagementChainResponse `json:"managementChain"`
}

// GetOrgChartResponse represents the response for getting organizational chart
type GetOrgChartResponse struct {
	OrgChart []OrgChartNodeResponse `json:"orgChart"`
}

// GetOrgChartFlatResponse represents the response for getting flat organizational chart
type GetOrgChartFlatResponse struct {
	DirectReports []DirectReportResponse `json:"directReports"`
}

// OrgChartNodeResponse represents a node in the organizational chart
type OrgChartNodeResponse struct {
	Member        UserResponse           `json:"member"`
	DirectReports []OrgChartNodeResponse `json:"directReports,omitempty"`
	Depth         int                    `json:"depth"`
}

// ManagementChainResponse represents the management chain for a member
type ManagementChainResponse struct {
	Member  UserResponse  `json:"member"`
	Manager *UserResponse `json:"manager,omitempty"`
	Depth   int           `json:"depth"`
}
