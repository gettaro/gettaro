package directs

// ListDirectReportsQuery represents query parameters for listing direct reports
type ListDirectReportsQuery struct {
	ManagerID *string `form:"managerId"`
	ReportID  *string `form:"reportId"`
	Depth     *int    `form:"depth"`
	Limit     *int    `form:"limit"`
	Offset    *int    `form:"offset"`
}

// ListDirectReportsResponse represents the response for listing direct reports
type ListDirectReportsResponse struct {
	DirectReports []DirectReportResponse `json:"directReports"`
}

// GetDirectReportResponse represents the response for getting a single direct report
type GetDirectReportResponse struct {
	DirectReport DirectReportResponse `json:"directReport"`
}

// UpdateDirectReportRequest represents the request body for updating a direct report
type UpdateDirectReportRequest struct {
	Depth *int `json:"depth,omitempty"`
}

// UpdateDirectReportResponse represents the response for updating a direct report
type UpdateDirectReportResponse struct {
	DirectReport DirectReportResponse `json:"directReport"`
}
