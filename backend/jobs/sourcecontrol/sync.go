package sourcecontrol

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"ems.dev/backend/jobs/sourcecontrol/providers"
	intapi "ems.dev/backend/services/integration/api"
	"ems.dev/backend/services/integration/types"
	orgapi "ems.dev/backend/services/organization/api"
)

type SyncJob struct {
	integrationAPI  intapi.IntegrationAPI
	orgAPI          orgapi.OrganizationAPI
	providerFactory providers.ProviderFactory
}

func NewSyncJob(integrationAPI intapi.IntegrationAPI, orgAPI orgapi.OrganizationAPI, providerFactory providers.ProviderFactory) *SyncJob {
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
			// Only process source control integrations
			if integration.ProviderType != types.IntegrationProviderTypeSourceControl {
				continue
			}

			// Get provider implementation
			provider, err := j.providerFactory.GetProvider(string(integration.ProviderName))
			if err != nil {
				fmt.Printf("Failed to create provider for %s: %v\n", integration.ProviderName, err)
				continue
			}

			// Parse repositories from metadata
			var metadata map[string]interface{}
			if err := json.Unmarshal(integration.Metadata, &metadata); err != nil {
				fmt.Printf("Failed to parse metadata for integration %s: %v\n", integration.ID, err)
				continue
			}

			reposStr, ok := metadata["repositories"].(string)
			if !ok || reposStr == "" {
				fmt.Printf("No repositories found in metadata for integration %s\n", integration.ID)
				continue
			}

			repositories := strings.Split(reposStr, ",")
			for i := range repositories {
				repositories[i] = strings.TrimSpace(repositories[i])
			}

			// Sync repositories
			if err := provider.SyncRepositories(ctx, &integration, repositories); err != nil {
				fmt.Printf("Failed to sync repositories for integration %s: %v\n", integration.ID, err)
				continue
			}
		}
	}

	return nil
}
