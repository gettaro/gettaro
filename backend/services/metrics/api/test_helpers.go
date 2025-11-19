package api

import (
	"context"

	membertypes "ems.dev/backend/services/member/types"
	sourcecontroltypes "ems.dev/backend/services/sourcecontrol/types"
	teamtypes "ems.dev/backend/services/team/types"
	aicodeassistanttypes "ems.dev/backend/services/aicodeassistant/types"
	"github.com/stretchr/testify/mock"
)

// MockMemberAPI is a mock implementation of MemberAPI
type MockMemberAPI struct {
	mock.Mock
}

func (m *MockMemberAPI) GetExternalAccounts(ctx context.Context, params *membertypes.ExternalAccountParams) ([]membertypes.ExternalAccount, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]membertypes.ExternalAccount), args.Error(1)
}

// Stub implementations for other MemberAPI methods
func (m *MockMemberAPI) AddOrganizationMember(ctx context.Context, req membertypes.AddMemberRequest, member *membertypes.OrganizationMember) (*membertypes.OrganizationMember, error) {
	args := m.Called(ctx, req, member)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*membertypes.OrganizationMember), args.Error(1)
}

func (m *MockMemberAPI) RemoveOrganizationMember(ctx context.Context, orgID string, userID string) error {
	args := m.Called(ctx, orgID, userID)
	return args.Error(0)
}

func (m *MockMemberAPI) GetOrganizationMembers(ctx context.Context, orgID string, params *membertypes.OrganizationMemberParams) ([]membertypes.OrganizationMember, error) {
	args := m.Called(ctx, orgID, params)
	return args.Get(0).([]membertypes.OrganizationMember), args.Error(1)
}

func (m *MockMemberAPI) GetOrganizationMemberByID(ctx context.Context, memberID string) (*membertypes.OrganizationMember, error) {
	args := m.Called(ctx, memberID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*membertypes.OrganizationMember), args.Error(1)
}

func (m *MockMemberAPI) IsOrganizationOwner(ctx context.Context, orgID string, userID string) (bool, error) {
	args := m.Called(ctx, orgID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockMemberAPI) UpdateOrganizationMember(ctx context.Context, orgID string, memberID string, req membertypes.UpdateMemberRequest) error {
	args := m.Called(ctx, orgID, memberID, req)
	return args.Error(0)
}

func (m *MockMemberAPI) CalculateSourceControlMemberMetrics(ctx context.Context, organizationID string, memberID string, params sourcecontroltypes.MemberMetricsParams) (*sourcecontroltypes.MetricsResponse, error) {
	args := m.Called(ctx, organizationID, memberID, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sourcecontroltypes.MetricsResponse), args.Error(1)
}

func (m *MockMemberAPI) CreateExternalAccounts(ctx context.Context, accounts []*membertypes.ExternalAccount) error {
	args := m.Called(ctx, accounts)
	return args.Error(0)
}

func (m *MockMemberAPI) GetExternalAccount(ctx context.Context, id string) (*membertypes.ExternalAccount, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*membertypes.ExternalAccount), args.Error(1)
}

func (m *MockMemberAPI) UpdateExternalAccount(ctx context.Context, account *membertypes.ExternalAccount) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockMemberAPI) UpdateExternalAccountMemberID(ctx context.Context, organizationID string, accountID string, memberID *string) (*membertypes.ExternalAccount, error) {
	args := m.Called(ctx, organizationID, accountID, memberID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*membertypes.ExternalAccount), args.Error(1)
}

// MockTeamAPI is a mock implementation of TeamAPI
type MockTeamAPI struct {
	mock.Mock
}

func (m *MockTeamAPI) ListTeams(ctx context.Context, params teamtypes.TeamSearchParams) ([]teamtypes.Team, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]teamtypes.Team), args.Error(1)
}

func (m *MockTeamAPI) GetTeamByOrganization(ctx context.Context, teamID, organizationID string) (*teamtypes.Team, error) {
	args := m.Called(ctx, teamID, organizationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*teamtypes.Team), args.Error(1)
}

