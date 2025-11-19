package api

import (
	"context"
	"time"

	metrictypes "ems.dev/backend/services/sourcecontrol/metrics/types"
	"ems.dev/backend/services/sourcecontrol/types"
	"github.com/stretchr/testify/mock"
)

// MockDB is a mock implementation of the database.DB interface
type MockDB struct {
	mock.Mock
}

func (m *MockDB) GetPullRequests(ctx context.Context, params *types.PullRequestParams) ([]*types.PullRequest, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*types.PullRequest), args.Error(1)
}

func (m *MockDB) CreatePullRequest(ctx context.Context, pr *types.PullRequest) (*types.PullRequest, error) {
	args := m.Called(ctx, pr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.PullRequest), args.Error(1)
}

func (m *MockDB) UpdatePullRequest(ctx context.Context, pr *types.PullRequest) error {
	args := m.Called(ctx, pr)
	return args.Error(0)
}

func (m *MockDB) CreatePRComments(ctx context.Context, comments []*types.PRComment) error {
	args := m.Called(ctx, comments)
	return args.Error(0)
}

func (m *MockDB) GetPullRequestComments(ctx context.Context, prID string) ([]*types.PRComment, error) {
	args := m.Called(ctx, prID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*types.PRComment), args.Error(1)
}

func (m *MockDB) GetMemberPullRequests(ctx context.Context, params *types.MemberPullRequestParams) ([]*types.PullRequestWithComments, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*types.PullRequestWithComments), args.Error(1)
}

func (m *MockDB) GetMemberPullRequestReviews(ctx context.Context, params *types.MemberPullRequestReviewsParams) ([]*types.MemberActivity, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*types.MemberActivity), args.Error(1)
}

func (m *MockDB) CalculateTimeToMerge(ctx context.Context, organizationID string, sourceControlAccountIDs []string, prPrefixes []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation) (*int, error) {
	args := m.Called(ctx, organizationID, sourceControlAccountIDs, prPrefixes, startDate, endDate, metricOperation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int), args.Error(1)
}

func (m *MockDB) CalculateTimeToMergeGraph(ctx context.Context, organizationID string, sourceControlAccountIDs []string, prPrefixes []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation, metricLabel string, interval string) ([]types.TimeSeriesEntry, error) {
	args := m.Called(ctx, organizationID, sourceControlAccountIDs, prPrefixes, startDate, endDate, metricOperation, metricLabel, interval)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.TimeSeriesEntry), args.Error(1)
}

func (m *MockDB) CalculatePRsMerged(ctx context.Context, organizationID string, sourceControlAccountIDs []string, prPrefixes []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation) (*int, error) {
	args := m.Called(ctx, organizationID, sourceControlAccountIDs, prPrefixes, startDate, endDate, metricOperation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int), args.Error(1)
}

func (m *MockDB) CalculatePRsMergedGraph(ctx context.Context, organizationID string, sourceControlAccountIDs []string, prPrefixes []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation, metricLabel string, interval string) ([]types.TimeSeriesEntry, error) {
	args := m.Called(ctx, organizationID, sourceControlAccountIDs, prPrefixes, startDate, endDate, metricOperation, metricLabel, interval)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.TimeSeriesEntry), args.Error(1)
}

func (m *MockDB) CalculatePRsReviewed(ctx context.Context, organizationID string, sourceControlAccountIDs []string, prPrefixes []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation) (*int, error) {
	args := m.Called(ctx, organizationID, sourceControlAccountIDs, prPrefixes, startDate, endDate, metricOperation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int), args.Error(1)
}

func (m *MockDB) CalculatePRsReviewedGraph(ctx context.Context, organizationID string, sourceControlAccountIDs []string, prPrefixes []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation, metricLabel string, interval string) ([]types.TimeSeriesEntry, error) {
	args := m.Called(ctx, organizationID, sourceControlAccountIDs, prPrefixes, startDate, endDate, metricOperation, metricLabel, interval)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.TimeSeriesEntry), args.Error(1)
}

func (m *MockDB) CalculateLOCAdded(ctx context.Context, organizationID string, sourceControlAccountIDs []string, prPrefixes []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation) (*int, error) {
	args := m.Called(ctx, organizationID, sourceControlAccountIDs, prPrefixes, startDate, endDate, metricOperation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int), args.Error(1)
}

func (m *MockDB) CalculateLOCAddedGraph(ctx context.Context, organizationID string, sourceControlAccountIDs []string, prPrefixes []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation, metricLabel string, interval string) ([]types.TimeSeriesEntry, error) {
	args := m.Called(ctx, organizationID, sourceControlAccountIDs, prPrefixes, startDate, endDate, metricOperation, metricLabel, interval)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.TimeSeriesEntry), args.Error(1)
}

func (m *MockDB) CalculateLOCRemoved(ctx context.Context, organizationID string, sourceControlAccountIDs []string, prPrefixes []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation) (*int, error) {
	args := m.Called(ctx, organizationID, sourceControlAccountIDs, prPrefixes, startDate, endDate, metricOperation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int), args.Error(1)
}

func (m *MockDB) CalculateLOCRemovedGraph(ctx context.Context, organizationID string, sourceControlAccountIDs []string, prPrefixes []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation, metricLabel string, interval string) ([]types.TimeSeriesEntry, error) {
	args := m.Called(ctx, organizationID, sourceControlAccountIDs, prPrefixes, startDate, endDate, metricOperation, metricLabel, interval)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.TimeSeriesEntry), args.Error(1)
}

