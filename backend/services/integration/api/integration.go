package api

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	liberrors "ems.dev/backend/libraries/errors"
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

// encryptToken encrypts the token using AES-GCM
func (a *Api) encryptToken(token string) (string, error) {
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

// DecryptToken decrypts the token using AES-GCM
func (a *Api) DecryptToken(encryptedToken string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedToken)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(a.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

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

// GetIntegrationConfig retrieves an integration config by ID
func (a *Api) GetIntegrationConfig(ctx context.Context, id string) (*types.IntegrationConfig, error) {
	return a.db.GetIntegrationConfig(id)
}

// GetOrganizationIntegrationConfigs retrieves all integration configs for an organization
func (a *Api) GetOrganizationIntegrationConfigs(ctx context.Context, orgID string) ([]types.IntegrationConfig, error) {
	return a.db.GetOrganizationIntegrationConfigs(orgID)
}

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

// DeleteIntegrationConfig deletes an integration config
func (a *Api) DeleteIntegrationConfig(ctx context.Context, id string) error {
	return a.db.DeleteIntegrationConfig(id)
}
