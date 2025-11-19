package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"ems.dev/backend/services/integration/types"
	"github.com/stretchr/testify/assert"
)

func TestGetOrganizationIntegrationConfigs(t *testing.T) {
	validKey := make([]byte, 32)
	for i := range validKey {
		validKey[i] = byte(i)
	}

	tests := []struct {
		name           string
		orgID          string
		mockConfigs    []types.IntegrationConfig
		mockError      error
		expectedConfigs []types.IntegrationConfig
		expectedError  error
	}{
		{
			name:  "success - multiple configs",
			orgID: "org-1",
			mockConfigs: []types.IntegrationConfig{
				{
					ID:             "config-1",
					OrganizationID: "org-1",
					ProviderName:   types.IntegrationProviderGithub,
					ProviderType:   types.IntegrationProviderTypeSourceControl,
					EncryptedToken: "encrypted-token-1",
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				},
				{
					ID:             "config-2",
					OrganizationID: "org-1",
					ProviderName:   types.IntegrationProviderCursor,
					ProviderType:   types.IntegrationProviderTypeAICodeAssistant,
					EncryptedToken: "encrypted-token-2",
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				},
			},
			expectedConfigs: []types.IntegrationConfig{
				{
					ID:             "config-1",
					OrganizationID: "org-1",
					ProviderName:   types.IntegrationProviderGithub,
					ProviderType:   types.IntegrationProviderTypeSourceControl,
					EncryptedToken: "encrypted-token-1",
				},
				{
					ID:             "config-2",
					OrganizationID: "org-1",
					ProviderName:   types.IntegrationProviderCursor,
					ProviderType:   types.IntegrationProviderTypeAICodeAssistant,
					EncryptedToken: "encrypted-token-2",
				},
			},
		},
		{
			name:           "success - empty list",
			orgID:          "org-1",
			mockConfigs:    []types.IntegrationConfig{},
			expectedConfigs: []types.IntegrationConfig{},
		},
		{
			name:          "error - database error",
			orgID:         "org-1",
			mockError:     errors.New("database connection failed"),
			expectedError: errors.New("database connection failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			api := NewApi(mockDB, validKey)

			mockDB.On("GetOrganizationIntegrationConfigs", tt.orgID).Return(tt.mockConfigs, tt.mockError)

			configs, err := api.GetOrganizationIntegrationConfigs(context.Background(), tt.orgID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, configs)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, configs)
				assert.Equal(t, len(tt.expectedConfigs), len(configs))
				for i, expected := range tt.expectedConfigs {
					assert.Equal(t, expected.ID, configs[i].ID)
					assert.Equal(t, expected.OrganizationID, configs[i].OrganizationID)
					assert.Equal(t, expected.ProviderName, configs[i].ProviderName)
					assert.Equal(t, expected.ProviderType, configs[i].ProviderType)
					assert.Equal(t, expected.EncryptedToken, configs[i].EncryptedToken)
				}
			}

			mockDB.AssertExpectations(t)
		})
	}
}
