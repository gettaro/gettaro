package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"ems.dev/backend/services/integration/types"
	"github.com/stretchr/testify/assert"
)

func TestGetIntegrationConfig(t *testing.T) {
	validKey := make([]byte, 32)
	for i := range validKey {
		validKey[i] = byte(i)
	}

	tests := []struct {
		name          string
		id            string
		mockConfig    *types.IntegrationConfig
		mockError     error
		expectedConfig *types.IntegrationConfig
		expectedError error
	}{
		{
			name: "success",
			id:   "config-1",
			mockConfig: &types.IntegrationConfig{
				ID:             "config-1",
				OrganizationID: "org-1",
				ProviderName:   types.IntegrationProviderGithub,
				ProviderType:   types.IntegrationProviderTypeSourceControl,
				EncryptedToken: "encrypted-token",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			expectedConfig: &types.IntegrationConfig{
				ID:             "config-1",
				OrganizationID: "org-1",
				ProviderName:   types.IntegrationProviderGithub,
				ProviderType:   types.IntegrationProviderTypeSourceControl,
				EncryptedToken: "encrypted-token",
			},
		},
		{
			name:          "error - config not found",
			id:            "non-existent",
			mockError:     errors.New("record not found"),
			expectedError: errors.New("record not found"),
		},
		{
			name:          "error - database error",
			id:            "config-1",
			mockError:     errors.New("database connection failed"),
			expectedError: errors.New("database connection failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			api := NewApi(mockDB, validKey)

			mockDB.On("GetIntegrationConfig", tt.id).Return(tt.mockConfig, tt.mockError)

			config, err := api.GetIntegrationConfig(context.Background(), tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, config)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, config)
				assert.Equal(t, tt.expectedConfig.ID, config.ID)
				assert.Equal(t, tt.expectedConfig.OrganizationID, config.OrganizationID)
				assert.Equal(t, tt.expectedConfig.ProviderName, config.ProviderName)
				assert.Equal(t, tt.expectedConfig.ProviderType, config.ProviderType)
				assert.Equal(t, tt.expectedConfig.EncryptedToken, config.EncryptedToken)
			}

			mockDB.AssertExpectations(t)
		})
	}
}
