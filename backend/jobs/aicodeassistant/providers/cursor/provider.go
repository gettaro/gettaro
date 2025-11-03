package cursor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"ems.dev/backend/libraries/cursor"
	cursortypes "ems.dev/backend/libraries/cursor/types"
	aicodeassistantapi "ems.dev/backend/services/aicodeassistant/api"
	aicodeassistanttypes "ems.dev/backend/services/aicodeassistant/types"
	"ems.dev/backend/services/integration/api"
	inttypes "ems.dev/backend/services/integration/types"
	memberapi "ems.dev/backend/services/member/api"
	membertypes "ems.dev/backend/services/member/types"
	"gorm.io/datatypes"
)

type CursorProvider struct {
	cursorClient       cursor.CursorClient
	integrationAPI     api.IntegrationAPI
	aiCodeAssistantAPI aicodeassistantapi.AICodeAssistantAPI
	memberAPI          memberapi.MemberAPI
}

func NewProvider(
	cursorClient cursor.CursorClient,
	integrationAPI api.IntegrationAPI,
	aiCodeAssistantAPI aicodeassistantapi.AICodeAssistantAPI,
	memberAPI memberapi.MemberAPI,
) *CursorProvider {
	return &CursorProvider{
		cursorClient:       cursorClient,
		integrationAPI:     integrationAPI,
		aiCodeAssistantAPI: aiCodeAssistantAPI,
		memberAPI:          memberAPI,
	}
}

func (p *CursorProvider) Name() string {
	return "cursor"
}

