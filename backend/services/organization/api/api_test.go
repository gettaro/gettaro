package api

import (
	"context"
	"errors"
	"testing"

	"ems.dev/backend/services/organization/types"
	usertypes "ems.dev/backend/services/user/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOrganizationDB is a mock implementation of the OrganizationDB type
type MockOrganizationDB struct {
	mock.Mock
}

// MockUserAPI is a mock implementation of the UserAPI type
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

func (m *MockUserAPI) GetOrCreateUserFromAuthProvider(provider string, providerID string, email string, name string) (*usertypes.User, error) {
	args := m.Called(provider, providerID, email, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usertypes.User), args.Error(1)
}

func (m *MockOrganizationDB) CreateOrganization(org *types.Organization, ownerID string) error {
	args := m.Called(org, ownerID)
	return args.Error(0)
}

func (m *MockOrganizationDB) GetUserOrganizations(userID string) ([]types.Organization, error) {
	args := m.Called(userID)
	return args.Get(0).([]types.Organization), args.Error(1)
}

func (m *MockOrganizationDB) GetOrganizationByID(id string) (*types.Organization, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Organization), args.Error(1)
}

func (m *MockOrganizationDB) UpdateOrganization(org *types.Organization) error {
	args := m.Called(org)
	return args.Error(0)
}

func (m *MockOrganizationDB) DeleteOrganization(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockOrganizationDB) AddOrganizationMember(orgID string, userID string) error {
	args := m.Called(orgID, userID)
	return args.Error(0)
}

func (m *MockOrganizationDB) RemoveOrganizationMember(orgID string, userID string) error {
	args := m.Called(orgID, userID)
	return args.Error(0)
}

func (m *MockOrganizationDB) GetOrganizationMembers(orgID string) ([]types.OrganizationMember, error) {
	args := m.Called(orgID)
	return args.Get(0).([]types.OrganizationMember), args.Error(1)
}

func (m *MockOrganizationDB) IsOrganizationOwner(orgID string, userID string) (bool, error) {
	args := m.Called(orgID, userID)
	return args.Bool(0), args.Error(1)
}

