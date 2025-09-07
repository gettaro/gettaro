package api

import (
	"context"

	"ems.dev/backend/libraries/errors"
	"ems.dev/backend/services/directs/database"
	"ems.dev/backend/services/directs/types"
)

// DirectReportsAPI defines the interface for direct reports operations
type DirectReportsAPI interface {
	CreateDirectReport(ctx context.Context, params types.CreateDirectReportParams) (*types.DirectReport, error)
	GetDirectReport(ctx context.Context, id string) (*types.DirectReport, error)
	ListDirectReports(ctx context.Context, params types.DirectReportSearchParams) ([]types.DirectReport, error)
	UpdateDirectReport(ctx context.Context, id string, params types.UpdateDirectReportParams) error
	DeleteDirectReport(ctx context.Context, id string) error

	// Manager operations
	GetManagerDirectReports(ctx context.Context, managerMemberID, orgID string) ([]types.DirectReport, error)
	GetManagerTree(ctx context.Context, managerMemberID, orgID string) ([]types.OrgChartNode, error)

	// Employee operations
	GetMemberManager(ctx context.Context, reportMemberID, orgID string) (*types.DirectReport, error)
	GetMemberManagementChain(ctx context.Context, reportMemberID, orgID string) ([]types.ManagementChain, error)

	// Organizational structure
	GetOrgChart(ctx context.Context, orgID string) ([]types.OrgChartNode, error)
	GetOrgChartFlat(ctx context.Context, orgID string) ([]types.DirectReport, error)
}

type DirectReportsAPIImpl struct {
	db database.DB
}

func NewDirectReportsAPI(db database.DB) DirectReportsAPI {
	return &DirectReportsAPIImpl{
		db: db,
	}
}

// CreateDirectReport creates a new direct report relationship
func (a *DirectReportsAPIImpl) CreateDirectReport(ctx context.Context, params types.CreateDirectReportParams) (*types.DirectReport, error) {
	// Validate that all required fields are not empty
	if params.ManagerMemberID == "" {
		return nil, errors.NewBadRequestError("manager member ID cannot be empty")
	}
	if params.ReportMemberID == "" {
		return nil, errors.NewBadRequestError("report member ID cannot be empty")
	}
	if params.OrganizationID == "" {
		return nil, errors.NewBadRequestError("organization ID cannot be empty")
	}

	directReport := &types.DirectReport{
		ManagerMemberID: params.ManagerMemberID,
		ReportMemberID:  params.ReportMemberID,
		OrganizationID:  params.OrganizationID,
		Depth:           params.Depth,
	}

	err := a.db.CreateDirectReport(ctx, directReport)
	if err != nil {
		return nil, err
	}

	return a.db.GetDirectReport(ctx, directReport.ID)
}

// GetDirectReport retrieves a direct report by ID
func (a *DirectReportsAPIImpl) GetDirectReport(ctx context.Context, id string) (*types.DirectReport, error) {
	return a.db.GetDirectReport(ctx, id)
}

// ListDirectReports retrieves direct reports based on search parameters
func (a *DirectReportsAPIImpl) ListDirectReports(ctx context.Context, params types.DirectReportSearchParams) ([]types.DirectReport, error) {
	return a.db.ListDirectReports(ctx, params)
}

// UpdateDirectReport updates a direct report relationship
func (a *DirectReportsAPIImpl) UpdateDirectReport(ctx context.Context, id string, params types.UpdateDirectReportParams) error {
	return a.db.UpdateDirectReport(ctx, id, params)
}

// DeleteDirectReport removes a direct report relationship
func (a *DirectReportsAPIImpl) DeleteDirectReport(ctx context.Context, id string) error {
	return a.db.DeleteDirectReport(ctx, id)
}

// GetManagerDirectReports retrieves all direct reports for a specific manager
func (a *DirectReportsAPIImpl) GetManagerDirectReports(ctx context.Context, managerMemberID, orgID string) ([]types.DirectReport, error) {
	return a.db.GetManagerDirectReports(ctx, managerMemberID, orgID)
}

// GetManagerTree retrieves the full management tree for a manager
func (a *DirectReportsAPIImpl) GetManagerTree(ctx context.Context, managerMemberID, orgID string) ([]types.OrgChartNode, error) {
	return a.db.GetManagerTree(ctx, managerMemberID, orgID)
}

// GetMemberManager retrieves the manager of a specific member
func (a *DirectReportsAPIImpl) GetMemberManager(ctx context.Context, reportMemberID, orgID string) (*types.DirectReport, error) {
	return a.db.GetMemberManager(ctx, reportMemberID, orgID)
}

// GetMemberManagementChain retrieves the full management chain for a member
func (a *DirectReportsAPIImpl) GetMemberManagementChain(ctx context.Context, reportMemberID, orgID string) ([]types.ManagementChain, error) {
	return a.db.GetMemberManagementChain(ctx, reportMemberID, orgID)
}

// GetOrgChart retrieves the complete organizational chart
func (a *DirectReportsAPIImpl) GetOrgChart(ctx context.Context, orgID string) ([]types.OrgChartNode, error) {
	return a.db.GetOrgChart(ctx, orgID)
}

// GetOrgChartFlat retrieves a flat list of all manager-direct relationships
func (a *DirectReportsAPIImpl) GetOrgChartFlat(ctx context.Context, orgID string) ([]types.DirectReport, error) {
	return a.db.GetOrgChartFlat(ctx, orgID)
}
