package types

import (
	"time"
)

// OrganizationMetricsParams represents the parameters for getting organization metrics
type OrganizationMetricsParams struct {
	OrganizationID string     `json:"organization_id"`
	TeamIDs        []string   `json:"team_ids,omitempty"` // Optional: filter by team IDs
	StartDate      *time.Time `json:"start_date,omitempty"`
	EndDate        *time.Time `json:"end_date,omitempty"`
	Interval       string     `json:"interval,omitempty"` // daily, weekly, monthly
}

