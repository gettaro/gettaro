package cursor

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"ems.dev/backend/libraries/cursor/types"
)

// CursorClient defines the interface for Cursor API operations
type CursorClient interface {
	// GetTeamMembers retrieves all team members
	// GET /teams/members
	GetTeamMembers(ctx context.Context, apiKey string) (*types.TeamMembersResponse, error)

	// GetDailyUsageData retrieves daily usage metrics
	// POST /teams/daily-usage-data
	GetDailyUsageData(ctx context.Context, apiKey string, params *types.DailyUsageDataParams) (*types.DailyUsageDataResponse, error)

	// GetFilteredUsageEvents retrieves detailed usage events with filtering
	// POST /teams/filtered-usage-events
	GetFilteredUsageEvents(ctx context.Context, apiKey string, params *types.FilteredUsageEventsParams) (*types.FilteredUsageEventsResponse, error)
}

// Client represents a Cursor API client
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Cursor API client
func NewClient() *Client {
	return &Client{
		baseURL:    "https://api.cursor.com",
		httpClient: &http.Client{},
	}
}

// getAuthHeader creates a Basic Auth header with API key as username
func (c *Client) getAuthHeader(apiKey string) string {
	auth := apiKey + ":"
	encoded := base64.StdEncoding.EncodeToString([]byte(auth))
	return "Basic " + encoded
}

// GetTeamMembers fetches team members from Cursor Admin API
func (c *Client) GetTeamMembers(ctx context.Context, apiKey string) (*types.TeamMembersResponse, error) {
	url := fmt.Sprintf("%s/teams/members", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", c.getAuthHeader(apiKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var data types.TeamMembersResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &data, nil
}

// GetDailyUsageData fetches daily usage data from Cursor Admin API
// POST /teams/daily-usage-data
func (c *Client) GetDailyUsageData(ctx context.Context, apiKey string, params *types.DailyUsageDataParams) (*types.DailyUsageDataResponse, error) {
	url := fmt.Sprintf("%s/teams/daily-usage-data", c.baseURL)

	// Build request body - use empty object if params are nil or empty
	var reqBody []byte
	var err error
	if params == nil || (params.StartDate == nil && params.EndDate == nil) {
		// Send empty JSON object if no params provided
		reqBody = []byte("{}")
	} else {
		reqBody, err = json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", c.getAuthHeader(apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var data types.DailyUsageDataResponse

	// Try to unmarshal as an object with a "data" field first
	if err := json.Unmarshal(bodyBytes, &data); err == nil {
		// Check if we got data (either as an object with Data field, or empty but valid structure)
		if len(data.Data) > 0 {
			return &data, nil
		}
		// If data field exists but is empty, still return it (valid response)
		return &data, nil
	}

	// If that fails, try to unmarshal as a direct array (response might be just an array)
	var entries []types.DailyUsageEntry
	if err := json.Unmarshal(bodyBytes, &entries); err == nil {
		data.Data = entries
		return &data, nil
	}

	// If both fail, try to decode as generic map for debugging
	var rawData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &rawData); err == nil {
		data.RawData = rawData
		return &data, fmt.Errorf("unexpected response format, raw data stored for debugging")
	}

	return nil, fmt.Errorf("failed to decode response: expected array or object with 'data' field, got: %s", string(bodyBytes))
}

// GetFilteredUsageEvents fetches filtered usage events from Cursor Analytics API
func (c *Client) GetFilteredUsageEvents(ctx context.Context, apiKey string, params *types.FilteredUsageEventsParams) (*types.FilteredUsageEventsResponse, error) {
	url := fmt.Sprintf("%s/teams/filtered-usage-events", c.baseURL)

	reqBody, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", c.getAuthHeader(apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var data types.FilteredUsageEventsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &data, nil
}