func (m *MockDB) CalculatePRReviewComplexity(ctx context.Context, organizationID string, sourceControlAccountIDs []string, prPrefixes []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation) (*float64, error) {
	args := m.Called(ctx, organizationID, sourceControlAccountIDs, prPrefixes, startDate, endDate, metricOperation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*float64), args.Error(1)
}

func (m *MockDB) CalculatePRReviewComplexityGraph(ctx context.Context, organizationID string, sourceControlAccountIDs []string, prPrefixes []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation, metricLabel string, interval string) ([]types.TimeSeriesEntry, error) {
	args := m.Called(ctx, organizationID, sourceControlAccountIDs, prPrefixes, startDate, endDate, metricOperation, metricLabel, interval)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.TimeSeriesEntry), args.Error(1)
}

func (m *MockDB) CalculateLOCAddedForAccounts(ctx context.Context, organizationID string, sourceControlAccountIDs []string, prPrefixes []string, startDate, endDate time.Time) (*float64, error) {
	args := m.Called(ctx, organizationID, sourceControlAccountIDs, prPrefixes, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*float64), args.Error(1)
}

func (m *MockDB) CalculateLOCRemovedForAccounts(ctx context.Context, organizationID string, sourceControlAccountIDs []string, prPrefixes []string, startDate, endDate time.Time) (*float64, error) {
	args := m.Called(ctx, organizationID, sourceControlAccountIDs, prPrefixes, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*float64), args.Error(1)
}

func (m *MockDB) CalculatePRsMergedForAccounts(ctx context.Context, organizationID string, sourceControlAccountIDs []string, prPrefixes []string, startDate, endDate time.Time) (*float64, error) {
	args := m.Called(ctx, organizationID, sourceControlAccountIDs, prPrefixes, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*float64), args.Error(1)
}

func (m *MockDB) CalculatePRsReviewedForAccounts(ctx context.Context, organizationID string, sourceControlAccountIDs []string, prPrefixes []string, startDate, endDate time.Time) (*float64, error) {
	args := m.Called(ctx, organizationID, sourceControlAccountIDs, prPrefixes, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*float64), args.Error(1)
}

func (m *MockDB) CalculateTimeToMergeForAccounts(ctx context.Context, organizationID string, sourceControlAccountIDs []string, prPrefixes []string, startDate, endDate time.Time) (*float64, error) {
	args := m.Called(ctx, organizationID, sourceControlAccountIDs, prPrefixes, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*float64), args.Error(1)
}

func (m *MockDB) CalculatePRReviewComplexityForAccounts(ctx context.Context, organizationID string, sourceControlAccountIDs []string, prPrefixes []string, startDate, endDate time.Time) (*float64, error) {
	args := m.Called(ctx, organizationID, sourceControlAccountIDs, prPrefixes, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*float64), args.Error(1)
}

func (m *MockDB) CalculateLOCAddedGraphForAccounts(ctx context.Context, organizationID string, sourceControlAccountIDs []string, prPrefixes []string, startDate, endDate time.Time, interval string) ([]types.TimeSeriesEntry, error) {
	args := m.Called(ctx, organizationID, sourceControlAccountIDs, prPrefixes, startDate, endDate, interval)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.TimeSeriesEntry), args.Error(1)
}

func (m *MockDB) CalculateLOCRemovedGraphForAccounts(ctx context.Context, organizationID string, sourceControlAccountIDs []string, prPrefixes []string, startDate, endDate time.Time, interval string) ([]types.TimeSeriesEntry, error) {
	args := m.Called(ctx, organizationID, sourceControlAccountIDs, prPrefixes, startDate, endDate, interval)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.TimeSeriesEntry), args.Error(1)
}

func (m *MockDB) CalculatePRsMergedGraphForAccounts(ctx context.Context, organizationID string, sourceControlAccountIDs []string, prPrefixes []string, startDate, endDate time.Time, interval string) ([]types.TimeSeriesEntry, error) {
	args := m.Called(ctx, organizationID, sourceControlAccountIDs, prPrefixes, startDate, endDate, interval)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.TimeSeriesEntry), args.Error(1)
}

func (m *MockDB) CalculatePRsReviewedGraphForAccounts(ctx context.Context, organizationID string, sourceControlAccountIDs []string, prPrefixes []string, startDate, endDate time.Time, interval string) ([]types.TimeSeriesEntry, error) {
	args := m.Called(ctx, organizationID, sourceControlAccountIDs, prPrefixes, startDate, endDate, interval)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.TimeSeriesEntry), args.Error(1)
}

func (m *MockDB) CalculateTimeToMergeGraphForAccounts(ctx context.Context, organizationID string, sourceControlAccountIDs []string, prPrefixes []string, startDate, endDate time.Time, interval string) ([]types.TimeSeriesEntry, error) {
	args := m.Called(ctx, organizationID, sourceControlAccountIDs, prPrefixes, startDate, endDate, interval)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.TimeSeriesEntry), args.Error(1)
}

func (m *MockDB) CalculatePRReviewComplexityGraphForAccounts(ctx context.Context, organizationID string, sourceControlAccountIDs []string, prPrefixes []string, startDate, endDate time.Time, interval string) ([]types.TimeSeriesEntry, error) {
	args := m.Called(ctx, organizationID, sourceControlAccountIDs, prPrefixes, startDate, endDate, interval)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.TimeSeriesEntry), args.Error(1)
}

// MockMetricsEngine is a mock implementation of the metrics.MetricsEngine interface
type MockMetricsEngine struct {
	mock.Mock
}

func (m *MockMetricsEngine) CalculateMetrics(ctx context.Context, params types.MetricRuleParams) (*types.MetricsResponse, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.MetricsResponse), args.Error(1)
}
