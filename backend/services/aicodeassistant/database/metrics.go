package database

import (
	"context"
	"fmt"
	"time"

	"ems.dev/backend/services/aicodeassistant/metrics/types"
	aicodeassistanttypes "ems.dev/backend/services/aicodeassistant/types"
)

// CalculateLinesOfCodeAccepted calculates the lines of code accepted metric
func (d *AICodeAssistantDB) CalculateLinesOfCodeAccepted(ctx context.Context, organizationID string, externalAccountIDs []string, toolNames []string, startDate, endDate time.Time, metricOperation types.MetricOperation) (*int, error) {
	selectStatement := ""
	switch metricOperation {
	case types.MetricOperationCount:
		selectStatement = "COALESCE(SUM(lines_of_code_accepted), 0)"
	default:
		return nil, fmt.Errorf("invalid metric operation for lines of code accepted: %s", metricOperation)
	}

	query := `
		SELECT ` + selectStatement + ` as loc_accepted_count
		FROM ai_code_assistant_daily_metrics
		WHERE organization_id = ?
		AND metric_date >= ?
		AND metric_date <= ?
	`

	var args []any
	args = append(args, organizationID, startDate, endDate)

	if len(externalAccountIDs) > 0 {
		query += " AND external_account_id IN ?"
		args = append(args, externalAccountIDs)
	}

	if len(toolNames) > 0 {
		query += " AND tool_name IN ?"
		args = append(args, toolNames)
	}

	var count int64
	if err := d.db.WithContext(ctx).Raw(query, args...).Scan(&count).Error; err != nil {
		return nil, err
	}

	value := int(count)
	return &value, nil
}

// CalculateLinesOfCodeAcceptedGraph calculates the lines of code accepted metric for a graph
func (d *AICodeAssistantDB) CalculateLinesOfCodeAcceptedGraph(ctx context.Context, organizationID string, externalAccountIDs []string, toolNames []string, startDate, endDate time.Time, metricOperation types.MetricOperation, metricLabel string, interval string) ([]aicodeassistanttypes.TimeSeriesEntry, error) {
	selectStatement := ""
	switch metricOperation {
	case types.MetricOperationCount:
		selectStatement = "COALESCE(SUM(lines_of_code_accepted), 0)"
	default:
		return nil, fmt.Errorf("invalid metric operation for lines of code accepted: %s", metricOperation)
	}

	// Map interval values to PostgreSQL DATE_TRUNC units
	postgresInterval := interval
	switch interval {
	case "daily":
		postgresInterval = "day"
	case "weekly":
		postgresInterval = "week"
	case "monthly":
		postgresInterval = "month"
	}

	query := `
		SELECT 
			DATE_TRUNC('` + postgresInterval + `', metric_date) as date,
			` + selectStatement + ` as loc_accepted_count
		FROM ai_code_assistant_daily_metrics
		WHERE organization_id = ?
		AND metric_date >= ?
		AND metric_date <= ?
	`

	var args []any
	args = append(args, organizationID, startDate, endDate)

	if len(externalAccountIDs) > 0 {
		query += " AND external_account_id IN ?"
		args = append(args, externalAccountIDs)
	}

	if len(toolNames) > 0 {
		query += " AND tool_name IN ?"
		args = append(args, toolNames)
	}

	query += " GROUP BY DATE_TRUNC('" + postgresInterval + "', metric_date)"
	query += " ORDER BY date"

	var result struct {
		Date            time.Time `json:"date"`
		LOCAcceptedCount float64   `json:"loc_accepted_count"`
	}

	rows, err := d.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dataPoints := []aicodeassistanttypes.TimeSeriesEntry{}
	for rows.Next() {
		if err := rows.Scan(&result.Date, &result.LOCAcceptedCount); err != nil {
			return nil, err
		}
		dataPoints = append(dataPoints, aicodeassistanttypes.TimeSeriesEntry{
			Date: result.Date.Format("2006-01-02"),
			Data: []aicodeassistanttypes.TimeSeriesDataPoint{
				{
					Key:   metricLabel,
					Value: result.LOCAcceptedCount,
				},
			},
		})
	}

	return dataPoints, nil
}

