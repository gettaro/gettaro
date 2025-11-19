package api

import (
	"context"

	"ems.dev/backend/services/integration/types"
)

// UpdateIntegrationConfig updates an existing integration config
func (a *Api) UpdateIntegrationConfig(ctx context.Context, id string, req *types.UpdateIntegrationConfigRequest) (*types.IntegrationConfig, error) {
	config, err := a.db.GetIntegrationConfig(id)
	if err != nil {
		return nil, err
	}

	if req.Token != "" {
		encryptedToken, err := a.encryptToken(req.Token)
		if err != nil {
			return nil, err
		}
		config.EncryptedToken = encryptedToken
	}

	if req.Metadata != nil {
		config.Metadata = req.Metadata
	}

	if err := a.db.UpdateIntegrationConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}
