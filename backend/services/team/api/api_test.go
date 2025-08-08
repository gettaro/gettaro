package api

import (
	"context"
	"errors"
	"testing"

	orgtypes "ems.dev/backend/services/organization/types"
	"ems.dev/backend/services/team/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTeamDB is a mock implementation of the TeamDB type
type MockTeamDB struct {
	mock.Mock
}

// MockOrganizationAPI is a mock implementation of the OrganizationAPI type
type MockOrganizationAPI struct {
	mock.Mock
}

func (m *MockTeamDB) CreateTeam(ctx context.Context, team *types.Team) error {
	args := m.Called(ctx, team)
	return args.Error(0)
}

func (m *MockTeamDB) ListTeams(ctx context.Context, params types.TeamSearchParams) ([]types.Team, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]types.Team), args.Error(1)
}

func (m *MockTeamDB) GetTeam(ctx context.Context, id string) (*types.Team, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Team), args.Error(1)
}

func (m *MockTeamDB) UpdateTeam(ctx context.Context, id string, team *types.Team) error {
	args := m.Called(ctx, id, team)
	return args.Error(0)
}

func (m *MockTeamDB) DeleteTeam(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTeamDB) AddTeamMember(ctx context.Context, teamID string, member *types.TeamMember) error {
	args := m.Called(ctx, teamID, member)
	return args.Error(0)
}

func (m *MockTeamDB) RemoveTeamMember(ctx context.Context, teamID string, userID string) error {
	args := m.Called(ctx, teamID, userID)
	return args.Error(0)
}

func (m *MockOrganizationAPI) CreateOrganization(ctx context.Context, org *orgtypes.Organization, ownerID string) error {
	args := m.Called(ctx, org, ownerID)
	return args.Error(0)
}

func (m *MockOrganizationAPI) GetOrganizations(ctx context.Context) ([]orgtypes.Organization, error) {
	args := m.Called(ctx)
	return args.Get(0).([]orgtypes.Organization), args.Error(1)
}

func (m *MockOrganizationAPI) GetUserOrganizations(ctx context.Context, userID string) ([]orgtypes.Organization, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]orgtypes.Organization), args.Error(1)
}

func (m *MockOrganizationAPI) GetOrganizationByID(ctx context.Context, id string) (*orgtypes.Organization, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*orgtypes.Organization), args.Error(1)
}

func (m *MockOrganizationAPI) UpdateOrganization(ctx context.Context, org *orgtypes.Organization) error {
	args := m.Called(ctx, org)
	return args.Error(0)
}

func (m *MockOrganizationAPI) DeleteOrganization(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrganizationAPI) AddOrganizationMemberByEmail(ctx context.Context, orgID string, email string) error {
	args := m.Called(ctx, orgID, email)
	return args.Error(0)
}

func (m *MockOrganizationAPI) RemoveOrganizationMember(ctx context.Context, orgID string, userID string) error {
	args := m.Called(ctx, orgID, userID)
	return args.Error(0)
}

func (m *MockOrganizationAPI) GetOrganizationMembers(ctx context.Context, orgID string) ([]orgtypes.OrganizationMember, error) {
	args := m.Called(ctx, orgID)
	return args.Get(0).([]orgtypes.OrganizationMember), args.Error(1)
}

func (m *MockOrganizationAPI) IsOrganizationOwner(ctx context.Context, orgID string, userID string) (bool, error) {
	args := m.Called(ctx, orgID, userID)
	return args.Bool(0), args.Error(1)
}

func TestCreateTeam(t *testing.T) {
	tests := []struct {
		name          string
		team          *types.Team
		mockError     error
		expectedError error
	}{
		{
			name: "successful creation",
			team: &types.Team{
				ID:             "team-1",
				Name:           "Test Team",
				OrganizationID: "org-1",
			},
		},
		{
			name: "database error",
			team: &types.Team{
				ID:             "team-1",
				Name:           "Test Team",
				OrganizationID: "org-1",
			},
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockTeamDB)
			mockOrgAPI := new(MockOrganizationAPI)
			api := NewApi(mockDB, mockOrgAPI)

			ctx := context.Background()
			mockOrgAPI.On("GetOrganizationByID", ctx, tt.team.OrganizationID).Return(&orgtypes.Organization{ID: tt.team.OrganizationID}, nil)
			mockDB.On("CreateTeam", ctx, tt.team).Return(tt.mockError)

			err := api.CreateTeam(ctx, tt.team)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
			mockOrgAPI.AssertExpectations(t)
		})
	}
}