// CalculateLinesOfCodeSuggested calculates the lines of code suggested metric
func (d *AICodeAssistantDB) CalculateLinesOfCodeSuggested(ctx context.Context, organizationID string, externalAccountIDs []string, toolNames []string, startDate, endDate time.Time, metricOperation types.MetricOperation) (*int, error) {
	selectStatement := ""
	switch metricOperation {
	case types.MetricOperationCount:
		selectStatement = "COALESCE(SUM(lines_of_code_suggested), 0)"
	default:
		return nil, fmt.Errorf("invalid metric operation for lines of code suggested: %s", metricOperation)
	}

	query := `
		SELECT ` + selectStatement + ` as loc_suggested_count
		FROM ai_code_assistant_daily_metrics
		WHERE organization_id = ?
		AND metric_date >= ?
		AND metric_date <= ?
	`

	var args []any
	args = append(args, organizationID, startDate, endDate)

	if len(externalAccountIDs) > 0 {
		query += " AND external_account_id IN ?"
		args = append(args, externalAccountIDs)
	}

	if len(toolNames) > 0 {
		query += " AND tool_name IN ?"
		args = append(args, toolNames)
	}

	var count int64
	if err := d.db.WithContext(ctx).Raw(query, args...).Scan(&count).Error; err != nil {
		return nil, err
	}

	value := int(count)
	return &value, nil
}

// CalculateLinesOfCodeSuggestedGraph calculates the lines of code suggested metric for a graph
func (d *AICodeAssistantDB) CalculateLinesOfCodeSuggestedGraph(ctx context.Context, organizationID string, externalAccountIDs []string, toolNames []string, startDate, endDate time.Time, metricOperation types.MetricOperation, metricLabel string, interval string) ([]aicodeassistanttypes.TimeSeriesEntry, error) {
	selectStatement := ""
	switch metricOperation {
	case types.MetricOperationCount:
		selectStatement = "COALESCE(SUM(lines_of_code_suggested), 0)"
	default:
		return nil, fmt.Errorf("invalid metric operation for lines of code suggested: %s", metricOperation)
	}

	// Map interval values to PostgreSQL DATE_TRUNC units
	postgresInterval := interval
	switch interval {
	case "daily":
		postgresInterval = "day"
	case "weekly":
		postgresInterval = "week"
	case "monthly":
		postgresInterval = "month"
	}

	query := `
		SELECT 
			DATE_TRUNC('` + postgresInterval + `', metric_date) as date,
			` + selectStatement + ` as loc_suggested_count
		FROM ai_code_assistant_daily_metrics
		WHERE organization_id = ?
		AND metric_date >= ?
		AND metric_date <= ?
	`

	var args []any
	args = append(args, organizationID, startDate, endDate)

	if len(externalAccountIDs) > 0 {
		query += " AND external_account_id IN ?"
		args = append(args, externalAccountIDs)
	}

	if len(toolNames) > 0 {
		query += " AND tool_name IN ?"
		args = append(args, toolNames)
	}

	query += " GROUP BY DATE_TRUNC('" + postgresInterval + "', metric_date)"
	query += " ORDER BY date"

	var result struct {
		Date             time.Time `json:"date"`
		LOCSuggestedCount float64   `json:"loc_suggested_count"`
	}

	rows, err := d.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dataPoints := []aicodeassistanttypes.TimeSeriesEntry{}
	for rows.Next() {
		if err := rows.Scan(&result.Date, &result.LOCSuggestedCount); err != nil {
			return nil, err
		}
		dataPoints = append(dataPoints, aicodeassistanttypes.TimeSeriesEntry{
			Date: result.Date.Format("2006-01-02"),
			Data: []aicodeassistanttypes.TimeSeriesDataPoint{
				{
					Key:   metricLabel,
					Value: result.LOCSuggestedCount,
				},
			},
		})
	}

	return dataPoints, nil
}

