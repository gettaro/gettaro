package auth0

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Auth0Client defines the interface for Auth0-related operations
type Auth0Client interface {
	// GetUserInfo retrieves user information from Auth0's userinfo endpoint
	GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error)
}

// UserInfo represents the user information returned by Auth0
type UserInfo struct {
	Sub      string `json:"sub"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	Provider string `json:"provider"`
}

// Client represents an Auth0 client
type Client struct {
	domain   string
	clientID string
	client   *http.Client
}

// NewClient creates a new Auth0 client
func NewClient(issuerUrl string, clientID string) *Client {
	return &Client{
		domain:   issuerUrl,
		clientID: clientID,
		client:   &http.Client{},
	}
}

// GetUserInfo retrieves user information from Auth0's userinfo endpoint
func (c *Client) GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%suserinfo", c.domain), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &userInfo, nil
}
