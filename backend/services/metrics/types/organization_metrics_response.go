package types

import (
	sourcecontroltypes "ems.dev/backend/services/sourcecontrol/types"
)

// TeamMetricsBreakdown represents metrics for a specific team
type TeamMetricsBreakdown struct {
	TeamID          string                           `json:"team_id"`
	TeamName        string                           `json:"team_name"`
	SnapshotMetrics []*sourcecontroltypes.SnapshotCategory `json:"snapshot_metrics"`
	GraphMetrics    []*sourcecontroltypes.GraphCategory    `json:"graph_metrics"`
}

// OrganizationMetricsResponse represents the response for organization metrics
// It includes both cumulative metrics and a breakdown by team
type OrganizationMetricsResponse struct {
	// Cumulative metrics across all selected teams/members
	SnapshotMetrics []*sourcecontroltypes.SnapshotCategory `json:"snapshot_metrics"`
	GraphMetrics    []*sourcecontroltypes.GraphCategory    `json:"graph_metrics"`
	
	// Breakdown by team (only included if teams are specified)
	TeamsBreakdown  []TeamMetricsBreakdown `json:"teams_breakdown,omitempty"`
}

