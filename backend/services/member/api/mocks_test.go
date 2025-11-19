package api

import (
	"context"

	directstypes "ems.dev/backend/services/directs/types"
	"ems.dev/backend/services/member/types"
	sourcecontroltypes "ems.dev/backend/services/sourcecontrol/types"
	titletypes "ems.dev/backend/services/title/types"
	usertypes "ems.dev/backend/services/user/types"
	"github.com/stretchr/testify/mock"
)

// MockMemberDB is a mock implementation of the member database interface
type MockMemberDB struct {
	mock.Mock
}

func (m *MockMemberDB) AddOrganizationMember(member *types.OrganizationMember) error {
	args := m.Called(member)
	return args.Error(0)
}

func (m *MockMemberDB) RemoveOrganizationMember(orgID string, userID string) error {
	args := m.Called(orgID, userID)
	return args.Error(0)
}

func (m *MockMemberDB) GetOrganizationMembers(orgID string, params *types.OrganizationMemberParams) ([]types.OrganizationMember, error) {
	args := m.Called(orgID, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.OrganizationMember), args.Error(1)
}

func (m *MockMemberDB) GetOrganizationMember(orgID string, userID string) (*types.OrganizationMember, error) {
	args := m.Called(orgID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.OrganizationMember), args.Error(1)
}

func (m *MockMemberDB) GetOrganizationMemberByID(ctx context.Context, memberID string) (*types.OrganizationMember, error) {
	args := m.Called(ctx, memberID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.OrganizationMember), args.Error(1)
}

func (m *MockMemberDB) IsOrganizationOwner(orgID string, userID string) (bool, error) {
	args := m.Called(orgID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockMemberDB) UpdateOrganizationMember(orgID string, userID string, username string, titleID *string) error {
	args := m.Called(orgID, userID, username, titleID)
	return args.Error(0)
}

func (m *MockMemberDB) GetExternalAccounts(ctx context.Context, params *types.ExternalAccountParams) ([]types.ExternalAccount, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.ExternalAccount), args.Error(1)
}

func (m *MockMemberDB) CreateExternalAccounts(ctx context.Context, accounts []*types.ExternalAccount) error {
	args := m.Called(ctx, accounts)
	return args.Error(0)
}

func (m *MockMemberDB) GetExternalAccount(ctx context.Context, id string) (*types.ExternalAccount, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.ExternalAccount), args.Error(1)
}

func (m *MockMemberDB) UpdateExternalAccount(ctx context.Context, account *types.ExternalAccount) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

// MockUserAPI is a mock implementation of the UserAPI interface
type MockUserAPI struct {
	mock.Mock
}

func (m *MockUserAPI) FindUser(params usertypes.UserSearchParams) (*usertypes.User, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usertypes.User), args.Error(1)
}

func (m *MockUserAPI) CreateUser(user *usertypes.User) (*usertypes.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usertypes.User), args.Error(1)
}

// MockDirectReportsAPI is a mock implementation of the DirectReportsAPI interface
type MockDirectReportsAPI struct {
	mock.Mock
}

func (m *MockDirectReportsAPI) CreateDirectReport(ctx context.Context, params directstypes.CreateDirectReportParams) (*directstypes.DirectReport, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*directstypes.DirectReport), args.Error(1)
}

func (m *MockDirectReportsAPI) GetMemberManager(ctx context.Context, reportMemberID, orgID string) (*directstypes.DirectReport, error) {
	args := m.Called(ctx, reportMemberID, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*directstypes.DirectReport), args.Error(1)
}

func (m *MockDirectReportsAPI) DeleteDirectReport(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDirectReportsAPI) GetDirectReport(ctx context.Context, id string) (*directstypes.DirectReport, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*directstypes.DirectReport), args.Error(1)
}

func (m *MockDirectReportsAPI) ListDirectReports(ctx context.Context, params directstypes.DirectReportSearchParams) ([]directstypes.DirectReport, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]directstypes.DirectReport), args.Error(1)
}

func (m *MockDirectReportsAPI) UpdateDirectReport(ctx context.Context, id string, params directstypes.UpdateDirectReportParams) error {
	args := m.Called(ctx, id, params)
	return args.Error(0)
}

func (m *MockDirectReportsAPI) GetManagerDirectReports(ctx context.Context, managerMemberID, orgID string) ([]directstypes.DirectReport, error) {
	args := m.Called(ctx, managerMemberID, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]directstypes.DirectReport), args.Error(1)
}

