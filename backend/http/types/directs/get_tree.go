package directs

// GetManagerDirectReportsResponse represents the response for getting manager's direct reports
type GetManagerDirectReportsResponse struct {
	DirectReports []DirectReportResponse `json:"direct_reports"`
}

// GetManagerTreeResponse represents the response for getting manager's tree
type GetManagerTreeResponse struct {
	OrgChart []OrgChartNodeResponse `json:"org_chart"`
}

// GetMemberManagerResponse represents the response for getting member's manager
type GetMemberManagerResponse struct {
	Manager *DirectReportResponse `json:"manager,omitempty"`
}

// GetMemberManagementChainResponse represents the response for getting member's management chain
type GetMemberManagementChainResponse struct {
	ManagementChain []ManagementChainResponse `json:"management_chain"`
}

// GetOrgChartResponse represents the response for getting organizational chart
type GetOrgChartResponse struct {
	OrgChart []OrgChartNodeResponse `json:"org_chart"`
}

// GetOrgChartFlatResponse represents the response for getting flat organizational chart
type GetOrgChartFlatResponse struct {
	DirectReports []DirectReportResponse `json:"direct_reports"`
}

// OrgChartNodeResponse represents a node in the organizational chart
type OrgChartNodeResponse struct {
	Member        MemberResponse         `json:"member"`
	DirectReports []OrgChartNodeResponse `json:"direct_reports,omitempty"`
	Depth         int                    `json:"depth"`
}

// ManagementChainResponse represents the management chain for a member
type ManagementChainResponse struct {
	Member  MemberResponse  `json:"member"`
	Manager *MemberResponse `json:"manager,omitempty"`
	Depth   int             `json:"depth"`
}
