package api

import (
	"context"

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
