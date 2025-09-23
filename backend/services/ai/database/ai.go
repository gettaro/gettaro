package database

import (
	"context"
	"time"

	"ems.dev/backend/services/ai/types"
	"gorm.io/gorm"
)

type AIDB struct {
	db *gorm.DB
}

func NewAIDB(db *gorm.DB) *AIDB {
	return &AIDB{db: db}
}

// SaveQueryHistory saves a query to the history
func (d *AIDB) SaveQueryHistory(ctx context.Context, history *types.AIQueryHistory) error {
	return d.db.WithContext(ctx).Create(history).Error
}

// GetQueryHistory retrieves query history for a user or organization
func (d *AIDB) GetQueryHistory(ctx context.Context, organizationID string, userID *string, limit int) ([]*types.AIQueryHistory, error) {
	var history []*types.AIQueryHistory

	db := d.db.WithContext(ctx).Where("organization_id = ?", organizationID)
	if userID != nil {
		db = db.Where("user_id = ?", *userID)
	}

	err := db.Order("created_at DESC").Limit(limit).Find(&history).Error
	return history, err
}

// GetQueryStats returns statistics about AI usage
func (d *AIDB) GetQueryStats(ctx context.Context, organizationID string, userID *string, days int) (*types.AIQueryStats, error) {
	stats := &types.AIQueryStats{
		QueriesByEntity:  make(map[string]int),
		QueriesByContext: make(map[string]int),
	}

	// Calculate date range
	startDate := time.Now().AddDate(0, 0, -days)

	db := d.db.WithContext(ctx).Model(&types.AIQueryHistory{}).Where("organization_id = ? AND created_at >= ?", organizationID, startDate)
	if userID != nil {
		db = db.Where("user_id = ?", *userID)
	}

	// Total queries
	var totalCount int64
	err := db.Count(&totalCount).Error
	if err != nil {
		return nil, err
	}
	stats.TotalQueries = int(totalCount)

	// Queries by entity type
	var entityResults []struct {
		EntityType string
		Count      int
	}
	err = db.Select("entity_type, COUNT(*) as count").Group("entity_type").Find(&entityResults).Error
	if err != nil {
		return nil, err
	}

	for _, result := range entityResults {
		stats.QueriesByEntity[result.EntityType] = result.Count
	}

	// Queries by context
	var contextResults []struct {
		Context string
		Count   int
	}
	err = db.Select("context, COUNT(*) as count").Group("context").Find(&contextResults).Error
	if err != nil {
		return nil, err
	}

	for _, result := range contextResults {
		stats.QueriesByContext[result.Context] = result.Count
	}

	// Average confidence
	var avgConfidence float64
	err = db.Select("AVG(confidence)").Scan(&avgConfidence).Error
	if err != nil {
		return nil, err
	}
	stats.AverageConfidence = avgConfidence

	// Recent queries
	recentQueries, err := d.GetQueryHistory(ctx, organizationID, userID, 10)
	if err != nil {
		return nil, err
	}
	stats.RecentQueries = make([]types.AIQueryHistory, len(recentQueries))
	for i, query := range recentQueries {
		stats.RecentQueries[i] = *query
	}

	return stats, nil
}

// DeleteQueryHistory deletes old query history
func (d *AIDB) DeleteQueryHistory(ctx context.Context, organizationID string, olderThan time.Time) error {
	return d.db.WithContext(ctx).Where("organization_id = ? AND created_at < ?", organizationID, olderThan).Delete(&types.AIQueryHistory{}).Error
}

// GetEntityQueryCount returns the number of queries for a specific entity
func (d *AIDB) GetEntityQueryCount(ctx context.Context, organizationID string, entityType string, entityID string) (int, error) {
	var count int64
	err := d.db.WithContext(ctx).Model(&types.AIQueryHistory{}).Where("organization_id = ? AND entity_type = ? AND entity_id = ?", organizationID, entityType, entityID).Count(&count).Error
	return int(count), err
}
