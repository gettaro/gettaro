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

	// Calculate lines of code suggested peers value
	peersLOCSuggestedValue, err := r.aicodeassistantDB.CalculateLinesOfCodeSuggestedForAccounts(ctx, *organizationID, peersExternalAccountIDs, toolNames, *startDate, *endDate)
	if err != nil {
		return nil, nil, err
	}

	snapshotMetric := types.SnapshotMetric{
		Label:          r.Name,
		Description:    r.Description,
		Unit:           r.Unit,
		Value:          float64(*locSuggestedValue),
		PeersValue:     float64(*peersLOCSuggestedValue),
		IconIdentifier: r.IconIdentifier,
		IconColor:      r.IconColor,
	}

	// Calculate lines of code suggested graph value
	locSuggestedGraphValue, err := r.aicodeassistantDB.CalculateLinesOfCodeSuggestedGraph(ctx, *organizationID, externalAccountIDs, toolNames, *startDate, *endDate, r.Operation, r.Name, params.Interval)
	if err != nil {
		return nil, nil, err
	}

	graphMetric := types.GraphMetric{
		Label:      r.Name,
		Type:       "line",
		Unit:       r.Unit,
		TimeSeries: locSuggestedGraphValue,
	}

	return &snapshotMetric, &graphMetric, nil
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
