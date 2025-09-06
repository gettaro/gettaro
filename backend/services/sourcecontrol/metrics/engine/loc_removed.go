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

type LOCRemovedRule struct {
	metrictypes.BaseMetricRule
	sourceControlDB database.DB
}

func NewLOCRemovedRule(baseMetricRule metrictypes.BaseMetricRule, sourceControlDB database.DB) *LOCRemovedRule {
	return &LOCRemovedRule{
		BaseMetricRule:  baseMetricRule,
		sourceControlDB: sourceControlDB,
	}
}

func (r *LOCRemovedRule) Calculate(ctx context.Context, params types.MetricRuleParams) (*types.SnapshotMetric, *types.GraphMetric, error) {
	// Validate params
	organizationID, startDate, endDate, sourceControlAccountIDs, peersSourceControlAccountIDs, err := r.extractParams(params)
	if err != nil {
		return nil, nil, err
	}

	// Calculate LOC removed value
	locRemovedValue, err := r.sourceControlDB.CalculateLOCRemoved(ctx, *organizationID, sourceControlAccountIDs, *startDate, *endDate, r.Operation)
	if err != nil {
		return nil, nil, err
	}

	// Calculate LOC removed peers value
	peersLOCRemovedValue, err := r.sourceControlDB.CalculateLOCRemovedForAccounts(ctx, *organizationID, peersSourceControlAccountIDs, *startDate, *endDate)
	if err != nil {
		return nil, nil, err
	}

	snapshotMetric := types.SnapshotMetric{
		Label:          r.Name,
		Description:    r.Description,
		Unit:           r.Unit,
		Value:          float64(*locRemovedValue),
		PeersValue:     float64(*peersLOCRemovedValue),
		IconIdentifier: r.IconIdentifier,
		IconColor:      r.IconColor,
	}

	// Calculate LOC removed graph value
	locRemovedGraphValue, err := r.sourceControlDB.CalculateLOCRemovedGraph(ctx, *organizationID, sourceControlAccountIDs, *startDate, *endDate, r.Operation, r.Name, params.Interval)
	if err != nil {
		return nil, nil, err
	}

	graphMetric := types.GraphMetric{
		Label:      r.Name,
		Unit:       r.Unit,
		TimeSeries: locRemovedGraphValue,
	}

	return &snapshotMetric, &graphMetric, nil
}

func (r *LOCRemovedRule) Category() string {
	return r.BaseMetricRule.Category
}

func (r *LOCRemovedRule) extractParams(params types.MetricRuleParams) (*string, *time.Time, *time.Time, []string, []string, error) {
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
						return nil, nil, nil, nil, nil, errors.NewBadRequestError("invalid peers source control account id")
					}
					peersSourceControlAccountIDs = append(peersSourceControlAccountIDs, id)
				}
			}
		}
	}

	return &organizationID, params.StartDate, params.EndDate, sourceControlAccountIDs, peersSourceControlAccountIDs, nil
}