func (m *MockDirectReportsAPI) GetManagerTree(ctx context.Context, managerMemberID, orgID string) ([]directstypes.OrgChartNode, error) {
	args := m.Called(ctx, managerMemberID, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]directstypes.OrgChartNode), args.Error(1)
}

func (m *MockDirectReportsAPI) GetMemberManagementChain(ctx context.Context, reportMemberID, orgID string) ([]directstypes.ManagementChain, error) {
	args := m.Called(ctx, reportMemberID, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]directstypes.ManagementChain), args.Error(1)
}

func (m *MockDirectReportsAPI) GetOrgChart(ctx context.Context, orgID string) ([]directstypes.OrgChartNode, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]directstypes.OrgChartNode), args.Error(1)
}

func (m *MockDirectReportsAPI) GetOrgChartFlat(ctx context.Context, orgID string) ([]directstypes.DirectReport, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]directstypes.DirectReport), args.Error(1)
}

// MockSourceControlAPI is a mock implementation of the SourceControlAPI interface
type MockSourceControlAPI struct {
	mock.Mock
}

func (m *MockSourceControlAPI) CalculateMetrics(ctx context.Context, params sourcecontroltypes.MetricRuleParams) (*sourcecontroltypes.MetricsResponse, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sourcecontroltypes.MetricsResponse), args.Error(1)
}

func (m *MockSourceControlAPI) GetPullRequests(ctx context.Context, params *sourcecontroltypes.PullRequestParams) ([]*sourcecontroltypes.PullRequest, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*sourcecontroltypes.PullRequest), args.Error(1)
}

func (m *MockSourceControlAPI) CreatePullRequest(ctx context.Context, pr *sourcecontroltypes.PullRequest) (*sourcecontroltypes.PullRequest, error) {
	args := m.Called(ctx, pr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sourcecontroltypes.PullRequest), args.Error(1)
}

func (m *MockSourceControlAPI) UpdatePullRequest(ctx context.Context, pr *sourcecontroltypes.PullRequest) error {
	args := m.Called(ctx, pr)
	return args.Error(0)
}

func (m *MockSourceControlAPI) CreatePRComments(ctx context.Context, comments []*sourcecontroltypes.PRComment) error {
	args := m.Called(ctx, comments)
	return args.Error(0)
}

func (m *MockSourceControlAPI) GetPullRequestComments(ctx context.Context, prID string) ([]*sourcecontroltypes.PRComment, error) {
	args := m.Called(ctx, prID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*sourcecontroltypes.PRComment), args.Error(1)
}

func (m *MockSourceControlAPI) GetMemberPullRequests(ctx context.Context, params *sourcecontroltypes.MemberPullRequestParams) ([]*sourcecontroltypes.PullRequestWithComments, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*sourcecontroltypes.PullRequestWithComments), args.Error(1)
}

func (m *MockSourceControlAPI) GetMemberPullRequestReviews(ctx context.Context, params *sourcecontroltypes.MemberPullRequestReviewsParams) ([]*sourcecontroltypes.MemberActivity, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*sourcecontroltypes.MemberActivity), args.Error(1)
}

// MockTitleAPI is a mock implementation of the TitleAPI interface
type MockTitleAPI struct {
	mock.Mock
}

func (m *MockTitleAPI) CreateTitle(ctx context.Context, title titletypes.Title) (*titletypes.Title, error) {
	args := m.Called(ctx, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*titletypes.Title), args.Error(1)
}

func (m *MockTitleAPI) GetTitle(ctx context.Context, id string) (*titletypes.Title, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*titletypes.Title), args.Error(1)
}

func (m *MockTitleAPI) ListTitles(ctx context.Context, orgID string) ([]titletypes.Title, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]titletypes.Title), args.Error(1)
}

func (m *MockTitleAPI) UpdateTitle(ctx context.Context, id string, title titletypes.Title) (*titletypes.Title, error) {
	args := m.Called(ctx, id, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*titletypes.Title), args.Error(1)
}

func (m *MockTitleAPI) DeleteTitle(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTitleAPI) AssignMemberTitle(ctx context.Context, memberTitle titletypes.MemberTitle) error {
	args := m.Called(ctx, memberTitle)
	return args.Error(0)
}

func (m *MockTitleAPI) GetMemberTitle(ctx context.Context, memberID string, orgID string) (*titletypes.MemberTitle, error) {
	args := m.Called(ctx, memberID, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*titletypes.MemberTitle), args.Error(1)
}

func (m *MockTitleAPI) RemoveMemberTitle(ctx context.Context, memberID string, orgID string) error {
	args := m.Called(ctx, memberID, orgID)
	return args.Error(0)
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