func TestCreateOrganization(t *testing.T) {
	tests := []struct {
		name          string
		org           *types.Organization
		ownerID       string
		mockError     error
		expectedError error
	}{
		{
			name: "successful creation",
			org: &types.Organization{
				Name: "Test Org",
				Slug: "test-org",
			},
			ownerID: "user-1",
		},
		{
			name: "database error",
			org: &types.Organization{
				Name: "Test Org",
				Slug: "test-org",
			},
			ownerID:       "user-1",
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockOrganizationDB)
			mockUserAPI := new(MockUserAPI)
			api := NewApi(mockDB, mockUserAPI)

			mockDB.On("CreateOrganization", tt.org, tt.ownerID).Return(tt.mockError)

			err := api.CreateOrganization(context.Background(), tt.org, tt.ownerID)

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

func TestGetUserOrganizations(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		mockOrgs      []types.Organization
		mockError     error
		expectedOrgs  []types.Organization
		expectedError error
	}{
		{
			name:   "successful retrieval",
			userID: "user-1",
			mockOrgs: []types.Organization{
				{
					ID:   "org-1",
					Name: "Test Org 1",
				},
				{
					ID:   "org-2",
					Name: "Test Org 2",
				},
			},
			expectedOrgs: []types.Organization{
				{
					ID:   "org-1",
					Name: "Test Org 1",
				},
				{
					ID:   "org-2",
					Name: "Test Org 2",
				},
			},
		},
		{
			name:         "no organizations",
			userID:       "user-1",
			mockOrgs:     []types.Organization{},
			expectedOrgs: []types.Organization{},
		},
		{
			name:          "database error",
			userID:        "user-1",
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockOrganizationDB)
			mockUserAPI := new(MockUserAPI)
			api := NewApi(mockDB, mockUserAPI)

			mockDB.On("GetUserOrganizations", tt.userID).Return(tt.mockOrgs, tt.mockError)

			orgs, err := api.GetUserOrganizations(context.Background(), tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, orgs)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOrgs, orgs)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func TestGetOrganizationByID(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockOrg       *types.Organization
		mockError     error
		expectedOrg   *types.Organization
		expectedError error
	}{
		{
			name: "successful retrieval",
			id:   "org-1",
			mockOrg: &types.Organization{
				ID:   "org-1",
				Name: "Test Org",
			},
			expectedOrg: &types.Organization{
				ID:   "org-1",
				Name: "Test Org",
			},
		},
		{
			name:          "organization not found",
			id:            "org-1",
			mockOrg:       nil,
			expectedOrg:   nil,
			expectedError: nil,
		},
		{
			name:          "database error",
			id:            "org-1",
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockOrganizationDB)
			mockUserAPI := new(MockUserAPI)
			api := NewApi(mockDB, mockUserAPI)

			mockDB.On("GetOrganizationByID", tt.id).Return(tt.mockOrg, tt.mockError)

			org, err := api.GetOrganizationByID(context.Background(), tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, org)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOrg, org)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func TestUpdateOrganization(t *testing.T) {
	tests := []struct {
		name          string
		org           *types.Organization
		mockError     error
		expectedError error
	}{
		{
			name: "successful update",
			org: &types.Organization{
				ID:   "org-1",
				Name: "Updated Org",
			},
		},
		{
			name: "database error",
			org: &types.Organization{
				ID:   "org-1",
				Name: "Updated Org",
			},
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockOrganizationDB)
			mockUserAPI := new(MockUserAPI)
			api := NewApi(mockDB, mockUserAPI)

			mockDB.On("UpdateOrganization", tt.org).Return(tt.mockError)

			err := api.UpdateOrganization(context.Background(), tt.org)

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

func TestDeleteOrganization(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockError     error
		expectedError error
	}{
		{
			name: "successful deletion",
			id:   "org-1",
		},
		{
			name:          "database error",
			id:            "org-1",
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockOrganizationDB)
			mockUserAPI := new(MockUserAPI)
			api := NewApi(mockDB, mockUserAPI)

			mockDB.On("DeleteOrganization", tt.id).Return(tt.mockError)

			err := api.DeleteOrganization(context.Background(), tt.id)

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

func TestAddOrganizationMemberByEmail(t *testing.T) {
	tests := []struct {
		name          string
		orgID         string
		email         string
		user          *usertypes.User
		userError     error
		mockError     error
		expectedError error
	}{
		{
			name:  "successful addition",
			orgID: "org-1",
			email: "user@example.com",
			user: &usertypes.User{
				ID: "user-1",
			},
		},
		{
			name:          "user not found",
			orgID:         "org-1",
			email:         "nonexistent@example.com",
			user:          nil,
			expectedError: errors.New("user not found"),
		},
		{
			name:          "database error",
			orgID:         "org-1",
			email:         "user@example.com",
			user:          &usertypes.User{ID: "user-1"},
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
		{
			name:          "user lookup error",
			orgID:         "org-1",
			email:         "user@example.com",
			userError:     errors.New("user lookup error"),
			expectedError: errors.New("user lookup error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockOrganizationDB)
			mockUserAPI := new(MockUserAPI)
			api := NewApi(mockDB, mockUserAPI)

			mockUserAPI.On("FindUser", usertypes.UserSearchParams{Email: &tt.email}).Return(tt.user, tt.userError)
			if tt.user != nil {
				mockDB.On("AddOrganizationMember", tt.orgID, tt.user.ID).Return(tt.mockError)
			}

			err := api.AddOrganizationMemberByEmail(context.Background(), tt.orgID, tt.email)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
			mockUserAPI.AssertExpectations(t)
		})
	}
}

func TestRemoveOrganizationMember(t *testing.T) {
	tests := []struct {
		name          string
		orgID         string
		userID        string
		mockError     error
		expectedError error
	}{
		{
			name:   "successful removal",
			orgID:  "org-1",
			userID: "user-1",
		},
		{
			name:          "database error",
			orgID:         "org-1",
			userID:        "user-1",
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockOrganizationDB)
			mockUserAPI := new(MockUserAPI)
			api := NewApi(mockDB, mockUserAPI)

			mockDB.On("RemoveOrganizationMember", tt.orgID, tt.userID).Return(tt.mockError)

			err := api.RemoveOrganizationMember(context.Background(), tt.orgID, tt.userID)

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

func TestGetOrganizationMembers(t *testing.T) {
	tests := []struct {
		name            string
		orgID           string
		mockMembers     []types.OrganizationMember
		mockError       error
		expectedMembers []types.OrganizationMember
		expectedError   error
	}{
		{
			name:  "successful retrieval",
			orgID: "org-1",
			mockMembers: []types.OrganizationMember{
				{
					ID:             "member-1",
					UserID:         "user-1",
					OrganizationID: "org-1",
				},
			},
			expectedMembers: []types.OrganizationMember{
				{
					ID:             "member-1",
					UserID:         "user-1",
					OrganizationID: "org-1",
				},
			},
		},
		{
			name:            "no members",
			orgID:           "org-1",
			mockMembers:     []types.OrganizationMember{},
			expectedMembers: []types.OrganizationMember{},
		},
		{
			name:          "database error",
			orgID:         "org-1",
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockOrganizationDB)
			mockUserAPI := new(MockUserAPI)
			api := NewApi(mockDB, mockUserAPI)

			mockDB.On("GetOrganizationMembers", tt.orgID).Return(tt.mockMembers, tt.mockError)

			members, err := api.GetOrganizationMembers(context.Background(), tt.orgID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, members)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMembers, members)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func TestIsOrganizationOwner(t *testing.T) {
	tests := []struct {
		name          string
		orgID         string
		userID        string
		mockIsOwner   bool
		mockError     error
		expectedOwner bool
		expectedError error
	}{
		{
			name:          "user is owner",
			orgID:         "org-1",
			userID:        "user-1",
			mockIsOwner:   true,
			expectedOwner: true,
		},
		{
			name:          "user is not owner",
			orgID:         "org-1",
			userID:        "user-1",
			mockIsOwner:   false,
			expectedOwner: false,
		},
		{
			name:          "database error",
			orgID:         "org-1",
			userID:        "user-1",
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockOrganizationDB)
			mockUserAPI := new(MockUserAPI)
			api := NewApi(mockDB, mockUserAPI)

			mockDB.On("IsOrganizationOwner", tt.orgID, tt.userID).Return(tt.mockIsOwner, tt.mockError)

			isOwner, err := api.IsOrganizationOwner(context.Background(), tt.orgID, tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOwner, isOwner)
			}

			mockDB.AssertExpectations(t)
		})
	}
}
