package database

import (
	"context"
	"errors"

	"ems.dev/backend/services/directs/types"
	"gorm.io/gorm"
)

// DB defines the interface for direct reports database operations
type DB interface {
	CreateDirectReport(ctx context.Context, directReport *types.DirectReport) error
	GetDirectReport(ctx context.Context, id string) (*types.DirectReport, error)
	ListDirectReports(ctx context.Context, params types.DirectReportSearchParams) ([]types.DirectReport, error)
	UpdateDirectReport(ctx context.Context, id string, params types.UpdateDirectReportParams) error
	DeleteDirectReport(ctx context.Context, id string) error

	// Manager operations
	GetManagerDirectReports(ctx context.Context, managerID, orgID string) ([]types.DirectReport, error)
	GetManagerTree(ctx context.Context, managerID, orgID string) ([]types.OrgChartNode, error)

	// Employee operations
	GetMemberManager(ctx context.Context, reportID, orgID string) (*types.DirectReport, error)
	GetMemberManagementChain(ctx context.Context, reportID, orgID string) ([]types.ManagementChain, error)

	// Organizational structure
	GetOrgChart(ctx context.Context, orgID string) ([]types.OrgChartNode, error)
	GetOrgChartFlat(ctx context.Context, orgID string) ([]types.DirectReport, error)
}

type DirectReportsDB struct {
	db *gorm.DB
}

func NewDirectReportsDB(db *gorm.DB) *DirectReportsDB {
	return &DirectReportsDB{
		db: db,
	}
}

// CreateDirectReport creates a new direct report relationship
func (d *DirectReportsDB) CreateDirectReport(ctx context.Context, directReport *types.DirectReport) error {
	return d.db.WithContext(ctx).Create(directReport).Error
}