// CalculateActiveSessions calculates the active sessions metric
func (d *AICodeAssistantDB) CalculateActiveSessions(ctx context.Context, organizationID string, externalAccountIDs []string, toolNames []string, startDate, endDate time.Time, metricOperation types.MetricOperation) (*int, error) {
	selectStatement := ""
	switch metricOperation {
	case types.MetricOperationCount:
		selectStatement = "COALESCE(SUM(active_sessions), 0)"
	default:
		return nil, fmt.Errorf("invalid metric operation for active sessions: %s", metricOperation)
	}

	query := `
		SELECT ` + selectStatement + ` as active_sessions_count
		FROM ai_code_assistant_daily_metrics
		WHERE organization_id = ?
		AND metric_date >= ?
		AND metric_date <= ?
	`

	var args []any
	args = append(args, organizationID, startDate, endDate)

	if len(externalAccountIDs) > 0 {
		query += " AND external_account_id IN ?"
		args = append(args, externalAccountIDs)
	}

	if len(toolNames) > 0 {
		query += " AND tool_name IN ?"
		args = append(args, toolNames)
	}

	var count int64
	if err := d.db.WithContext(ctx).Raw(query, args...).Scan(&count).Error; err != nil {
		return nil, err
	}

	value := int(count)
	return &value, nil
}

// CalculateActiveSessionsGraph calculates the active sessions metric for a graph
func (d *AICodeAssistantDB) CalculateActiveSessionsGraph(ctx context.Context, organizationID string, externalAccountIDs []string, toolNames []string, startDate, endDate time.Time, metricOperation types.MetricOperation, metricLabel string, interval string) ([]aicodeassistanttypes.TimeSeriesEntry, error) {
	selectStatement := ""
	switch metricOperation {
	case types.MetricOperationCount:
		selectStatement = "COALESCE(SUM(active_sessions), 0)"
	default:
		return nil, fmt.Errorf("invalid metric operation for active sessions: %s", metricOperation)
	}

	// Map interval values to PostgreSQL DATE_TRUNC units
	postgresInterval := interval
	switch interval {
	case "daily":
		postgresInterval = "day"
	case "weekly":
		postgresInterval = "week"
	case "monthly":
		postgresInterval = "month"
	}

	query := `
		SELECT 
			DATE_TRUNC('` + postgresInterval + `', metric_date) as date,
			` + selectStatement + ` as active_sessions_count
		FROM ai_code_assistant_daily_metrics
		WHERE organization_id = ?
		AND metric_date >= ?
		AND metric_date <= ?
	`

	var args []any
	args = append(args, organizationID, startDate, endDate)

	if len(externalAccountIDs) > 0 {
		query += " AND external_account_id IN ?"
		args = append(args, externalAccountIDs)
	}

	if len(toolNames) > 0 {
		query += " AND tool_name IN ?"
		args = append(args, toolNames)
	}

	query += " GROUP BY DATE_TRUNC('" + postgresInterval + "', metric_date)"
	query += " ORDER BY date"

	var result struct {
		Date               time.Time `json:"date"`
		ActiveSessionsCount float64   `json:"active_sessions_count"`
	}

	rows, err := d.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dataPoints := []aicodeassistanttypes.TimeSeriesEntry{}
	for rows.Next() {
		if err := rows.Scan(&result.Date, &result.ActiveSessionsCount); err != nil {
			return nil, err
		}
		dataPoints = append(dataPoints, aicodeassistanttypes.TimeSeriesEntry{
			Date: result.Date.Format("2006-01-02"),
			Data: []aicodeassistanttypes.TimeSeriesDataPoint{
				{
					Key:   metricLabel,
					Value: result.ActiveSessionsCount,
				},
			},
		})
	}

	return dataPoints, nil
}

// CalculateAcceptRate calculates the accept rate metric
func (d *AICodeAssistantDB) CalculateAcceptRate(ctx context.Context, organizationID string, externalAccountIDs []string, toolNames []string, startDate, endDate time.Time, metricOperation types.MetricOperation) (*float64, error) {
	selectStatement := ""
	switch metricOperation {
	case types.MetricOperationAverage:
		selectStatement = "AVG(CASE WHEN lines_of_code_suggested > 0 THEN (lines_of_code_accepted::float / lines_of_code_suggested::float) * 100 ELSE 0 END)"
	default:
		return nil, fmt.Errorf("invalid metric operation for accept rate: %s", metricOperation)
	}

	query := `
		SELECT ` + selectStatement + ` as accept_rate
		FROM ai_code_assistant_daily_metrics
		WHERE organization_id = ?
		AND metric_date >= ?
		AND metric_date <= ?
	`

	var args []any
	args = append(args, organizationID, startDate, endDate)

	if len(externalAccountIDs) > 0 {
		query += " AND external_account_id IN ?"
		args = append(args, externalAccountIDs)
	}

	if len(toolNames) > 0 {
		query += " AND tool_name IN ?"
		args = append(args, toolNames)
	}

	var result struct {
		AcceptRate *float64 `json:"accept_rate"`
	}

	if err := d.db.WithContext(ctx).Raw(query, args...).Scan(&result).Error; err != nil {
		return nil, err
	}

	value := 0.0
	if result.AcceptRate != nil {
		value = *result.AcceptRate
	}

	return &value, nil
}

