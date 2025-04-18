package api

import (
	"errors"
	"testing"

	orgtypes "ems.dev/backend/services/organization/types"
	"ems.dev/backend/services/user/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserDB is a mock implementation of the UserDB type
type MockUserDB struct {
	mock.Mock
}

func (m *MockUserDB) FindUser(params types.UserSearchParams) (*types.User, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.User), args.Error(1)
}

func (m *MockUserDB) GetOrCreateUserFromAuthProvider(provider string, providerID string, email string, name string) (*types.User, error) {
	args := m.Called(provider, providerID, email, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.User), args.Error(1)
}

func (m *MockUserDB) CreateOrganizationWithOwner(org *orgtypes.Organization, userID string) error {
	args := m.Called(org, userID)
	return args.Error(0)
}

func (m *MockUserDB) GetUserOrganizations(userID string) ([]orgtypes.Organization, error) {
	args := m.Called(userID)
	return args.Get(0).([]orgtypes.Organization), args.Error(1)
}

func (m *MockUserDB) CreateUser(user *types.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserDB) UpdateUser(user *types.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserDB) DeleteUser(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserDB) GetUserByID(userID string) (*types.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.User), args.Error(1)
}

func (m *MockUserDB) GetUserByEmail(email string) (*types.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.User), args.Error(1)
}

func (m *MockUserDB) ListUsers() ([]types.User, error) {
	args := m.Called()
	return args.Get(0).([]types.User), args.Error(1)
}

func TestFindUser(t *testing.T) {
	tests := []struct {
		name          string
		params        types.UserSearchParams
		mockUser      *types.User
		mockError     error
		expectedUser  *types.User
		expectedError error
	}{
		{
			name: "find by email",
			params: types.UserSearchParams{
				Email: stringPtr("test@example.com"),
			},
			mockUser: &types.User{
				ID:    "user-1",
				Email: "test@example.com",
			},
			expectedUser: &types.User{
				ID:    "user-1",
				Email: "test@example.com",
			},
		},
		{
			name: "find by id",
			params: types.UserSearchParams{
				ID: stringPtr("user-1"),
			},
			mockUser: &types.User{
				ID:    "user-1",
				Email: "test@example.com",
			},
			expectedUser: &types.User{
				ID:    "user-1",
				Email: "test@example.com",
			},
		},
		{
			name: "user not found",
			params: types.UserSearchParams{
				Email: stringPtr("notfound@example.com"),
			},
			mockError:     errors.New("user not found"),
			expectedError: errors.New("user not found"),
		},
		{
			name: "database error",
			params: types.UserSearchParams{
				Email: stringPtr("test@example.com"),
			},
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockUserDB)
			api := NewApi(mockDB)

			mockDB.On("FindUser", tt.params).Return(tt.mockUser, tt.mockError)

			user, err := api.FindUser(tt.params)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func TestGetOrCreateUserFromAuthProvider(t *testing.T) {
	tests := []struct {
		name          string
		provider      string
		providerID    string
		email         string
		userName      string
		mockUser      *types.User
		mockError     error
		expectedUser  *types.User
		expectedError error
	}{
		{
			name:       "existing user",
			provider:   "auth0",
			providerID: "auth0|123",
			email:      "test@example.com",
			userName:   "Test User",
			mockUser: &types.User{
				ID:    "user-1",
				Email: "test@example.com",
				Name:  "Test User",
			},
			expectedUser: &types.User{
				ID:    "user-1",
				Email: "test@example.com",
				Name:  "Test User",
			},
		},
		{
			name:       "create new user",
			provider:   "auth0",
			providerID: "auth0|456",
			email:      "new@example.com",
			userName:   "New User",
			mockUser: &types.User{
				ID:    "user-2",
				Email: "new@example.com",
				Name:  "New User",
			},
			expectedUser: &types.User{
				ID:    "user-2",
				Email: "new@example.com",
				Name:  "New User",
			},
		},
		{
			name:          "database error",
			provider:      "auth0",
			providerID:    "auth0|789",
			email:         "error@example.com",
			userName:      "Error User",
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockUserDB)
			api := NewApi(mockDB)

			mockDB.On("GetOrCreateUserFromAuthProvider", tt.provider, tt.providerID, tt.email, tt.userName).Return(tt.mockUser, tt.mockError)

			user, err := api.GetOrCreateUserFromAuthProvider(tt.provider, tt.providerID, tt.email, tt.userName)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
