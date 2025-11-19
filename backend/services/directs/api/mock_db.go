package api

import (
	"context"

	"ems.dev/backend/services/directs/database"
	"ems.dev/backend/services/directs/types"
	"github.com/stretchr/testify/mock"
)

// MockDB is a mock implementation of the database.DB interface
type MockDB struct {
	mock.Mock
}

func (m *MockDB) CreateDirectReport(ctx context.Context, directReport *types.DirectReport) error {
	args := m.Called(ctx, directReport)
	return args.Error(0)
}

func (m *MockDB) GetDirectReport(ctx context.Context, id string) (*types.DirectReport, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.DirectReport), args.Error(1)
}

func (m *MockDB) ListDirectReports(ctx context.Context, params types.DirectReportSearchParams) ([]types.DirectReport, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.DirectReport), args.Error(1)
}

func (m *MockDB) UpdateDirectReport(ctx context.Context, id string, params types.UpdateDirectReportParams) error {
	args := m.Called(ctx, id, params)
	return args.Error(0)
}

func (m *MockDB) DeleteDirectReport(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDB) GetManagerDirectReports(ctx context.Context, managerID, orgID string) ([]types.DirectReport, error) {
	args := m.Called(ctx, managerID, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.DirectReport), args.Error(1)
}

func (m *MockDB) GetManagerTree(ctx context.Context, managerID, orgID string) ([]types.OrgChartNode, error) {
	args := m.Called(ctx, managerID, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.OrgChartNode), args.Error(1)
}

func (m *MockDB) GetMemberManager(ctx context.Context, reportID, orgID string) (*types.DirectReport, error) {
	args := m.Called(ctx, reportID, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.DirectReport), args.Error(1)
}

func (m *MockDB) GetMemberManagementChain(ctx context.Context, reportID, orgID string) ([]types.ManagementChain, error) {
	args := m.Called(ctx, reportID, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.ManagementChain), args.Error(1)
}

func (m *MockDB) GetOrgChart(ctx context.Context, orgID string) ([]types.OrgChartNode, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.OrgChartNode), args.Error(1)
}

func (m *MockDB) GetOrgChartFlat(ctx context.Context, orgID string) ([]types.DirectReport, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.DirectReport), args.Error(1)
}

// Ensure MockDB implements database.DB interface
var _ database.DB = (*MockDB)(nil)