// GetDirectReport retrieves a direct report by ID
func (d *DirectReportsDB) GetDirectReport(ctx context.Context, id string) (*types.DirectReport, error) {
	var directReport types.DirectReport
	err := d.db.WithContext(ctx).
		Preload("Manager").
		Preload("Report").
		First(&directReport, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &directReport, nil
}

// ListDirectReports retrieves direct reports based on search parameters
func (d *DirectReportsDB) ListDirectReports(ctx context.Context, params types.DirectReportSearchParams) ([]types.DirectReport, error) {
	var directReports []types.DirectReport
	query := d.db.WithContext(ctx).
		Preload("Manager").
		Preload("Report")

	if params.ID != nil {
		query = query.Where("id = ?", *params.ID)
	}
	if params.ManagerMemberID != nil {
		query = query.Where("manager_member_id = ?", *params.ManagerMemberID)
	}
	if params.ReportMemberID != nil {
		query = query.Where("report_member_id = ?", *params.ReportMemberID)
	}
	if params.OrganizationID != nil {
		query = query.Where("organization_id = ?", *params.OrganizationID)
	}
	if params.Depth != nil {
		query = query.Where("depth = ?", *params.Depth)
	}

	err := query.Find(&directReports).Error
	return directReports, err
}

// UpdateDirectReport updates a direct report relationship
func (d *DirectReportsDB) UpdateDirectReport(ctx context.Context, id string, params types.UpdateDirectReportParams) error {
	updates := make(map[string]interface{})
	if params.Depth != nil {
		updates["depth"] = *params.Depth
	}

	if len(updates) == 0 {
		return nil
	}

	return d.db.WithContext(ctx).
		Model(&types.DirectReport{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// DeleteDirectReport removes a direct report relationship
func (d *DirectReportsDB) DeleteDirectReport(ctx context.Context, id string) error {
	return d.db.WithContext(ctx).Delete(&types.DirectReport{}, "id = ?", id).Error
}

// GetManagerDirectReports retrieves all direct reports for a specific manager
func (d *DirectReportsDB) GetManagerDirectReports(ctx context.Context, managerMemberID, orgID string) ([]types.DirectReport, error) {
	var directReports []types.DirectReport
	err := d.db.WithContext(ctx).
		Preload("Manager").
		Preload("Report").
		Where("manager_member_id = ? AND organization_id = ?", managerMemberID, orgID).
		Order("depth ASC").
		Find(&directReports).Error
	return directReports, err
}

// GetManagerTree retrieves the full management tree for a manager
func (d *DirectReportsDB) GetManagerTree(ctx context.Context, managerMemberID, orgID string) ([]types.OrgChartNode, error) {
	// Get all direct reports recursively
	return d.buildOrgChartNode(ctx, managerMemberID, orgID, 0)
}

// GetMemberManager retrieves the manager of a specific member
func (d *DirectReportsDB) GetMemberManager(ctx context.Context, reportMemberID, orgID string) (*types.DirectReport, error) {
	var directReport types.DirectReport
	err := d.db.WithContext(ctx).
		Preload("Manager").
		Preload("Report").
		Where("report_member_id = ? AND organization_id = ?", reportMemberID, orgID).
		First(&directReport).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &directReport, nil
}

// GetMemberManagementChain retrieves the full management chain for a member
func (d *DirectReportsDB) GetMemberManagementChain(ctx context.Context, reportMemberID, orgID string) ([]types.ManagementChain, error) {
	var chain []types.ManagementChain
	currentReportMemberID := reportMemberID
	depth := 0

	for {
		var directReport types.DirectReport
		err := d.db.WithContext(ctx).
			Preload("Manager").
			Preload("Report").
			Where("report_member_id = ? AND organization_id = ?", currentReportMemberID, orgID).
			First(&directReport).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				break
			}
			return nil, err
		}

		chain = append(chain, types.ManagementChain{
			Member:  directReport.Report,
			Manager: &directReport.Manager,
			Depth:   depth,
		})

		currentReportMemberID = directReport.ManagerMemberID
		depth++

		// Prevent infinite loops
		if depth > 10 {
			break
		}
	}

	return chain, nil
}

// GetOrgChart retrieves the complete organizational chart
func (d *DirectReportsDB) GetOrgChart(ctx context.Context, orgID string) ([]types.OrgChartNode, error) {
	// Find all top-level managers (those who are not reports of anyone)
	var topLevelManagers []types.DirectReport
	err := d.db.WithContext(ctx).
		Preload("Manager").
		Preload("Report").
		Where("organization_id = ?", orgID).
		Find(&topLevelManagers).Error
	if err != nil {
		return nil, err
	}

	// Find managers who are not reports
	managerIDs := make(map[string]bool)
	reportIDs := make(map[string]bool)

	for _, dr := range topLevelManagers {
		managerIDs[dr.ManagerMemberID] = true
		reportIDs[dr.ReportMemberID] = true
	}

	// Find top-level managers (managers who are not reports)
	var topManagers []string
	for managerMemberID := range managerIDs {
		if !reportIDs[managerMemberID] {
			topManagers = append(topManagers, managerMemberID)
		}
	}

	// Build org chart for each top-level manager
	var orgChart []types.OrgChartNode
	for _, managerMemberID := range topManagers {
		node, err := d.buildOrgChartNode(ctx, managerMemberID, orgID, 0)
		if err != nil {
			return nil, err
		}
		orgChart = append(orgChart, node...)
	}

	return orgChart, nil
}

// GetOrgChartFlat retrieves a flat list of all manager-direct relationships
func (d *DirectReportsDB) GetOrgChartFlat(ctx context.Context, orgID string) ([]types.DirectReport, error) {
	var directReports []types.DirectReport
	err := d.db.WithContext(ctx).
		Preload("Manager").
		Preload("Report").
		Where("organization_id = ?", orgID).
		Order("depth ASC, manager_member_id ASC").
		Find(&directReports).Error
	return directReports, err
}

// buildOrgChartNode recursively builds org chart nodes
func (d *DirectReportsDB) buildOrgChartNode(ctx context.Context, managerMemberID, orgID string, depth int) ([]types.OrgChartNode, error) {
	// Get direct reports for this manager
	directReports, err := d.GetManagerDirectReports(ctx, managerMemberID, orgID)
	if err != nil {
		return nil, err
	}

	if len(directReports) == 0 {
		return nil, nil
	}

	var nodes []types.OrgChartNode
	for _, dr := range directReports {
		// Recursively get direct reports for this member
		subNodes, err := d.buildOrgChartNode(ctx, dr.ReportMemberID, orgID, depth+1)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, types.OrgChartNode{
			Member:        dr.Report,
			DirectReports: subNodes,
			Depth:         depth,
		})
	}

	return nodes, nil
}