// SyncUsageData syncs usage data from Cursor API
func (p *CursorProvider) SyncUsageData(ctx context.Context, config *inttypes.IntegrationConfig) error {
	// 1. Decrypt API key from config.EncryptedToken
	apiKey, err := p.integrationAPI.DecryptToken(config.EncryptedToken)
	if err != nil {
		return fmt.Errorf("failed to decrypt API key: %w", err)
	}

	// 2. First, sync all team members
	teamMembers, err := p.cursorClient.GetTeamMembers(ctx, apiKey)
	if err != nil {
		return fmt.Errorf("failed to fetch team members: %w", err)
	}

	// Create external accounts for all team members
	emailToAccountID := make(map[string]string)
	for _, member := range teamMembers.TeamMembers {
		account, err := p.upsertCursorAccountFromTeamMember(ctx, config.OrganizationID, member)
		if err != nil {
			fmt.Printf("Warning: failed to upsert account for team member %s (%s): %v\n", member.Name, member.Email, err)
			continue
		}
		emailToAccountID[member.Email] = account.ID
	}

	// 3. Fetch daily usage data from Cursor Admin API
	// Use date range based on last sync time or default to last 30 days
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30) // Default to 30 days

	if config.LastSyncedAt != nil {
		startDate = *config.LastSyncedAt
	}

	// Build params for daily usage data endpoint
	dailyUsageParams := &cursortypes.DailyUsageDataParams{
		StartDate: intPtr(startDate.UnixMilli()),
		EndDate:   intPtr(endDate.UnixMilli()),
	}

	// Fetch daily usage data
	dailyUsageData, err := p.cursorClient.GetDailyUsageData(ctx, apiKey, dailyUsageParams)
	if err != nil {
		return fmt.Errorf("failed to fetch daily usage data: %w", err)
	}

	// 4. Process daily usage data and create daily metrics
	// Each entry in the Data array is already user-specific (one entry per user per day)
	userDateMetrics := make(map[string]*aicodeassistanttypes.AICodeAssistantDailyMetric)

	// Process daily usage entries
	for _, entry := range dailyUsageData.Data {
		// Skip entries without email (shouldn't happen, but be safe)
		if entry.Email == "" {
			fmt.Printf("Warning: entry without email, skipping\n")
			continue
		}

		// Find account ID for this email
		accountID, exists := emailToAccountID[entry.Email]
		if !exists {
			fmt.Printf("Warning: no account found for email %s\n", entry.Email)
			continue
		}

		// Convert date from epoch milliseconds to time.Time
		metricDate := time.Unix(0, entry.Date*int64(time.Millisecond)).Truncate(24 * time.Hour)
		dateKey := metricDate.Format("2006-01-02")
		key := fmt.Sprintf("%s_%s", accountID, dateKey)

		metric, exists := userDateMetrics[key]
		if !exists {
			metric = &aicodeassistanttypes.AICodeAssistantDailyMetric{
				OrganizationID:       config.OrganizationID,
				ExternalAccountID:    accountID,
				ToolName:             "cursor",
				MetricDate:           metricDate,
				LinesOfCodeAccepted:  0,
				LinesOfCodeSuggested: 0,
				ActiveSessions:       0,
			}
			userDateMetrics[key] = metric
		}

		// Map Cursor API fields to our metrics
		// Lines accepted = acceptedLinesAdded (lines from accepted AI suggestions)
		metric.LinesOfCodeAccepted += entry.AcceptedLinesAdded + entry.AcceptedLinesDeleted

		// Lines suggested = totalLinesAdded + totalLinesDeleted
		linesSuggested := entry.TotalLinesAdded + entry.TotalLinesDeleted
		metric.LinesOfCodeSuggested += linesSuggested

		// Suggestion accept rate = (acceptedLinesAdded + acceptedLinesDeleted) / (totalLinesAdded + totalLinesDeleted)
		acceptRate := float64(entry.AcceptedLinesAdded+entry.AcceptedLinesDeleted) / float64(entry.TotalLinesAdded+entry.TotalLinesDeleted) * 100
		metric.SuggestionAcceptRate = &acceptRate

		// Active sessions: use isActive as indicator (we could also track unique session IDs if available)
		if entry.IsActive {
			metric.ActiveSessions = 1
		}

		// Build metadata with all entry fields (all requests and metrics)
		entryMetadata := map[string]interface{}{
			"date":                     entry.Date,
			"email":                    entry.Email,
			"isActive":                 entry.IsActive,
			"totalLinesAdded":          entry.TotalLinesAdded,
			"totalLinesDeleted":        entry.TotalLinesDeleted,
			"acceptedLinesAdded":       entry.AcceptedLinesAdded,
			"acceptedLinesDeleted":     entry.AcceptedLinesDeleted,
			"totalApplies":             entry.TotalApplies,
			"totalAccepts":             entry.TotalAccepts,
			"totalRejects":             entry.TotalRejects,
			"totalTabsShown":           entry.TotalTabsShown,
			"totalTabsAccepted":        entry.TotalTabsAccepted,
			"composerRequests":         entry.ComposerRequests,
			"chatRequests":             entry.ChatRequests,
			"agentRequests":            entry.AgentRequests,
			"cmdkUsages":               entry.CmdkUsages,
			"subscriptionIncludedReqs": entry.SubscriptionIncludedReqs,
			"apiKeyReqs":               entry.APIKeyReqs,
			"usageBasedReqs":           entry.UsageBasedReqs,
			"bugbotUsages":             entry.BugbotUsages,
		}

		// Add optional fields if present
		if entry.MostUsedModel != "" {
			entryMetadata["mostUsedModel"] = entry.MostUsedModel
		}
		if entry.ApplyMostUsedExtension != "" {
			entryMetadata["applyMostUsedExtension"] = entry.ApplyMostUsedExtension
		}
		if entry.TabMostUsedExtension != "" {
			entryMetadata["tabMostUsedExtension"] = entry.TabMostUsedExtension
		}
		if entry.ClientVersion != "" {
			entryMetadata["clientVersion"] = entry.ClientVersion
		}

		// Store metadata in the metric
		metadataBytes, _ := json.Marshal(entryMetadata)
		metric.Metadata = datatypes.JSON(metadataBytes)
	}

	// 5. Save daily metrics
	for _, metric := range userDateMetrics {
		// Create or update daily metric
		if err := p.aiCodeAssistantAPI.CreateOrUpdateDailyMetric(ctx, metric); err != nil {
			fmt.Printf("Warning: failed to create/update daily metric for user %s on %s: %v\n",
				metric.ExternalAccountID, metric.MetricDate.Format("2006-01-02"), err)
			continue
		}
	}

	return nil
}

