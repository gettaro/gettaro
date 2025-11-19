package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"ems.dev/backend/services/sourcecontrol/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCalculateMetrics(t *testing.T) {
	now := time.Now()
	startDate := now.AddDate(0, -1, 0)
	endDate := now

	tests := []struct {
		name           string
		params         types.MetricRuleParams
		mockResponse   *types.MetricsResponse
		mockError      error
		expectedResponse *types.MetricsResponse
		expectedError  error
	}{
		{
			name: "success - calculates metrics",
			params: types.MetricRuleParams{
				StartDate: &startDate,
				EndDate:   &endDate,
				Interval:  "daily",
			},
			mockResponse: &types.MetricsResponse{
				SnapshotMetrics: []*types.SnapshotCategory{
					{
						Category: types.MetricRuleCategory{
							Name:     "Activity",
							Priority: 1,
						},
						Metrics: []types.SnapshotMetric{
							{
								Label:       "PRs Merged",
								Description: "Total PRs merged",
								Value:       10,
								Unit:        types.UnitCount,
							},
						},
					},
				},
				GraphMetrics: []*types.GraphCategory{},
			},
			expectedResponse: &types.MetricsResponse{
				SnapshotMetrics: []*types.SnapshotCategory{
					{
						Category: types.MetricRuleCategory{
							Name:     "Activity",
							Priority: 1,
						},
						Metrics: []types.SnapshotMetric{
							{
								Label:       "PRs Merged",
								Description: "Total PRs merged",
								Value:       10,
								Unit:        types.UnitCount,
							},
						},
					},
				},
				GraphMetrics: []*types.GraphCategory{},
			},
		},
		{
			name: "success - with graph metrics",
			params: types.MetricRuleParams{
				StartDate: &startDate,
				EndDate:   &endDate,
				Interval:  "weekly",
			},
			mockResponse: &types.MetricsResponse{
				SnapshotMetrics: []*types.SnapshotCategory{},
				GraphMetrics: []*types.GraphCategory{
					{
						Category: types.MetricRuleCategory{
							Name:     "Activity",
							Priority: 1,
						},
						Metrics: []types.GraphMetric{
							{
								Label: "PRs Merged",
								Type:  "line",
								Unit:  types.UnitCount,
							},
						},
					},
				},
			},
			expectedResponse: &types.MetricsResponse{
				SnapshotMetrics: []*types.SnapshotCategory{},
				GraphMetrics: []*types.GraphCategory{
					{
						Category: types.MetricRuleCategory{
							Name:     "Activity",
							Priority: 1,
						},
						Metrics: []types.GraphMetric{
							{
								Label: "PRs Merged",
								Type:  "line",
								Unit:  types.UnitCount,
							},
						},
					},
				},
			},
		},
		{
			name: "error - metrics engine error",
			params: types.MetricRuleParams{
				StartDate: &startDate,
				EndDate:   &endDate,
			},
			mockError:     errors.New("metrics calculation failed"),
			expectedError: errors.New("metrics calculation failed"),
		},
		{
			name: "error - invalid date range",
			params: types.MetricRuleParams{
				StartDate: &endDate,
				EndDate:   &startDate, // End date before start date
			},
			mockError:     errors.New("invalid date range"),
			expectedError: errors.New("invalid date range"),
		},
		{
			name: "error - missing required params",
			params: types.MetricRuleParams{
				StartDate: nil,
				EndDate:   nil,
			},
			mockError:     errors.New("start date and end date are required"),
			expectedError: errors.New("start date and end date are required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			mockMetricsEngine := new(MockMetricsEngine)
			api := &Api{
				db:            mockDB,
				metricsEngine: mockMetricsEngine,
			}

			mockMetricsEngine.On("CalculateMetrics", mock.Anything, tt.params).Return(tt.mockResponse, tt.mockError)

			response, err := api.CalculateMetrics(context.Background(), tt.params)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, len(tt.expectedResponse.SnapshotMetrics), len(response.SnapshotMetrics))
				assert.Equal(t, len(tt.expectedResponse.GraphMetrics), len(response.GraphMetrics))
			}

			mockMetricsEngine.AssertExpectations(t)
		})
	}
}
