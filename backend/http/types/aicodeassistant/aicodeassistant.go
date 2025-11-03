package aicodeassistant

import (
	"encoding/json"
	"time"

	"ems.dev/backend/services/aicodeassistant/types"
)

type AICodeAssistantDailyMetricResponse struct {
	ID                   string      `json:"id"`
	OrganizationID       string      `json:"organization_id"`
	ExternalAccountID    string      `json:"external_account_id"`
	ToolName             string      `json:"tool_name"`
	MetricDate           time.Time   `json:"metric_date"`
	LinesOfCodeAccepted  int         `json:"lines_of_code_accepted"`
	LinesOfCodeSuggested int         `json:"lines_of_code_suggested"`
	SuggestionAcceptRate *float64    `json:"suggestion_accept_rate,omitempty"`
	ActiveSessions       int         `json:"active_sessions"`
	Metadata             interface{} `json:"metadata,omitempty"`
	CreatedAt            time.Time   `json:"created_at"`
	UpdatedAt            time.Time   `json:"updated_at"`
}

// MarshalAICodeAssistantDailyMetric converts a service type to HTTP response type
func MarshalAICodeAssistantDailyMetric(m *types.AICodeAssistantDailyMetric) *AICodeAssistantDailyMetricResponse {
	var metadata interface{}
	if len(m.Metadata) > 0 {
		// Metadata is already JSONB, unmarshal to interface{} for JSON marshaling
		var metadataMap map[string]interface{}
		if err := json.Unmarshal(m.Metadata, &metadataMap); err == nil {
			metadata = metadataMap
		} else {
			// If unmarshaling fails, return nil
			metadata = nil
		}
	}

	return &AICodeAssistantDailyMetricResponse{
		ID:                   m.ID,
		OrganizationID:       m.OrganizationID,
		ExternalAccountID:    m.ExternalAccountID,
		ToolName:             m.ToolName,
		MetricDate:           m.MetricDate,
		LinesOfCodeAccepted:  m.LinesOfCodeAccepted,
		LinesOfCodeSuggested: m.LinesOfCodeSuggested,
		SuggestionAcceptRate: m.SuggestionAcceptRate,
		ActiveSessions:       m.ActiveSessions,
		Metadata:             metadata,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}
}

