package api

import (
	"context"
	"errors"
	"testing"

	liberrors "ems.dev/backend/libraries/errors"
	"ems.dev/backend/services/integration/types"
	"gorm.io/datatypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDB is a mock implementation of the database.DB interface
type MockDB struct {
	mock.Mock
}

func (m *MockDB) CreateIntegrationConfig(config *types.IntegrationConfig) error {
	args := m.Called(config)
	return args.Error(0)
}

func (m *MockDB) GetIntegrationConfig(id string) (*types.IntegrationConfig, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.IntegrationConfig), args.Error(1)
}

func (m *MockDB) GetOrganizationIntegrationConfigs(orgID string) ([]types.IntegrationConfig, error) {
	args := m.Called(orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.IntegrationConfig), args.Error(1)
}

func (m *MockDB) UpdateIntegrationConfig(config *types.IntegrationConfig) error {
	args := m.Called(config)
	return args.Error(0)
}

func (m *MockDB) DeleteIntegrationConfig(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateIntegrationConfig(t *testing.T) {
	// Generate a valid AES-256 key (32 bytes)
	validKey := make([]byte, 32)
	for i := range validKey {
		validKey[i] = byte(i)
	}

	tests := []struct {
		name          string
		orgID         string
		req           *types.CreateIntegrationConfigRequest
		mockError     error
		expectedError error
		validateFunc  func(t *testing.T, config *types.IntegrationConfig)
	}{
		{
			name:  "success - github provider",
			orgID: "org-1",
			req: &types.CreateIntegrationConfigRequest{
				ProviderName: types.IntegrationProviderGithub,
				Token:        "test-token",
				Metadata:     createMetadataJSON(`{"repositories": "repo1,repo2"}`),
			},
			validateFunc: func(t *testing.T, config *types.IntegrationConfig) {
				assert.Equal(t, "org-1", config.OrganizationID)
				assert.Equal(t, types.IntegrationProviderGithub, config.ProviderName)
				assert.Equal(t, types.IntegrationProviderTypeSourceControl, config.ProviderType)
				assert.NotEmpty(t, config.EncryptedToken)
				assert.NotEqual(t, "test-token", config.EncryptedToken) // Should be encrypted
			},
		},
		{
			name:  "success - cursor provider",
			orgID: "org-1",
			req: &types.CreateIntegrationConfigRequest{
				ProviderName: types.IntegrationProviderCursor,
				Token:        "test-token",
				Metadata:     nil,
			},
			validateFunc: func(t *testing.T, config *types.IntegrationConfig) {
				assert.Equal(t, types.IntegrationProviderCursor, config.ProviderName)
				assert.Equal(t, types.IntegrationProviderTypeAICodeAssistant, config.ProviderType)
			},
		},
		{
			name:  "success - default provider type",
			orgID: "org-1",
			req: &types.CreateIntegrationConfigRequest{
				ProviderName: "unknown-provider",
				Token:        "test-token",
				Metadata:     createMetadataJSON(`{"repositories": "repo1"}`),
			},
			validateFunc: func(t *testing.T, config *types.IntegrationConfig) {
				assert.Equal(t, types.IntegrationProviderTypeSourceControl, config.ProviderType)
			},
		},
		{
			name:  "error - missing metadata for source control provider",
			orgID: "org-1",
			req: &types.CreateIntegrationConfigRequest{
				ProviderName: types.IntegrationProviderGithub,
				Token:        "test-token",
				Metadata:     nil,
			},
			expectedError: liberrors.NewBadRequestError("metadata is required for source control providers"),
		},
		{
			name:  "error - missing repositories in metadata",
			orgID: "org-1",
			req: &types.CreateIntegrationConfigRequest{
				ProviderName: types.IntegrationProviderGithub,
				Token:        "test-token",
				Metadata:     createMetadataJSON(`{}`),
			},
			expectedError: liberrors.NewBadRequestError("repositories is required for source control providers"),
		},
		{
			name:  "error - empty repositories string",
			orgID: "org-1",
			req: &types.CreateIntegrationConfigRequest{
				ProviderName: types.IntegrationProviderGithub,
				Token:        "test-token",
				Metadata:     createMetadataJSON(`{"repositories": ""}`),
			},
			expectedError: liberrors.NewBadRequestError("repositories is required for source control providers"),
		},
		{
			name:  "error - invalid metadata JSON",
			orgID: "org-1",
			req: &types.CreateIntegrationConfigRequest{
				ProviderName: types.IntegrationProviderGithub,
				Token:        "test-token",
				Metadata:     datatypes.JSON("invalid json"),
			},
			expectedError: errors.New("invalid character 'i' looking for beginning of value"),
		},
		{
			name:  "error - database error",
			orgID: "org-1",
			req: &types.CreateIntegrationConfigRequest{
				ProviderName: types.IntegrationProviderGithub,
				Token:        "test-token",
				Metadata:     createMetadataJSON(`{"repositories": "repo1"}`),
			},
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			api := NewApi(mockDB, validKey)

			if tt.mockError == nil && tt.expectedError == nil {
				mockDB.On("CreateIntegrationConfig", mock.AnythingOfType("*types.IntegrationConfig")).Return(nil).Run(func(args mock.Arguments) {
					config := args.Get(0).(*types.IntegrationConfig)
					// Store the config for validation
					if tt.validateFunc != nil {
						tt.validateFunc(t, config)
					}
				})
			} else if tt.mockError != nil {
				mockDB.On("CreateIntegrationConfig", mock.AnythingOfType("*types.IntegrationConfig")).Return(tt.mockError)
			}

			config, err := api.CreateIntegrationConfig(context.Background(), tt.orgID, tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if tt.expectedError.Error() != "invalid character 'i' looking for beginning of value" {
					assert.Equal(t, tt.expectedError.Error(), err.Error())
				} else {
					// For JSON errors, just check that an error occurred
					assert.NotNil(t, err)
				}
				assert.Nil(t, config)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, config)
				if tt.validateFunc != nil {
					tt.validateFunc(t, config)
				}
			}

			mockDB.AssertExpectations(t)
		})
	}
}

// Helper function to create datatypes.JSON from JSON string
func createMetadataJSON(jsonStr string) datatypes.JSON {
	return datatypes.JSON(jsonStr)
}
