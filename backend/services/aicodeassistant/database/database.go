package database

import (
	"context"
	"fmt"
	"time"

	"ems.dev/backend/services/aicodeassistant/metrics/types"
	aicodeassistanttypes "ems.dev/backend/services/aicodeassistant/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AICodeAssistantDB defines the interface for AI code assistant database operations
type DB interface {
	CreateOrUpdateDailyMetric(ctx context.Context, metric *aicodeassistanttypes.AICodeAssistantDailyMetric) error
	GetDailyMetrics(ctx context.Context, params *aicodeassistanttypes.AICodeAssistantDailyMetricParams) ([]*aicodeassistanttypes.AICodeAssistantDailyMetric, error)

	// Calculate metrics
	CalculateLinesOfCodeAccepted(ctx context.Context, organizationID string, externalAccountIDs []string, toolNames []string, startDate, endDate time.Time, metricOperation types.MetricOperation) (*int, error)
	CalculateLinesOfCodeAcceptedGraph(ctx context.Context, organizationID string, externalAccountIDs []string, toolNames []string, startDate, endDate time.Time, metricOperation types.MetricOperation, metricLabel string, interval string) ([]aicodeassistanttypes.TimeSeriesEntry, error)
	CalculateLinesOfCodeSuggested(ctx context.Context, organizationID string, externalAccountIDs []string, toolNames []string, startDate, endDate time.Time, metricOperation types.MetricOperation) (*int, error)
	CalculateLinesOfCodeSuggestedGraph(ctx context.Context, organizationID string, externalAccountIDs []string, toolNames []string, startDate, endDate time.Time, metricOperation types.MetricOperation, metricLabel string, interval string) ([]aicodeassistanttypes.TimeSeriesEntry, error)
	CalculateActiveSessions(ctx context.Context, organizationID string, externalAccountIDs []string, toolNames []string, startDate, endDate time.Time, metricOperation types.MetricOperation) (*int, error)
	CalculateActiveSessionsGraph(ctx context.Context, organizationID string, externalAccountIDs []string, toolNames []string, startDate, endDate time.Time, metricOperation types.MetricOperation, metricLabel string, interval string) ([]aicodeassistanttypes.TimeSeriesEntry, error)
	CalculateAcceptRate(ctx context.Context, organizationID string, externalAccountIDs []string, toolNames []string, startDate, endDate time.Time, metricOperation types.MetricOperation) (*float64, error)
	CalculateAcceptRateGraph(ctx context.Context, organizationID string, externalAccountIDs []string, toolNames []string, startDate, endDate time.Time, metricOperation types.MetricOperation, metricLabel string, interval string) ([]aicodeassistanttypes.TimeSeriesEntry, error)

	// Calculate peer metrics (median across peers)
	CalculateLinesOfCodeAcceptedForAccounts(ctx context.Context, organizationID string, externalAccountIDs []string, toolNames []string, startDate, endDate time.Time) (*float64, error)
	CalculateLinesOfCodeSuggestedForAccounts(ctx context.Context, organizationID string, externalAccountIDs []string, toolNames []string, startDate, endDate time.Time) (*float64, error)
	CalculateActiveSessionsForAccounts(ctx context.Context, organizationID string, externalAccountIDs []string, toolNames []string, startDate, endDate time.Time) (*float64, error)
	CalculateAcceptRateForAccounts(ctx context.Context, organizationID string, externalAccountIDs []string, toolNames []string, startDate, endDate time.Time) (*float64, error)
}

// AICodeAssistantDB implements the database operations
type AICodeAssistantDB struct {
	db *gorm.DB
}

// NewAICodeAssistantDB creates a new instance of AICodeAssistantDB
func NewAICodeAssistantDB(db *gorm.DB) DB {
	return &AICodeAssistantDB{
		db: db,
	}
}

// CreateOrUpdateDailyMetric creates or updates a daily metric
// Uses ON CONFLICT to handle upsert logic based on unique constraint
func (d *AICodeAssistantDB) CreateOrUpdateDailyMetric(ctx context.Context, metric *aicodeassistanttypes.AICodeAssistantDailyMetric) error {
	// Use Clauses with OnConflict to handle upsert on unique constraint
	// The unique constraint is on (organization_id, external_account_id, tool_name, metric_date)
	return d.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "organization_id"},
				{Name: "external_account_id"},
				{Name: "tool_name"},
				{Name: "metric_date"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"lines_of_code_accepted",
				"lines_of_code_suggested",
				"suggestion_accept_rate",
				"active_sessions",
				"metadata",
				"updated_at",
			}),
		}).
		Create(metric).
		Error
}

// GetDailyMetrics retrieves daily metrics based on the given parameters
func (d *AICodeAssistantDB) GetDailyMetrics(ctx context.Context, params *aicodeassistanttypes.AICodeAssistantDailyMetricParams) ([]*aicodeassistanttypes.AICodeAssistantDailyMetric, error) {
	var metrics []*aicodeassistanttypes.AICodeAssistantDailyMetric
	query := d.db.WithContext(ctx).Model(&aicodeassistanttypes.AICodeAssistantDailyMetric{})

	// Filter by organization ID (required)
	query = query.Where("organization_id = ?", params.OrganizationID)

	// Filter by external account IDs if provided
	if len(params.ExternalAccountIDs) > 0 {
		query = query.Where("external_account_id IN ?", params.ExternalAccountIDs)
	}

	// Filter by tool name if provided
	if params.ToolName != nil && *params.ToolName != "" {
		query = query.Where("tool_name = ?", *params.ToolName)
	}

	// Filter by date range if provided
	if params.StartDate != nil {
		query = query.Where("metric_date >= ?", *params.StartDate)
	}
	if params.EndDate != nil {
		query = query.Where("metric_date <= ?", *params.EndDate)
	}

	// Order by date (oldest first)
	query = query.Order("metric_date ASC")

	if err := query.Find(&metrics).Error; err != nil {
		return nil, fmt.Errorf("failed to get daily metrics: %w", err)
	}

	return metrics, nil
}