// upsertCursorAccountFromTeamMember creates or retrieves an external account from a TeamMember
func (p *CursorProvider) upsertCursorAccountFromTeamMember(ctx context.Context, organizationID string, member cursortypes.TeamMember) (*membertypes.ExternalAccount, error) {
	// Check if account exists by email (which serves as provider_id for team members)
	accountType := "ai-code-assistant"
	accounts, err := p.memberAPI.GetExternalAccounts(ctx, &membertypes.ExternalAccountParams{
		OrganizationID: organizationID,
		AccountType:    &accountType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch external accounts: %w", err)
	}

	// Find existing account by email (used as provider_id for team members)
	for _, account := range accounts {
		if account.ProviderID == member.Email && account.ProviderName == "cursor" {
			return &account, nil
		}
	}

	// Create new account if not found
	metadata := map[string]interface{}{
		"cursor_email": member.Email,
		"cursor_role":  member.Role,
	}
	metadataBytes, _ := json.Marshal(metadata)

	newAccount := &membertypes.ExternalAccount{
		AccountType:    "ai-code-assistant",
		ProviderName:   "cursor",
		ProviderID:     member.Email, // Use email as provider_id for team members
		Username:       member.Name,
		OrganizationID: &organizationID,
		Metadata:       datatypes.JSON(metadataBytes),
	}

	if err := p.memberAPI.CreateExternalAccounts(ctx, []*membertypes.ExternalAccount{newAccount}); err != nil {
		return nil, fmt.Errorf("failed to create external account: %w", err)
	}

	// Fetch the created account to return it with ID
	accounts, err = p.memberAPI.GetExternalAccounts(ctx, &membertypes.ExternalAccountParams{
		OrganizationID: organizationID,
		AccountType:    &accountType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created account: %w", err)
	}

	for _, account := range accounts {
		if account.ProviderID == member.Email && account.ProviderName == "cursor" {
			return &account, nil
		}
	}

	return nil, fmt.Errorf("failed to find created account")
}

// upsertCursorAccount creates or retrieves an external account for a Cursor user
func (p *CursorProvider) upsertCursorAccount(ctx context.Context, organizationID string, cursorUser cursortypes.CursorUser) (*membertypes.ExternalAccount, error) {
	// Check if account exists by provider_id and provider_name
	accountType := "ai-code-assistant"
	accounts, err := p.memberAPI.GetExternalAccounts(ctx, &membertypes.ExternalAccountParams{
		OrganizationID: organizationID,
		AccountType:    &accountType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch external accounts: %w", err)
	}

	// Find existing account by provider ID
	providerID := fmt.Sprintf("%d", cursorUser.ID)
	for _, account := range accounts {
		if account.ProviderID == providerID && account.ProviderName == "cursor" {
			return &account, nil
		}
	}

	// Create new account if not found
	metadata := map[string]interface{}{
		"cursor_user_id": cursorUser.ID,
	}
	metadataBytes, _ := json.Marshal(metadata)

	newAccount := &membertypes.ExternalAccount{
		AccountType:    "ai-code-assistant",
		ProviderName:   "cursor",
		ProviderID:     providerID,
		Username:       cursorUser.Email, // Use email if available, otherwise username
		OrganizationID: &organizationID,
		Metadata:       datatypes.JSON(metadataBytes),
	}

	if err := p.memberAPI.CreateExternalAccounts(ctx, []*membertypes.ExternalAccount{newAccount}); err != nil {
		return nil, fmt.Errorf("failed to create external account: %w", err)
	}

	// Fetch the created account to return it with ID
	accounts, err = p.memberAPI.GetExternalAccounts(ctx, &membertypes.ExternalAccountParams{
		OrganizationID: organizationID,
		AccountType:    &accountType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created account: %w", err)
	}

	for _, account := range accounts {
		if account.ProviderID == providerID && account.ProviderName == "cursor" {
			return &account, nil
		}
	}

	return nil, fmt.Errorf("failed to find created account")
}

// isSuggestionType checks if the usage type represents a suggestion
func isSuggestionType(usageType string) bool {
	return usageType == "edit" || usageType == "write" || usageType == "notebook_edit"
}

// intPtr returns a pointer to an int64
func intPtr(i int64) *int64 {
	return &i
}
