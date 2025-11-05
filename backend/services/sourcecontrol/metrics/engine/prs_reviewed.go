package engine

import (
	"context"
	"encoding/json"
	"time"

	"ems.dev/backend/libraries/errors"
	"ems.dev/backend/services/sourcecontrol/database"
	metrictypes "ems.dev/backend/services/sourcecontrol/metrics/types"
	"ems.dev/backend/services/sourcecontrol/types"
	"github.com/google/uuid"
)

type PRsReviewedRule struct {
	metrictypes.BaseMetricRule
	sourceControlDB database.DB
}

func NewPRsReviewedRule(baseMetricRule metrictypes.BaseMetricRule, sourceControlDB database.DB) *PRsReviewedRule {
	return &PRsReviewedRule{
		BaseMetricRule:  baseMetricRule,
		sourceControlDB: sourceControlDB,
	}
}

func (r *PRsReviewedRule) Calculate(ctx context.Context, params types.MetricRuleParams) (*types.SnapshotMetric, *types.GraphMetric, error) {
	// Validate params
	organizationID, startDate, endDate, sourceControlAccountIDs, peersSourceControlAccountIDs, prPrefixes, err := r.extractParams(params)
	if err != nil {
		return nil, nil, err
	}

	// Calculate PRs reviewed value
	prsReviewedValue, err := r.sourceControlDB.CalculatePRsReviewed(ctx, *organizationID, sourceControlAccountIDs, prPrefixes, *startDate, *endDate, r.Operation)
	if err != nil {
		return nil, nil, err
	}

	// Calculate peer values only if peer account IDs are provided (member metrics)
	var peersValue float64
	var timeSeries []types.TimeSeriesEntry

	// Calculate PRs reviewed graph value
	prsReviewedGraphValue, err := r.sourceControlDB.CalculatePRsReviewedGraph(ctx, *organizationID, sourceControlAccountIDs, prPrefixes, *startDate, *endDate, r.Operation, r.Name, params.Interval)
	if err != nil {
		return nil, nil, err
	}

	// Only calculate peer values if peer account IDs are provided (member metrics only)
	if len(peersSourceControlAccountIDs) > 0 {
		// Calculate PRs reviewed peers value
		peersPRsReviewedValue, err := r.sourceControlDB.CalculatePRsReviewedForAccounts(ctx, *organizationID, peersSourceControlAccountIDs, nil, *startDate, *endDate)
		if err != nil {
			return nil, nil, err
		}
		peersValue = float64(*peersPRsReviewedValue)

		// Calculate peer PRs reviewed graph value
		peersPRsReviewedGraphValue, err := r.sourceControlDB.CalculatePRsReviewedGraphForAccounts(ctx, *organizationID, peersSourceControlAccountIDs, nil, *startDate, *endDate, params.Interval)
		if err != nil {
			return nil, nil, err
		}

		// Merge peer values into the member's time series
		timeSeries = mergeTimeSeriesWithPeers(prsReviewedGraphValue, peersPRsReviewedGraphValue, r.Name)
	} else {
		// No peer account IDs, use member's time series as-is
		timeSeries = prsReviewedGraphValue
	}

	snapshotMetric := types.SnapshotMetric{
		Label:          r.Name,
		Description:    r.Description,
		Unit:           r.Unit,
		Value:          float64(*prsReviewedValue),
		PeersValue:     peersValue,
		IconIdentifier: r.IconIdentifier,
		IconColor:      r.IconColor,
	}

	graphMetric := types.GraphMetric{
		Label:      r.Name,
		Type:       "line",
		Unit:       r.Unit,
		TimeSeries: timeSeries,
	}

	return &snapshotMetric, &graphMetric, nil
}

func (r *PRsReviewedRule) Category() types.MetricRuleCategory {
	return r.BaseMetricRule.Category
}

func (r *PRsReviewedRule) extractParams(params types.MetricRuleParams) (*string, *time.Time, *time.Time, []string, []string, []string, error) {
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

	// Check if organization ID is present, this is needed also a security measure to prevent unauthorized access to other organizations
	orgID, exists := metricParams["organizationId"]
	if !exists {
		return nil, nil, nil, nil, nil, nil, errors.NewBadRequestError("organization id is required")
	}

	organizationID, ok := orgID.(string)
	if !ok {
		return nil, nil, nil, nil, nil, nil, errors.NewBadRequestError("invalid organization id format")
	}

	// If sourcecontrolaccountids is present check that these are valid uuids
	srcControlAccountIDs, exists := metricParams["sourceControlAccountIDs"]
	var sourceControlAccountIDs []string
	if exists {
		if srcControlAccountIDsArray, ok := srcControlAccountIDs.([]interface{}); ok {
			for _, idInterface := range srcControlAccountIDsArray {
				if id, ok := idInterface.(string); ok {
					if _, err := uuid.Parse(id); err != nil {
						return nil, nil, nil, nil, nil, nil, errors.NewBadRequestError("invalid source control account id")
					}
					sourceControlAccountIDs = append(sourceControlAccountIDs, id)
				}
			}
		}
	}

	// If peerssourcecontrolaccountids is present check that these are valid uuids
	peersSrcControlAccountIDs, exists := metricParams["peersSourceControlAccountIDs"]
	var peersSourceControlAccountIDs []string
	if exists {
		if peersSrcControlAccountIDsArray, ok := peersSrcControlAccountIDs.([]interface{}); ok {
			for _, idInterface := range peersSrcControlAccountIDsArray {
				if id, ok := idInterface.(string); ok {
					if _, err := uuid.Parse(id); err != nil {
						return nil, nil, nil, nil, nil, nil, errors.NewBadRequestError("invalid peers source control account id")
					}
					peersSourceControlAccountIDs = append(peersSourceControlAccountIDs, id)
				}
			}
		}
	}

	// Extract pr_prefixes if present
	prPrefixesInterface, exists := metricParams["pr_prefixes"]
	var prPrefixes []string
	if exists {
		if prPrefixesArray, ok := prPrefixesInterface.([]interface{}); ok {
			for _, prefixInterface := range prPrefixesArray {
				if prefix, ok := prefixInterface.(string); ok {
					prPrefixes = append(prPrefixes, prefix)
				}
			}
		}
	}

	return &organizationID, params.StartDate, params.EndDate, sourceControlAccountIDs, peersSourceControlAccountIDs, prPrefixes, nil
}
