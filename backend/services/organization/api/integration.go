package api

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	"ems.dev/backend/services/errors"
	"ems.dev/backend/services/organization/database"
	"ems.dev/backend/services/organization/types"
)

type IntegrationAPI struct {
	db            database.IntegrationDB
	orgDB         database.DB
	encryptionKey []byte
}

func NewIntegrationAPI(db database.IntegrationDB, orgDB database.DB, encryptionKey []byte) *IntegrationAPI {
	return &IntegrationAPI{
		db:            db,
		orgDB:         orgDB,
		encryptionKey: encryptionKey,
	}
}

// encryptToken encrypts the token using AES-GCM
func (a *IntegrationAPI) encryptToken(token string) (string, error) {
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
func (a *IntegrationAPI) CreateIntegrationConfig(ctx context.Context, orgID string, userID string, req *types.CreateIntegrationConfigRequest) (*types.IntegrationConfig, error) {
	// Check if user is organization owner
	isOwner, err := a.orgDB.IsOrganizationOwner(orgID, userID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, &errors.ErrForbidden{Message: "only organization owners can create integration configs"}
	}

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
func (a *IntegrationAPI) GetIntegrationConfig(ctx context.Context, id string, userID string) (*types.IntegrationConfig, error) {
	config, err := a.db.GetIntegrationConfig(id)
	if err != nil {
		return nil, err
	}

	// Check if user has access to the organization
	orgs, err := a.orgDB.GetUserOrganizations(userID)
	if err != nil {
		return nil, err
	}

	hasAccess := false
	for _, org := range orgs {
		if org.ID == config.OrganizationID {
			hasAccess = true
			break
		}
	}

	if !hasAccess {
		return nil, &errors.ErrForbidden{Message: "user does not have access to this integration config"}
	}

	return config, nil
}

// GetOrganizationIntegrationConfigs retrieves all integration configs for an organization
func (a *IntegrationAPI) GetOrganizationIntegrationConfigs(ctx context.Context, orgID string, userID string) ([]types.IntegrationConfig, error) {
	// Check if user has access to the organization
	orgs, err := a.orgDB.GetUserOrganizations(userID)
	if err != nil {
		return nil, err
	}

	hasAccess := false
	for _, org := range orgs {
		if org.ID == orgID {
			hasAccess = true
			break
		}
	}

	if !hasAccess {
		return nil, &errors.ErrForbidden{Message: "user does not have access to this organization"}
	}

	return a.db.GetOrganizationIntegrationConfigs(orgID)
}

// UpdateIntegrationConfig updates an existing integration config
func (a *IntegrationAPI) UpdateIntegrationConfig(ctx context.Context, id string, userID string, req *types.UpdateIntegrationConfigRequest) (*types.IntegrationConfig, error) {
	config, err := a.db.GetIntegrationConfig(id)
	if err != nil {
		return nil, err
	}

	// Check if user is organization owner
	isOwner, err := a.orgDB.IsOrganizationOwner(config.OrganizationID, userID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, &errors.ErrForbidden{Message: "only organization owners can update integration configs"}
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
func (a *IntegrationAPI) DeleteIntegrationConfig(ctx context.Context, id string, userID string) error {
	config, err := a.db.GetIntegrationConfig(id)
	if err != nil {
		return err
	}

	// Check if user is organization owner
	isOwner, err := a.orgDB.IsOrganizationOwner(config.OrganizationID, userID)
	if err != nil {
		return err
	}
	if !isOwner {
		return &errors.ErrForbidden{Message: "only organization owners can delete integration configs"}
	}

	return a.db.DeleteIntegrationConfig(id)
}
