package aicodeassistant

import (
	"context"
	"fmt"

	"ems.dev/backend/jobs/aicodeassistant/providers"
	intapi "ems.dev/backend/services/integration/api"
	"ems.dev/backend/services/integration/types"
	orgapi "ems.dev/backend/services/organization/api"
)

type SyncJob struct {
	integrationAPI  intapi.IntegrationAPI
	orgAPI          orgapi.OrganizationAPI
	providerFactory *providers.Factory
}

func NewSyncJob(integrationAPI intapi.IntegrationAPI, orgAPI orgapi.OrganizationAPI, providerFactory *providers.Factory) *SyncJob {
	return &SyncJob{
		integrationAPI:  integrationAPI,
		orgAPI:          orgAPI,
		providerFactory: providerFactory,
	}
}

func (j *SyncJob) Run(ctx context.Context) error {
	// Get all organizations
	orgs, err := j.orgAPI.GetOrganizations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get organizations: %w", err)
	}

	for _, org := range orgs {
		// Get all integrations for the organization
		integrations, err := j.integrationAPI.GetOrganizationIntegrationConfigs(ctx, org.ID)
		if err != nil {
			fmt.Printf("Failed to get integrations for org %s: %v\n", org.ID, err)
			continue
		}

		for _, integration := range integrations {
			// Only process AI code assistant integrations
			if integration.ProviderType != types.IntegrationProviderTypeAICodeAssistant {
				continue
			}

			// Get provider implementation
			provider, err := j.providerFactory.GetProvider(string(integration.ProviderName))
			if err != nil {
				fmt.Printf("Failed to create provider for %s: %v\n", integration.ProviderName, err)
				continue
			}

			// Sync usage data
			// TODO: Move this to run in a go routine
			if err := provider.SyncUsageData(ctx, &integration); err != nil {
				fmt.Printf("Failed to sync usage data for integration %s: %v\n", integration.ID, err)
				continue
			}
		}
	}

	return nil
}
