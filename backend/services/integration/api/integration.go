package api

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	"ems.dev/backend/services/integration/database"
	"ems.dev/backend/services/integration/types"
)

type API struct {
	db            database.DB
	encryptionKey []byte
}

func NewAPI(db database.DB, encryptionKey []byte) *API {
	return &API{
		db:            db,
		encryptionKey: encryptionKey,
	}
}

// encryptToken encrypts the token using AES-GCM
func (a *API) encryptToken(token string) (string, error) {
	block, err := aes.NewCipher(a.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(token), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// CreateIntegrationConfig creates a new integration config
func (a *API) CreateIntegrationConfig(ctx context.Context, orgID string, userID string, req *types.CreateIntegrationConfigRequest) (*types.IntegrationConfig, error) {
	// Encrypt the token
	encryptedToken, err := a.encryptToken(req.Token)
	if err != nil {
		return nil, err
	}

	config := &types.IntegrationConfig{
		OrganizationID: orgID,
		ProviderName:   req.ProviderName,
		ProviderType:   string(req.ProviderName), // For now, type is same as name
		EncryptedToken: encryptedToken,
		Metadata:       req.Metadata,
	}

	if err := a.db.CreateIntegrationConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// GetIntegrationConfig retrieves an integration config by ID
func (a *API) GetIntegrationConfig(ctx context.Context, id string, userID string) (*types.IntegrationConfig, error) {
	return a.db.GetIntegrationConfig(id)
}

// GetOrganizationIntegrationConfigs retrieves all integration configs for an organization
func (a *API) GetOrganizationIntegrationConfigs(ctx context.Context, orgID string, userID string) ([]types.IntegrationConfig, error) {
	return a.db.GetOrganizationIntegrationConfigs(orgID)
}

// UpdateIntegrationConfig updates an existing integration config
func (a *API) UpdateIntegrationConfig(ctx context.Context, id string, userID string, req *types.UpdateIntegrationConfigRequest) (*types.IntegrationConfig, error) {
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

// DeleteIntegrationConfig deletes an integration config
func (a *API) DeleteIntegrationConfig(ctx context.Context, id string, userID string) error {
	return a.db.DeleteIntegrationConfig(id)
}