// CalculateAcceptRateGraph calculates the accept rate metric for a graph
func (d *AICodeAssistantDB) CalculateAcceptRateGraph(ctx context.Context, organizationID string, externalAccountIDs []string, toolNames []string, startDate, endDate time.Time, metricOperation types.MetricOperation, metricLabel string, interval string) ([]aicodeassistanttypes.TimeSeriesEntry, error) {
	selectStatement := ""
	switch metricOperation {
	case types.MetricOperationAverage:
		selectStatement = "AVG(CASE WHEN lines_of_code_suggested > 0 THEN (lines_of_code_accepted::float / lines_of_code_suggested::float) * 100 ELSE 0 END)"
	default:
		return nil, fmt.Errorf("invalid metric operation for accept rate: %s", metricOperation)
	}

	// Map interval values to PostgreSQL DATE_TRUNC units
	postgresInterval := interval
	switch interval {
	case "daily":
		postgresInterval = "day"
	case "weekly":
		postgresInterval = "week"
	case "monthly":
		postgresInterval = "month"
	}

	query := `
		SELECT 
			DATE_TRUNC('` + postgresInterval + `', metric_date) as date,
			` + selectStatement + ` as accept_rate
		FROM ai_code_assistant_daily_metrics
		WHERE organization_id = ?
		AND metric_date >= ?
		AND metric_date <= ?
	`

	var args []any
	args = append(args, organizationID, startDate, endDate)

	if len(externalAccountIDs) > 0 {
		query += " AND external_account_id IN ?"
		args = append(args, externalAccountIDs)
	}

	if len(toolNames) > 0 {
		query += " AND tool_name IN ?"
		args = append(args, toolNames)
	}

	query += " GROUP BY DATE_TRUNC('" + postgresInterval + "', metric_date)"
	query += " ORDER BY date"

	var result struct {
		Date       time.Time `json:"date"`
		AcceptRate float64   `json:"accept_rate"`
	}

	rows, err := d.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dataPoints := []aicodeassistanttypes.TimeSeriesEntry{}
	for rows.Next() {
		if err := rows.Scan(&result.Date, &result.AcceptRate); err != nil {
			return nil, err
		}
		dataPoints = append(dataPoints, aicodeassistanttypes.TimeSeriesEntry{
			Date: result.Date.Format("2006-01-02"),
			Data: []aicodeassistanttypes.TimeSeriesDataPoint{
				{
					Key:   metricLabel,
					Value: result.AcceptRate,
				},
			},
		})
	}

	return dataPoints, nil
}

// CalculateLinesOfCodeAcceptedForAccounts calculates the median lines of code accepted across accounts
func (d *AICodeAssistantDB) CalculateLinesOfCodeAcceptedForAccounts(ctx context.Context, organizationID string, externalAccountIDs []string, toolNames []string, startDate, endDate time.Time) (*float64, error) {
	query := `
		SELECT PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY member_total) as peer_loc_accepted
		FROM (
			SELECT external_account_id, COALESCE(SUM(lines_of_code_accepted), 0) as member_total
			FROM ai_code_assistant_daily_metrics
			WHERE organization_id = ?
			AND metric_date >= ?
			AND metric_date <= ?
	`

	var args []any
	args = append(args, organizationID, startDate, endDate)

	if len(externalAccountIDs) > 0 {
		query += "			AND external_account_id IN ?"
		args = append(args, externalAccountIDs)
	}

	if len(toolNames) > 0 {
		query += "			AND tool_name IN ?"
		args = append(args, toolNames)
	}

	query += `
			GROUP BY external_account_id
		) member_totals
	`

	var result *float64
	if err := d.db.WithContext(ctx).Raw(query, args...).Scan(&result).Error; err != nil {
		return nil, err
	}

	value := 0.0
	if result != nil {
		value = *result
	}

	return &value, nil
}

