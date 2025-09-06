package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"ems.dev/backend/libraries/errors"
	"ems.dev/backend/services/sourcecontrol/database"
	metrictypes "ems.dev/backend/services/sourcecontrol/metrics/types"
	"ems.dev/backend/services/sourcecontrol/types"
	"github.com/google/uuid"
)

// PRReviewComplexityRule calculates the average complexity of PRs reviewed by a member
type PRReviewComplexityRule struct {
	metrictypes.BaseMetricRule
	sourceControlDB database.DB
}

// NewPRReviewComplexityRule creates a new PR Review Complexity rule
func NewPRReviewComplexityRule(base metrictypes.BaseMetricRule, sourceControlDB database.DB) *PRReviewComplexityRule {
	return &PRReviewComplexityRule{
		BaseMetricRule:  base,
		sourceControlDB: sourceControlDB,
	}
}

// Calculate implements the MetricRule interface
func (r *PRReviewComplexityRule) Calculate(ctx context.Context, params types.MetricRuleParams) (*types.SnapshotMetric, *types.GraphMetric, error) {
	// Extract parameters
	organizationID, startDate, endDate, sourceControlAccountIDs, peersSourceControlAccountIDs, err := r.extractParams(params)
	if err != nil {
		return nil, nil, err
	}

	// Calculate PR review complexity for the member
	prReviewComplexityValue, err := r.sourceControlDB.CalculatePRReviewComplexity(
		ctx,
		*organizationID,
		sourceControlAccountIDs,
		*startDate,
		*endDate,
		r.Operation,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to calculate PR review complexity: %w", err)
	}

	// Calculate PR review complexity for peers (other members in the organization)
	peersPRReviewComplexityValue, err := r.sourceControlDB.CalculatePRReviewComplexityForAccounts(
		ctx,
		*organizationID,
		peersSourceControlAccountIDs,
		*startDate,
		*endDate,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to calculate peers PR review complexity: %w", err)
	}

	// Create snapshot metric
	snapshotMetric := types.SnapshotMetric{
		Label:          r.Name,
		Description:    r.Description,
		Unit:           r.Unit,
		Value:          *prReviewComplexityValue,
		PeersValue:     *peersPRReviewComplexityValue,
		IconIdentifier: r.IconIdentifier,
		IconColor:      r.IconColor,
	}

	// Calculate graph metric
	graphMetric, err := r.calculateGraphMetric(ctx, *organizationID, sourceControlAccountIDs, *startDate, *endDate, params.Interval)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to calculate graph metric: %w", err)
	}

	return &snapshotMetric, graphMetric, nil
}

// calculateGraphMetric calculates the time series data for the metric
func (r *PRReviewComplexityRule) calculateGraphMetric(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, interval string) (*types.GraphMetric, error) {
	// Calculate time series data
	timeSeriesData, err := r.sourceControlDB.CalculatePRReviewComplexityGraph(
		ctx,
		organizationID,
		sourceControlAccountIDs,
		startDate,
		endDate,
		r.Operation,
		r.Name,
		interval,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate PR review complexity graph: %w", err)
	}

	// Create graph metric
	graphMetric := types.GraphMetric{
		Label:      r.Name,
		TimeSeries: timeSeriesData,
	}

	return &graphMetric, nil
}

// extractParams extracts and validates the metric parameters
func (r *PRReviewComplexityRule) extractParams(params types.MetricRuleParams) (*string, *time.Time, *time.Time, []string, []string, error) {
	if params.Interval == "" {
		return nil, nil, nil, nil, nil, errors.NewBadRequestError("interval is required")
	}

	if params.Interval != "daily" && params.Interval != "weekly" && params.Interval != "monthly" {
		return nil, nil, nil, nil, nil, errors.NewBadRequestError("invalid interval")
	}

	if params.StartDate == nil {
		return nil, nil, nil, nil, nil, errors.NewBadRequestError("start date is required")
	}

	if params.EndDate == nil {
		return nil, nil, nil, nil, nil, errors.NewBadRequestError("end date is required")
	}

	if params.MetricParams == nil {
		return nil, nil, nil, nil, nil, errors.NewBadRequestError("metric params is required")
	}

	// Unmarshal MetricParams to check for organization ID
	var metricParams map[string]interface{}
	if err := json.Unmarshal(params.MetricParams, &metricParams); err != nil {
		return nil, nil, nil, nil, nil, errors.NewBadRequestError("invalid metric params format")
	}

	// Check if organization ID is present, this is needed also a security measure to prevent unauthorized access to other organizations
	orgID, exists := metricParams["organizationId"]
	if !exists {
		return nil, nil, nil, nil, nil, errors.NewBadRequestError("organization id is required")
	}

	organizationID, ok := orgID.(string)
	if !ok {
		return nil, nil, nil, nil, nil, errors.NewBadRequestError("invalid organization id format")
	}

	// If sourcecontrolaccountids is present check that these are valid uuids
	srcControlAccountIDs, exists := metricParams["sourceControlAccountIDs"]
	var sourceControlAccountIDs []string
	if exists {
		if srcControlAccountIDsArray, ok := srcControlAccountIDs.([]interface{}); ok {
			for _, idInterface := range srcControlAccountIDsArray {
				if id, ok := idInterface.(string); ok {
					if _, err := uuid.Parse(id); err != nil {
						return nil, nil, nil, nil, nil, errors.NewBadRequestError("invalid source control account id")
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
						return nil, nil, nil, nil, nil, errors.NewBadRequestError("invalid peer source control account id")
					}
					peersSourceControlAccountIDs = append(peersSourceControlAccountIDs, id)
				}
			}
		}
	}

	return &organizationID, params.StartDate, params.EndDate, sourceControlAccountIDs, peersSourceControlAccountIDs, nil
}

// Category returns the category of the metric
func (r *PRReviewComplexityRule) Category() types.MetricRuleCategory {
	return r.BaseMetricRule.Category
}