func TestListTeams(t *testing.T) {
	tests := []struct {
		name          string
		params        types.TeamSearchParams
		mockTeams     []types.Team
		mockError     error
		expectedTeams []types.Team
		expectedError error
	}{
		{
			name:   "successful retrieval",
			params: types.TeamSearchParams{},
			mockTeams: []types.Team{
				{
					ID:   "team-1",
					Name: "Test Team 1",
				},
				{
					ID:   "team-2",
					Name: "Test Team 2",
				},
			},
			expectedTeams: []types.Team{
				{
					ID:   "team-1",
					Name: "Test Team 1",
				},
				{
					ID:   "team-2",
					Name: "Test Team 2",
				},
			},
		},
		{
			name:          "no teams",
			params:        types.TeamSearchParams{},
			mockTeams:     []types.Team{},
			expectedTeams: []types.Team{},
		},
		{
			name:          "database error",
			params:        types.TeamSearchParams{},
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockTeamDB)
			mockOrgAPI := new(MockOrganizationAPI)
			api := NewApi(mockDB, mockOrgAPI)

			ctx := context.Background()
			mockDB.On("ListTeams", ctx, tt.params).Return(tt.mockTeams, tt.mockError)

			teams, err := api.ListTeams(ctx, tt.params)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, teams)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTeams, teams)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func TestGetTeam(t *testing.T) {
	tests := []struct {
		name          string
		teamID        string
		mockTeam      *types.Team
		mockError     error
		expectedTeam  *types.Team
		expectedError error
	}{
		{
			name:   "successful retrieval",
			teamID: "team-1",
			mockTeam: &types.Team{
				ID:   "team-1",
				Name: "Test Team",
			},
			expectedTeam: &types.Team{
				ID:   "team-1",
				Name: "Test Team",
			},
		},
		{
			name:          "team not found",
			teamID:        "team-1",
			mockError:     errors.New("team not found"),
			expectedError: errors.New("team not found"),
		},
		{
			name:          "database error",
			teamID:        "team-1",
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockTeamDB)
			mockOrgAPI := new(MockOrganizationAPI)
			api := NewApi(mockDB, mockOrgAPI)

			ctx := context.Background()
			mockDB.On("GetTeam", ctx, tt.teamID).Return(tt.mockTeam, tt.mockError)

			team, err := api.GetTeam(ctx, tt.teamID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, team)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTeam, team)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func TestUpdateTeam(t *testing.T) {
	tests := []struct {
		name          string
		teamID        string
		team          *types.Team
		mockError     error
		expectedError error
	}{
		{
			name:   "successful update",
			teamID: "team-1",
			team: &types.Team{
				ID:   "team-1",
				Name: "Updated Team",
			},
		},
		{
			name:   "database error",
			teamID: "team-1",
			team: &types.Team{
				ID:   "team-1",
				Name: "Updated Team",
			},
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockTeamDB)
			mockOrgAPI := new(MockOrganizationAPI)
			api := NewApi(mockDB, mockOrgAPI)

			ctx := context.Background()
			mockDB.On("UpdateTeam", ctx, tt.teamID, tt.team).Return(tt.mockError)

			err := api.UpdateTeam(ctx, tt.teamID, tt.team)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func TestDeleteTeam(t *testing.T) {
	tests := []struct {
		name          string
		teamID        string
		mockError     error
		expectedError error
	}{
		{
			name:   "successful deletion",
			teamID: "team-1",
		},
		{
			name:          "database error",
			teamID:        "team-1",
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockTeamDB)
			mockOrgAPI := new(MockOrganizationAPI)
			api := NewApi(mockDB, mockOrgAPI)

			ctx := context.Background()
			mockDB.On("DeleteTeam", ctx, tt.teamID).Return(tt.mockError)

			err := api.DeleteTeam(ctx, tt.teamID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func TestAddTeamMember(t *testing.T) {
	tests := []struct {
		name          string
		teamID        string
		member        *types.TeamMember
		mockError     error
		expectedError error
	}{
		{
			name:   "successful addition",
			teamID: "team-1",
			member: &types.TeamMember{
				ID:     "member-1",
				UserID: "user-1",
				TeamID: "team-1",
			},
		},
		{
			name:   "database error",
			teamID: "team-1",
			member: &types.TeamMember{
				ID:     "member-1",
				UserID: "user-1",
				TeamID: "team-1",
			},
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockTeamDB)
			mockOrgAPI := new(MockOrganizationAPI)
			api := NewApi(mockDB, mockOrgAPI)

			ctx := context.Background()
			mockDB.On("AddTeamMember", ctx, tt.teamID, tt.member).Return(tt.mockError)

			err := api.AddTeamMember(ctx, tt.teamID, tt.member)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func TestRemoveTeamMember(t *testing.T) {
	tests := []struct {
		name          string
		teamID        string
		userID        string
		mockError     error
		expectedError error
	}{
		{
			name:   "successful removal",
			teamID: "team-1",
			userID: "user-1",
		},
		{
			name:          "database error",
			teamID:        "team-1",
			userID:        "user-1",
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockTeamDB)
			mockOrgAPI := new(MockOrganizationAPI)
			api := NewApi(mockDB, mockOrgAPI)

			ctx := context.Background()
			mockDB.On("RemoveTeamMember", ctx, tt.teamID, tt.userID).Return(tt.mockError)

			err := api.RemoveTeamMember(ctx, tt.teamID, tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
		})
	}
}
