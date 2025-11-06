package engine

import (
	"context"
	"encoding/json"
	"time"

	"ems.dev/backend/libraries/errors"
	"ems.dev/backend/services/aicodeassistant/database"
	metrictypes "ems.dev/backend/services/aicodeassistant/metrics/types"
	"ems.dev/backend/services/aicodeassistant/types"
	"github.com/google/uuid"
)

type LinesOfCodeSuggestedRule struct {
	metrictypes.BaseMetricRule
	aicodeassistantDB database.DB
}

func NewLinesOfCodeSuggestedRule(baseMetricRule metrictypes.BaseMetricRule, aicodeassistantDB database.DB) *LinesOfCodeSuggestedRule {
	return &LinesOfCodeSuggestedRule{
		BaseMetricRule:    baseMetricRule,
		aicodeassistantDB: aicodeassistantDB,
	}
}

func (r *LinesOfCodeSuggestedRule) Calculate(ctx context.Context, params types.MetricRuleParams) (*types.SnapshotMetric, *types.GraphMetric, error) {
	// Validate params
	organizationID, startDate, endDate, externalAccountIDs, peersExternalAccountIDs, toolNames, err := r.extractParams(params)
	if err != nil {
		return nil, nil, err
	}

	// Calculate lines of code suggested value
	locSuggestedValue, err := r.aicodeassistantDB.CalculateLinesOfCodeSuggested(ctx, *organizationID, externalAccountIDs, toolNames, *startDate, *endDate, r.Operation)
	if err != nil {
		return nil, nil, err
	}

	// Calculate peer values only if peer account IDs are provided (member metrics)
	var peersValue float64
	if len(peersExternalAccountIDs) > 0 {
		peersLOCSuggestedValueFloat, err := r.aicodeassistantDB.CalculateLinesOfCodeSuggestedForAccounts(ctx, *organizationID, peersExternalAccountIDs, toolNames, *startDate, *endDate)
		if err != nil {
			return nil, nil, err
		}
		peersValue = *peersLOCSuggestedValueFloat
	}

	snapshotMetric := types.SnapshotMetric{
		Label:          r.Name,
		Description:    r.Description,
		Unit:           r.Unit,
		Value:          float64(*locSuggestedValue),
		PeersValue:     peersValue,
		IconIdentifier: r.IconIdentifier,
		IconColor:      r.IconColor,
	}

	// Calculate lines of code suggested graph value
	locSuggestedGraphValue, err := r.aicodeassistantDB.CalculateLinesOfCodeSuggestedGraph(ctx, *organizationID, externalAccountIDs, toolNames, *startDate, *endDate, r.Operation, r.Name, params.Interval)
	if err != nil {
		return nil, nil, err
	}

	// Calculate peers graph value and merge with main metric only if peer account IDs are provided
	var mergedTimeSeries []types.TimeSeriesEntry
	if len(peersExternalAccountIDs) > 0 {
		peersGraphValue, err := r.aicodeassistantDB.CalculateLinesOfCodeSuggestedGraphForAccounts(ctx, *organizationID, peersExternalAccountIDs, toolNames, *startDate, *endDate, params.Interval)
		if err != nil {
			return nil, nil, err
		}
		// Merge peers data into time series by date
		mergedTimeSeries = r.mergeTimeSeriesData(locSuggestedGraphValue, peersGraphValue)
	} else {
		mergedTimeSeries = locSuggestedGraphValue
	}

	graphMetric := types.GraphMetric{
		Label:      r.Name,
		Type:       "line",
		Unit:       r.Unit,
		TimeSeries: mergedTimeSeries,
	}

	return &snapshotMetric, &graphMetric, nil
}

