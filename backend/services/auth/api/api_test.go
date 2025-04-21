package api

import (
	"context"
	"errors"
	"testing"

	auth0client "ems.dev/backend/libraries/auth0"
	"ems.dev/backend/services/auth/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuth0Client is a mock implementation of the Auth0Client interface
type MockAuth0Client struct {
	mock.Mock
}

func (m *MockAuth0Client) GetUserInfo(ctx context.Context, accessToken string) (*auth0client.UserInfo, error) {
	args := m.Called(ctx, accessToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth0client.UserInfo), args.Error(1)
}

// MockAuthDB is a mock implementation of the AuthDB interface
type MockAuthDB struct {
	mock.Mock
}

func (m *MockAuthDB) GetExternalAuth(ctx context.Context, providerID string) (*types.AuthProvider, error) {
	args := m.Called(ctx, providerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.AuthProvider), args.Error(1)
}

func (m *MockAuthDB) CreateExternalAuth(ctx context.Context, authProvider *types.AuthProvider) error {
	args := m.Called(ctx, authProvider)
	return args.Error(0)
}

func TestGetUserInfo(t *testing.T) {
	tests := []struct {
		name          string
		accessToken   string
		mockUserInfo  *auth0client.UserInfo
		mockError     error
		expectedInfo  *auth0client.UserInfo
		expectedError error
	}{
		{
			name:        "successful retrieval",
			accessToken: "valid-token",
			mockUserInfo: &auth0client.UserInfo{
				Sub:      "auth0|123",
				Email:    "test@example.com",
				Name:     "Test User",
				Picture:  "https://example.com/picture.jpg",
				Provider: "auth0",
			},
			expectedInfo: &auth0client.UserInfo{
				Sub:      "auth0|123",
				Email:    "test@example.com",
				Name:     "Test User",
				Picture:  "https://example.com/picture.jpg",
				Provider: "auth0",
			},
		},
		{
			name:          "invalid token",
			accessToken:   "invalid-token",
			mockError:     errors.New("invalid token"),
			expectedError: errors.New("invalid token"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockAuth0Client)
			mockDB := new(MockAuthDB)
			api := NewApi(mockClient, mockDB)

			mockClient.On("GetUserInfo", context.Background(), tt.accessToken).Return(tt.mockUserInfo, tt.mockError)

			userInfo, err := api.GetUserInfo(context.Background(), tt.accessToken)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, userInfo)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedInfo, userInfo)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestGetExternalAuth(t *testing.T) {
	tests := []struct {
		name             string
		providerID       string
		mockAuthProvider *types.AuthProvider
		mockError        error
		expectedAuth     *types.AuthProvider
		expectedError    error
	}{
		{
			name:       "successful retrieval",
			providerID: "auth0|123",
			mockAuthProvider: &types.AuthProvider{
				ID:         "provider-1",
				UserID:     "user-1",
				Provider:   "auth0",
				ProviderID: "auth0|123",
			},
			expectedAuth: &types.AuthProvider{
				ID:         "provider-1",
				UserID:     "user-1",
				Provider:   "auth0",
				ProviderID: "auth0|123",
			},
		},
		{
			name:          "not found",
			providerID:    "auth0|456",
			mockError:     nil,
			expectedAuth:  nil,
			expectedError: nil,
		},
		{
			name:          "database error",
			providerID:    "auth0|789",
			mockError:     errors.New("database error"),
			expectedAuth:  nil,
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockAuth0Client)
			mockDB := new(MockAuthDB)
			api := NewApi(mockClient, mockDB)

			mockDB.On("GetExternalAuth", context.Background(), tt.providerID).Return(tt.mockAuthProvider, tt.mockError)

			authProvider, err := api.GetExternalAuth(context.Background(), tt.providerID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, authProvider)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAuth, authProvider)
			}

			mockDB.AssertExpectations(t)
		})
	}
}
