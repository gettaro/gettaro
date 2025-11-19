package api

import (
	"context"

	"ems.dev/backend/services/integration/database"
	"ems.dev/backend/services/integration/types"
)

type IntegrationAPI interface {
	CreateIntegrationConfig(ctx context.Context, orgID string, req *types.CreateIntegrationConfigRequest) (*types.IntegrationConfig, error)
	GetIntegrationConfig(ctx context.Context, id string) (*types.IntegrationConfig, error)
	GetOrganizationIntegrationConfigs(ctx context.Context, orgID string) ([]types.IntegrationConfig, error)
	UpdateIntegrationConfig(ctx context.Context, id string, req *types.UpdateIntegrationConfigRequest) (*types.IntegrationConfig, error)
	DeleteIntegrationConfig(ctx context.Context, id string) error
	DecryptToken(encryptedToken string) (string, error)
}

type Api struct {
	db            database.DB
	encryptionKey []byte
}

func NewApi(db database.DB, encryptionKey []byte) *Api {
	return &Api{
		db:            db,
		encryptionKey: encryptionKey,
	}
}
