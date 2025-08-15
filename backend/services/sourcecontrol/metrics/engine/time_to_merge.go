package engine

import (
	"context"
	"encoding/json"

	"ems.dev/backend/libraries/errors"
	"ems.dev/backend/services/sourcecontrol/database"
	metrictypes "ems.dev/backend/services/sourcecontrol/metrics/types"
	"ems.dev/backend/services/sourcecontrol/types"
	"github.com/google/uuid"
)

type TimeToMergeRule struct {
	metrictypes.BaseMetricRule
	sourceControlDB database.DB
}

func NewTimeToMergeRule(sourceControlDB database.DB) *TimeToMergeRule {
	return &TimeToMergeRule{
		sourceControlDB: sourceControlDB,
	}
}

func (r *TimeToMergeRule) Calculate(ctx context.Context, params types.MetricRuleParams) (*types.SnapshotMetric, *types.GraphMetric, error) {
	// Validate params
	if err := r.ValidateParams(ctx, params); err != nil {
		return nil, nil, err
	}

	// Unmarshal MetricParams to check for organization ID
	var metricParams map[string]interface{}
	if err := json.Unmarshal(params.MetricParams, &metricParams); err != nil {
		return nil, nil, errors.NewBadRequestError("invalid metric params format")
	}

	organizationID := metricParams["organizationId"].(string)
	sourceControlAccountIDs := metricParams["sourceControlAccountIDs"].([]string)
	peersSourceControlAccountIDs := metricParams["peersSourceControlAccountIDs"].([]string)

	// Calculate time to merge value
	timeToMergeValue, err := r.sourceControlDB.CalculateTimeToMerge(ctx, organizationID, sourceControlAccountIDs, *params.StartDate, *params.EndDate, r.Operation)
	if err != nil {
		return nil, nil, err
	}

	// Calculate time to merge peers value
	peersTimeToMergeValue, err := r.sourceControlDB.CalculateTimeToMerge(ctx, organizationID, peersSourceControlAccountIDs, *params.StartDate, *params.EndDate, r.Operation)
	if err != nil {
		return nil, nil, err
	}

	snapshotMetric := types.SnapshotMetric{
		Label:      r.Name,
		Unit:       r.Unit,
		Value:      float64(*timeToMergeValue),
		PeersValue: float64(*peersTimeToMergeValue),
	}

	// Calculate time to merge graph value
	timeToMergeGraphValue, err := r.sourceControlDB.CalculateTimeToMergeGraph(ctx, organizationID, sourceControlAccountIDs, *params.StartDate, *params.EndDate, r.Operation, r.Name, params.Interval)
	if err != nil {
		return nil, nil, err
	}

	graphMetric := types.GraphMetric{
		Label:      r.Name,
		Unit:       r.Unit,
		TimeSeries: timeToMergeGraphValue,
	}

	return &snapshotMetric, &graphMetric, nil
}

func (r *TimeToMergeRule) ValidateParams(ctx context.Context, params types.MetricRuleParams) error {
	if params.Interval == "" {
		return errors.NewBadRequestError("interval is required")
	}

	if params.Interval != "daily" && params.Interval != "weekly" && params.Interval != "monthly" {
		return errors.NewBadRequestError("invalid interval")
	}

	if params.StartDate == nil {
		return errors.NewBadRequestError("start date is required")
	}

	if params.EndDate == nil {
		return errors.NewBadRequestError("end date is required")
	}

	if params.MetricParams == nil {
		return errors.NewBadRequestError("metric params is required")
	}

	// Unmarshal MetricParams to check for organization ID
	var metricParams map[string]interface{}
	if err := json.Unmarshal(params.MetricParams, &metricParams); err != nil {
		return errors.NewBadRequestError("invalid metric params format")
	}

	// Check if organization ID is present, this is needed also a security measure to prevent unauthorized access to other organizations
	if orgID, exists := metricParams["organizationId"]; !exists || orgID == nil {
		return errors.NewBadRequestError("organization id is required")
	}

	// If sourcecontrolaccountids is present check that these are valid uuids
	if sourceControlAccountIDs, exists := metricParams["sourceControlAccountIDs"]; exists {
		if sourceControlAccountIDs, ok := sourceControlAccountIDs.([]string); ok {
			for _, sourceControlAccountID := range sourceControlAccountIDs {
				if _, err := uuid.Parse(sourceControlAccountID); err != nil {
					return errors.NewBadRequestError("invalid source control account id")
				}
			}
		}
	}

	// If peerssourcecontrolaccountids is present check that these are valid uuids
	if peersSourceControlAccountIDs, exists := metricParams["peersSourceControlAccountIDs"]; exists {
		if sourceControlAccountIDs, ok := peersSourceControlAccountIDs.([]string); ok {
			for _, sourceControlAccountID := range sourceControlAccountIDs {
				if _, err := uuid.Parse(sourceControlAccountID); err != nil {
					return errors.NewBadRequestError("invalid peers source control account id")
				}
			}
		}
	}

	return nil
}