// CalculateLinesOfCodeSuggestedForAccounts calculates the median lines of code suggested across accounts
func (d *AICodeAssistantDB) CalculateLinesOfCodeSuggestedForAccounts(ctx context.Context, organizationID string, externalAccountIDs []string, toolNames []string, startDate, endDate time.Time) (*float64, error) {
	query := `
		SELECT PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY member_total) as peer_loc_suggested
		FROM (
			SELECT external_account_id, COALESCE(SUM(lines_of_code_suggested), 0) as member_total
			FROM ai_code_assistant_daily_metrics
			WHERE organization_id = ?
			AND metric_date >= ?
			AND metric_date <= ?
	`

	var args []any
	args = append(args, organizationID, startDate, endDate)

	if len(externalAccountIDs) > 0 {
		query += "			AND external_account_id IN ?"
		args = append(args, externalAccountIDs)
	}

	if len(toolNames) > 0 {
		query += "			AND tool_name IN ?"
		args = append(args, toolNames)
	}

	query += `
			GROUP BY external_account_id
		) member_totals
	`

	var result *float64
	if err := d.db.WithContext(ctx).Raw(query, args...).Scan(&result).Error; err != nil {
		return nil, err
	}

	value := 0.0
	if result != nil {
		value = *result
	}

	return &value, nil
}

// CalculateActiveSessionsForAccounts calculates the median active sessions across accounts
func (d *AICodeAssistantDB) CalculateActiveSessionsForAccounts(ctx context.Context, organizationID string, externalAccountIDs []string, toolNames []string, startDate, endDate time.Time) (*float64, error) {
	query := `
		SELECT PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY member_total) as peer_active_sessions
		FROM (
			SELECT external_account_id, COALESCE(SUM(active_sessions), 0) as member_total
			FROM ai_code_assistant_daily_metrics
			WHERE organization_id = ?
			AND metric_date >= ?
			AND metric_date <= ?
	`

	var args []any
	args = append(args, organizationID, startDate, endDate)

	if len(externalAccountIDs) > 0 {
		query += "			AND external_account_id IN ?"
		args = append(args, externalAccountIDs)
	}

	if len(toolNames) > 0 {
		query += "			AND tool_name IN ?"
		args = append(args, toolNames)
	}

	query += `
			GROUP BY external_account_id
		) member_totals
	`

	var result *float64
	if err := d.db.WithContext(ctx).Raw(query, args...).Scan(&result).Error; err != nil {
		return nil, err
	}

	value := 0.0
	if result != nil {
		value = *result
	}

	return &value, nil
}

// CalculateAcceptRateForAccounts calculates the median accept rate across accounts
func (d *AICodeAssistantDB) CalculateAcceptRateForAccounts(ctx context.Context, organizationID string, externalAccountIDs []string, toolNames []string, startDate, endDate time.Time) (*float64, error) {
	query := `
		SELECT PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY member_avg) as peer_accept_rate
		FROM (
			SELECT external_account_id, AVG(CASE WHEN lines_of_code_suggested > 0 THEN (lines_of_code_accepted::float / lines_of_code_suggested::float) * 100 ELSE 0 END) as member_avg
			FROM ai_code_assistant_daily_metrics
			WHERE organization_id = ?
			AND metric_date >= ?
			AND metric_date <= ?
	`

	var args []any
	args = append(args, organizationID, startDate, endDate)

	if len(externalAccountIDs) > 0 {
		query += "			AND external_account_id IN ?"
		args = append(args, externalAccountIDs)
	}

	if len(toolNames) > 0 {
		query += "			AND tool_name IN ?"
		args = append(args, toolNames)
	}

	query += `
			GROUP BY external_account_id
		) member_averages
	`

	var result *float64
	if err := d.db.WithContext(ctx).Raw(query, args...).Scan(&result).Error; err != nil {
		return nil, err
	}

	value := 0.0
	if result != nil {
		value = *result
	}

	return &value, nil
}

