package api

import (
	"context"
	"encoding/json"

	liberrors "ems.dev/backend/libraries/errors"
	"ems.dev/backend/services/integration/types"
)

// CreateIntegrationConfig creates a new integration config
func (a *Api) CreateIntegrationConfig(ctx context.Context, orgID string, req *types.CreateIntegrationConfigRequest) (*types.IntegrationConfig, error) {
	// Encrypt the token
	encryptedToken, err := a.encryptToken(req.Token)
	if err != nil {
		return nil, err
	}

	// Determine provider type based on provider name
	var providerType types.IntegrationProviderType
	switch req.ProviderName {
	case "github":
		providerType = types.IntegrationProviderTypeSourceControl
	case "cursor":
		providerType = types.IntegrationProviderTypeAICodeAssistant
	default:
		// Default to source control for backward compatibility
		providerType = types.IntegrationProviderTypeSourceControl
	}

	// Validate repositories metadata for source control providers only
	if providerType == types.IntegrationProviderTypeSourceControl {
		if req.Metadata == nil {
			return nil, liberrors.NewBadRequestError("metadata is required for source control providers")
		}
		var metadata map[string]interface{}
		if err := json.Unmarshal(req.Metadata, &metadata); err != nil {
			return nil, err
		}

		repos, ok := metadata["repositories"].(string)
		if !ok || repos == "" {
			return nil, liberrors.NewBadRequestError("repositories is required for source control providers")
		}
	}

	config := &types.IntegrationConfig{
		OrganizationID: orgID,
		ProviderName:   req.ProviderName,
		ProviderType:   providerType,
		EncryptedToken: encryptedToken,
		Metadata:       req.Metadata,
	}

	if err := a.db.CreateIntegrationConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}
