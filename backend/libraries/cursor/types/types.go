package types

// DailyUsageDataParams represents the request parameters for POST /teams/daily-usage-data
type DailyUsageDataParams struct {
	StartDate *int64 `json:"startDate,omitempty"` // Unix timestamp in milliseconds
	EndDate  *int64 `json:"endDate,omitempty"`   // Unix timestamp in milliseconds
}

// DailyUsageDataResponse represents the response from POST /teams/daily-usage-data
// Based on Cursor Admin API documentation: https://cursor.com/docs/account/teams/admin-api#get-daily-usage-data
// The response is an array of daily usage entries, one per user per day
type DailyUsageDataResponse struct {
	Data []DailyUsageEntry `json:"data,omitempty"`
	// Raw data if structured parsing fails
	RawData map[string]interface{} `json:"-"`
}

// DailyUsageEntry represents daily usage data for a specific user on a specific date
// Each entry in the response array is already user-specific
type DailyUsageEntry struct {
	Date                      int64  `json:"date"`                        // Date in epoch milliseconds
	Email                     string `json:"email,omitempty"`              // User's email address
	IsActive                  bool   `json:"isActive"`                    // Whether user was active on this day
	TotalLinesAdded           int    `json:"totalLinesAdded"`             // Total lines of code added
	TotalLinesDeleted         int    `json:"totalLinesDeleted"`           // Total lines of code deleted
	AcceptedLinesAdded        int    `json:"acceptedLinesAdded"`          // Lines added from accepted AI suggestions
	AcceptedLinesDeleted      int    `json:"acceptedLinesDeleted"`        // Lines deleted from accepted AI suggestions
	TotalApplies              int    `json:"totalApplies"`                // Number of apply operations
	TotalAccepts               int    `json:"totalAccepts"`               // Number of accepted suggestions
	TotalRejects               int    `json:"totalRejects"`               // Number of rejected suggestions
	TotalTabsShown             int    `json:"totalTabsShown"`             // Number of tab completions shown
	TotalTabsAccepted          int    `json:"totalTabsAccepted"`          // Number of tab completions accepted
	ComposerRequests           int    `json:"composerRequests"`           // Number of composer requests
	ChatRequests              int    `json:"chatRequests"`                // Number of chat requests
	AgentRequests             int    `json:"agentRequests"`               // Number of agent requests
	CmdkUsages                int    `json:"cmdkUsages"`                  // Number of command palette (Cmd+K) uses
	SubscriptionIncludedReqs  int    `json:"subscriptionIncludedReqs"`    // Number of subscription requests
	APIKeyReqs                int    `json:"apiKeyReqs"`                  // Number of API key requests
	UsageBasedReqs            int    `json:"usageBasedReqs"`              // Number of pay-per-use requests
	BugbotUsages              int    `json:"bugbotUsages"`                // Number of bug detection uses
	MostUsedModel             string `json:"mostUsedModel,omitempty"`     // Most frequently used AI model
	ApplyMostUsedExtension    string `json:"applyMostUsedExtension,omitempty"` // Most used file extension for applies (optional)
	TabMostUsedExtension      string `json:"tabMostUsedExtension,omitempty"`   // Most used file extension for tabs (optional)
	ClientVersion             string `json:"clientVersion,omitempty"`     // Cursor client version (optional)
}

// FilteredUsageEventsParams for POST /teams/filtered-usage-events
type FilteredUsageEventsParams struct {
	StartDate *int64 `json:"startDate,omitempty"` // Unix timestamp in milliseconds
	EndDate   *int64 `json:"endDate,omitempty"`   // Unix timestamp in milliseconds
	UserID    *int   `json:"userId,omitempty"`
	// Add other filter parameters as needed based on API docs
}

// FilteredUsageEventsResponse represents individual API call events
type FilteredUsageEventsResponse struct {
	Events []UsageEvent `json:"events"`
	// May include pagination fields if API supports it
}

// UsageEvent represents a single usage event
type UsageEvent struct {
	UserID          int     `json:"userId"`
	Timestamp       int64   `json:"timestamp"`        // Unix timestamp in milliseconds
	Model           string  `json:"model"`            // e.g., "claude-3-opus", "gpt-4"
	PromptTokens    int     `json:"promptTokens"`
	CompletionTokens int   `json:"completionTokens"`
	TotalTokens     int     `json:"totalTokens"`
	Cost            float64 `json:"cost,omitempty"`
	UsageType       string  `json:"usageType"`        // e.g., "completion", "chat", "edit", "write", "notebook_edit"
	SessionID       string  `json:"sessionId,omitempty"`
	LinesAccepted   int     `json:"linesAccepted,omitempty"`   // Lines of code accepted
	Accepted        bool    `json:"accepted,omitempty"`        // Whether suggestion was accepted
	// Additional fields based on actual API response
}

// CursorUser represents a user from Cursor API
type CursorUser struct {
	ID       int    `json:"id"`
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	// Other user fields
}

// TeamMember represents a team member from the GET /teams/members endpoint
type TeamMember struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"` // "owner" | "member" | "free-owner"
}

// TeamMembersResponse represents the response from GET /teams/members
type TeamMembersResponse struct {
	TeamMembers []TeamMember `json:"teamMembers"`
}


