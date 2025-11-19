package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"ems.dev/backend/services/integration/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateIntegrationConfig(t *testing.T) {
	validKey := make([]byte, 32)
	for i := range validKey {
		validKey[i] = byte(i)
	}

	tests := []struct {
		name          string
		id            string
		req           *types.UpdateIntegrationConfigRequest
		mockConfig    *types.IntegrationConfig
		mockGetError  error
		mockUpdateError error
		expectedError error
		validateFunc  func(t *testing.T, config *types.IntegrationConfig)
	}{
		{
			name: "success - update token only",
			id:   "config-1",
			req: &types.UpdateIntegrationConfigRequest{
				Token: "new-token",
			},
			mockConfig: &types.IntegrationConfig{
				ID:             "config-1",
				OrganizationID: "org-1",
				ProviderName:   types.IntegrationProviderGithub,
				ProviderType:   types.IntegrationProviderTypeSourceControl,
				EncryptedToken: "old-encrypted-token",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			validateFunc: func(t *testing.T, config *types.IntegrationConfig) {
				assert.NotEqual(t, "old-encrypted-token", config.EncryptedToken)
				assert.NotEqual(t, "new-token", config.EncryptedToken) // Should be encrypted
				assert.NotEmpty(t, config.EncryptedToken)
			},
		},
		{
			name: "success - update metadata only",
			id:   "config-1",
			req: &types.UpdateIntegrationConfigRequest{
				Metadata: createMetadataJSON(`{"repositories": "repo1,repo2"}`),
			},
			mockConfig: &types.IntegrationConfig{
				ID:             "config-1",
				OrganizationID: "org-1",
				ProviderName:   types.IntegrationProviderGithub,
				ProviderType:   types.IntegrationProviderTypeSourceControl,
				EncryptedToken: "encrypted-token",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			validateFunc: func(t *testing.T, config *types.IntegrationConfig) {
				assert.Equal(t, "encrypted-token", config.EncryptedToken) // Should remain unchanged
				assert.NotNil(t, config.Metadata)
			},
		},
		{
			name: "success - update both token and metadata",
			id:   "config-1",
			req: &types.UpdateIntegrationConfigRequest{
				Token:    "new-token",
				Metadata: createMetadataJSON(`{"repositories": "repo1,repo2"}`),
			},
			mockConfig: &types.IntegrationConfig{
				ID:             "config-1",
				OrganizationID: "org-1",
				ProviderName:   types.IntegrationProviderGithub,
				ProviderType:   types.IntegrationProviderTypeSourceControl,
				EncryptedToken: "old-encrypted-token",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			validateFunc: func(t *testing.T, config *types.IntegrationConfig) {
				assert.NotEqual(t, "old-encrypted-token", config.EncryptedToken)
				assert.NotNil(t, config.Metadata)
			},
		},
		{
			name: "success - no updates (empty request)",
			id:   "config-1",
			req: &types.UpdateIntegrationConfigRequest{
				Token:    "",
				Metadata: nil,
			},
			mockConfig: &types.IntegrationConfig{
				ID:             "config-1",
				OrganizationID: "org-1",
				ProviderName:   types.IntegrationProviderGithub,
				ProviderType:   types.IntegrationProviderTypeSourceControl,
				EncryptedToken: "encrypted-token",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			validateFunc: func(t *testing.T, config *types.IntegrationConfig) {
				assert.Equal(t, "encrypted-token", config.EncryptedToken) // Should remain unchanged
			},
		},
		{
			name:          "error - config not found",
			id:            "non-existent",
			req:           &types.UpdateIntegrationConfigRequest{Token: "new-token"},
			mockGetError:  errors.New("record not found"),
			expectedError: errors.New("record not found"),
		},
		{
			name: "error - database update error",
			id:   "config-1",
			req: &types.UpdateIntegrationConfigRequest{
				Token: "new-token",
			},
			mockConfig: &types.IntegrationConfig{
				ID:             "config-1",
				OrganizationID: "org-1",
				ProviderName:   types.IntegrationProviderGithub,
				ProviderType:   types.IntegrationProviderTypeSourceControl,
				EncryptedToken: "old-encrypted-token",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			mockUpdateError: errors.New("database update failed"),
			expectedError:   errors.New("database update failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			api := NewApi(mockDB, validKey)

			mockDB.On("GetIntegrationConfig", tt.id).Return(tt.mockConfig, tt.mockGetError)

			if tt.mockGetError == nil {
				mockDB.On("UpdateIntegrationConfig", mock.AnythingOfType("*types.IntegrationConfig")).Return(tt.mockUpdateError).Run(func(args mock.Arguments) {
					config := args.Get(0).(*types.IntegrationConfig)
					if tt.validateFunc != nil && tt.mockUpdateError == nil {
						tt.validateFunc(t, config)
					}
				})
			}

			config, err := api.UpdateIntegrationConfig(context.Background(), tt.id, tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
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