// mergeTimeSeriesData merges main metric and peers data by date
func (r *LinesOfCodeSuggestedRule) mergeTimeSeriesData(mainSeries []types.TimeSeriesEntry, peersSeries []types.TimeSeriesEntry) []types.TimeSeriesEntry {
	// Create a map of dates to entries for quick lookup
	dateMap := make(map[string]*types.TimeSeriesEntry)

	// First, add all main series entries
	for i := range mainSeries {
		dateMap[mainSeries[i].Date] = &mainSeries[i]
	}

	// Then, merge peers data into existing entries or create new ones
	for _, peerEntry := range peersSeries {
		if entry, exists := dateMap[peerEntry.Date]; exists {
			// Merge peers data point into existing entry
			if len(peerEntry.Data) > 0 {
				entry.Data = append(entry.Data, peerEntry.Data[0])
			}
		} else {
			// Create new entry with peers data
			newEntry := types.TimeSeriesEntry{
				Date: peerEntry.Date,
				Data: peerEntry.Data,
			}
			dateMap[peerEntry.Date] = &newEntry
		}
	}

	// Convert map back to slice and sort by date
	result := make([]types.TimeSeriesEntry, 0, len(dateMap))
	for _, entry := range dateMap {
		result = append(result, *entry)
	}

	// Sort by date
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Date > result[j].Date {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result
}

func (r *LinesOfCodeSuggestedRule) Category() types.MetricRuleCategory {
	return r.BaseMetricRule.Category
}

func (r *LinesOfCodeSuggestedRule) extractParams(params types.MetricRuleParams) (*string, *time.Time, *time.Time, []string, []string, []string, error) {
	if params.Interval == "" {
		return nil, nil, nil, nil, nil, nil, errors.NewBadRequestError("interval is required")
	}

	if params.Interval != "daily" && params.Interval != "weekly" && params.Interval != "monthly" {
		return nil, nil, nil, nil, nil, nil, errors.NewBadRequestError("invalid interval")
	}

	if params.StartDate == nil {
		return nil, nil, nil, nil, nil, nil, errors.NewBadRequestError("start date is required")
	}

	if params.EndDate == nil {
		return nil, nil, nil, nil, nil, nil, errors.NewBadRequestError("end date is required")
	}

	if params.MetricParams == nil {
		return nil, nil, nil, nil, nil, nil, errors.NewBadRequestError("metric params is required")
	}

	// Unmarshal MetricParams to check for organization ID
	var metricParams map[string]interface{}
	if err := json.Unmarshal(params.MetricParams, &metricParams); err != nil {
		return nil, nil, nil, nil, nil, nil, errors.NewBadRequestError("invalid metric params format")
	}

	// Check if organization ID is present
	orgID, exists := metricParams["organizationId"]
	if !exists {
		return nil, nil, nil, nil, nil, nil, errors.NewBadRequestError("organization id is required")
	}

	organizationID, ok := orgID.(string)
	if !ok {
		return nil, nil, nil, nil, nil, nil, errors.NewBadRequestError("invalid organization id format")
	}

	// If externalAccountIDs is present check that these are valid uuids
	externalAccountIDsInterface, exists := metricParams["externalAccountIDs"]
	var externalAccountIDs []string
	if exists {
		if externalAccountIDsArray, ok := externalAccountIDsInterface.([]interface{}); ok {
			for _, idInterface := range externalAccountIDsArray {
				if id, ok := idInterface.(string); ok {
					if _, err := uuid.Parse(id); err != nil {
						return nil, nil, nil, nil, nil, nil, errors.NewBadRequestError("invalid external account id")
					}
					externalAccountIDs = append(externalAccountIDs, id)
				}
			}
		}
	}

	// If peersExternalAccountIDs is present check that these are valid uuids
	peersExternalAccountIDsInterface, exists := metricParams["peersExternalAccountIDs"]
	var peersExternalAccountIDs []string
	if exists {
		if peersExternalAccountIDsArray, ok := peersExternalAccountIDsInterface.([]interface{}); ok {
			for _, idInterface := range peersExternalAccountIDsArray {
				if id, ok := idInterface.(string); ok {
					if _, err := uuid.Parse(id); err != nil {
						return nil, nil, nil, nil, nil, nil, errors.NewBadRequestError("invalid peers external account id")
					}
					peersExternalAccountIDs = append(peersExternalAccountIDs, id)
				}
			}
		}
	}

	// Extract tool names if present
	toolNamesInterface, exists := metricParams["toolNames"]
	var toolNames []string
	if exists {
		if toolNamesArray, ok := toolNamesInterface.([]interface{}); ok {
			for _, toolNameInterface := range toolNamesArray {
				if toolName, ok := toolNameInterface.(string); ok {
					toolNames = append(toolNames, toolName)
				}
			}
		}
	}

	return &organizationID, params.StartDate, params.EndDate, externalAccountIDs, peersExternalAccountIDs, toolNames, nil
}