// Stub implementations for other TeamAPI methods
func (m *MockTeamAPI) CreateTeam(ctx context.Context, team *teamtypes.Team) error {
	args := m.Called(ctx, team)
	return args.Error(0)
}

func (m *MockTeamAPI) GetTeam(ctx context.Context, id string) (*teamtypes.Team, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*teamtypes.Team), args.Error(1)
}

func (m *MockTeamAPI) UpdateTeam(ctx context.Context, id string, team *teamtypes.Team) error {
	args := m.Called(ctx, id, team)
	return args.Error(0)
}

func (m *MockTeamAPI) DeleteTeam(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTeamAPI) AddTeamMember(ctx context.Context, teamID string, member *teamtypes.TeamMember) error {
	args := m.Called(ctx, teamID, member)
	return args.Error(0)
}

func (m *MockTeamAPI) RemoveTeamMember(ctx context.Context, teamID string, memberID string) error {
	args := m.Called(ctx, teamID, memberID)
	return args.Error(0)
}

// MockSourceControlAPI is a mock implementation of SourceControlAPI
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

// Stub implementations for other SourceControlAPI methods
func (m *MockSourceControlAPI) GetPullRequests(ctx context.Context, params *sourcecontroltypes.PullRequestParams) ([]*sourcecontroltypes.PullRequest, error) {
	args := m.Called(ctx, params)
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
	return args.Get(0).([]*sourcecontroltypes.PRComment), args.Error(1)
}

func (m *MockSourceControlAPI) GetMemberPullRequests(ctx context.Context, params *sourcecontroltypes.MemberPullRequestParams) ([]*sourcecontroltypes.PullRequestWithComments, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]*sourcecontroltypes.PullRequestWithComments), args.Error(1)
}

func (m *MockSourceControlAPI) GetMemberPullRequestReviews(ctx context.Context, params *sourcecontroltypes.MemberPullRequestReviewsParams) ([]*sourcecontroltypes.MemberActivity, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]*sourcecontroltypes.MemberActivity), args.Error(1)
}

// MockAICodeAssistantAPI is a mock implementation of AICodeAssistantAPI
type MockAICodeAssistantAPI struct {
	mock.Mock
}

func (m *MockAICodeAssistantAPI) CalculateMetrics(ctx context.Context, params aicodeassistanttypes.MetricRuleParams) (*aicodeassistanttypes.MetricsResponse, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*aicodeassistanttypes.MetricsResponse), args.Error(1)
}

// Stub implementations for other AICodeAssistantAPI methods
func (m *MockAICodeAssistantAPI) CreateOrUpdateDailyMetric(ctx context.Context, metric *aicodeassistanttypes.AICodeAssistantDailyMetric) error {
	args := m.Called(ctx, metric)
	return args.Error(0)
}

func (m *MockAICodeAssistantAPI) GetDailyMetrics(ctx context.Context, params *aicodeassistanttypes.AICodeAssistantDailyMetricParams) ([]*aicodeassistanttypes.AICodeAssistantDailyMetric, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]*aicodeassistanttypes.AICodeAssistantDailyMetric), args.Error(1)
}

func (m *MockAICodeAssistantAPI) GetMemberDailyMetrics(ctx context.Context, organizationID, memberID string, params *aicodeassistanttypes.AICodeAssistantMemberMetricsParams) ([]*aicodeassistanttypes.AICodeAssistantDailyMetric, error) {
	args := m.Called(ctx, organizationID, memberID, params)
	return args.Get(0).([]*aicodeassistanttypes.AICodeAssistantDailyMetric), args.Error(1)
}

func (m *MockAICodeAssistantAPI) CalculateMemberMetrics(ctx context.Context, organizationID string, memberID string, params aicodeassistanttypes.MemberMetricsParams) (*aicodeassistanttypes.MetricsResponse, error) {
	args := m.Called(ctx, organizationID, memberID, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*aicodeassistanttypes.MetricsResponse), args.Error(1)
}
