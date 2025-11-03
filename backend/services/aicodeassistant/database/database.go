package database

import (
	"context"
	"fmt"

	"ems.dev/backend/services/aicodeassistant/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AICodeAssistantDB defines the interface for AI code assistant database operations
type DB interface {
	CreateOrUpdateDailyMetric(ctx context.Context, metric *types.AICodeAssistantDailyMetric) error
	GetDailyMetrics(ctx context.Context, params *types.AICodeAssistantDailyMetricParams) ([]*types.AICodeAssistantDailyMetric, error)
	GetUsageStats(ctx context.Context, params *types.AICodeAssistantDailyMetricParams) (*types.AICodeAssistantUsageStats, error)
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
func (d *AICodeAssistantDB) CreateOrUpdateDailyMetric(ctx context.Context, metric *types.AICodeAssistantDailyMetric) error {
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
				"total_suggestions",
				"suggestions_accepted",
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
func (d *AICodeAssistantDB) GetDailyMetrics(ctx context.Context, params *types.AICodeAssistantDailyMetricParams) ([]*types.AICodeAssistantDailyMetric, error) {
	var metrics []*types.AICodeAssistantDailyMetric
	query := d.db.WithContext(ctx).Model(&types.AICodeAssistantDailyMetric{})

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

// GetUsageStats calculates aggregated statistics from daily metrics
func (d *AICodeAssistantDB) GetUsageStats(ctx context.Context, params *types.AICodeAssistantDailyMetricParams) (*types.AICodeAssistantUsageStats, error) {
	query := d.db.WithContext(ctx).Model(&types.AICodeAssistantDailyMetric{})

	// Apply same filters as GetDailyMetrics
	query = query.Where("organization_id = ?", params.OrganizationID)

	if len(params.ExternalAccountIDs) > 0 {
		query = query.Where("external_account_id IN ?", params.ExternalAccountIDs)
	}

	if params.ToolName != nil && *params.ToolName != "" {
		query = query.Where("tool_name = ?", *params.ToolName)
	}

	if params.StartDate != nil {
		query = query.Where("metric_date >= ?", *params.StartDate)
	}
	if params.EndDate != nil {
		query = query.Where("metric_date <= ?", *params.EndDate)
	}

	// Initialize stats with zero values
	stats := types.AICodeAssistantUsageStats{
		TotalLinesAccepted: 0,
		TotalSuggestions:   0,
		OverallAcceptRate:  0.0,
		ActiveUsers:        0,
	}

	// Calculate total lines accepted
	var totalLinesAccepted struct {
		Sum int64
	}
	if err := query.Select("COALESCE(SUM(lines_of_code_accepted), 0) as sum").
		Scan(&totalLinesAccepted).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate total lines accepted: %w", err)
	}
	stats.TotalLinesAccepted = int(totalLinesAccepted.Sum)

	// Calculate total suggestions
	var totalSuggestions struct {
		Sum int64
	}
	if err := query.Select("COALESCE(SUM(total_suggestions), 0) as sum").
		Scan(&totalSuggestions).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate total suggestions: %w", err)
	}
	stats.TotalSuggestions = int(totalSuggestions.Sum)

	// Calculate total accepted suggestions
	var totalAccepted struct {
		Sum int64
	}
	if err := query.Select("COALESCE(SUM(suggestions_accepted), 0) as sum").
		Scan(&totalAccepted).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate total accepted: %w", err)
	}

	// Calculate overall accept rate
	if stats.TotalSuggestions > 0 {
		stats.OverallAcceptRate = (float64(totalAccepted.Sum) / float64(stats.TotalSuggestions)) * 100
	}

	// Count unique active users (distinct external_account_ids)
	var activeUsers struct {
		Count int64
	}
	if err := query.Select("COALESCE(COUNT(DISTINCT external_account_id), 0) as count").
		Scan(&activeUsers).Error; err != nil {
		return nil, fmt.Errorf("failed to count active users: %w", err)
	}
	stats.ActiveUsers = int(activeUsers.Count)

	return &stats, nil
}
