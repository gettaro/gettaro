package database

import (
	"context"
	"fmt"
	"time"

	metrictypes "ems.dev/backend/services/sourcecontrol/metrics/types"
	"ems.dev/backend/services/sourcecontrol/types"
	"gorm.io/gorm"
)

// SourceControlDB defines the interface for source control database operations
type DB interface {
	// Source Control Accounts
	GetSourceControlAccounts(ctx context.Context, params *types.SourceControlAccountParams) ([]types.SourceControlAccount, error)
	CreateSourceControlAccounts(ctx context.Context, accounts []*types.SourceControlAccount) error
	GetSourceControlAccount(ctx context.Context, id string) (*types.SourceControlAccount, error)
	UpdateSourceControlAccount(ctx context.Context, account *types.SourceControlAccount) error

	// Pull Requests
	GetPullRequests(ctx context.Context, params *types.PullRequestParams) ([]*types.PullRequest, error)
	CreatePullRequest(ctx context.Context, pr *types.PullRequest) (*types.PullRequest, error)
	UpdatePullRequest(ctx context.Context, pr *types.PullRequest) error

	// Comments
	CreatePRComments(ctx context.Context, comments []*types.PRComment) error
	GetPullRequestComments(ctx context.Context, prID string) ([]*types.PRComment, error)

	// Member Activity
	GetMemberActivity(ctx context.Context, params *types.MemberActivityParams) ([]*types.MemberActivity, error)
	GetMemberPullRequests(ctx context.Context, params *types.MemberPullRequestParams) ([]*types.PullRequest, error)
	GetMemberPullRequestReviews(ctx context.Context, params *types.MemberPullRequestReviewsParams) ([]*types.MemberActivity, error)

	// Calculate time to merge metrics
	CalculateTimeToMerge(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation) (*int, error)
	CalculateTimeToMergeGraph(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation, metricLabel string, interval string) ([]types.TimeSeriesEntry, error)

	// Calculate PRs merged metrics
	CalculatePRsMerged(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation) (*int, error)
	CalculatePRsMergedGraph(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation, metricLabel string, interval string) ([]types.TimeSeriesEntry, error)

	// Calculate PRs reviewed metrics
	CalculatePRsReviewed(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation) (*int, error)
	CalculatePRsReviewedGraph(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation, metricLabel string, interval string) ([]types.TimeSeriesEntry, error)

	// Calculate LOC metrics
	CalculateLOCAdded(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation) (*int, error)
	CalculateLOCAddedGraph(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation, metricLabel string, interval string) ([]types.TimeSeriesEntry, error)
	CalculateLOCRemoved(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation) (*int, error)
	CalculateLOCRemovedGraph(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation, metricLabel string, interval string) ([]types.TimeSeriesEntry, error)

	// Calculate PR Review Complexity metrics
	CalculatePRReviewComplexity(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation) (*float64, error)
	CalculatePRReviewComplexityGraph(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation, metricLabel string, interval string) ([]types.TimeSeriesEntry, error)

	// Calculate peer metrics (median across peers)
	CalculateLOCAddedForAccounts(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time) (*float64, error)
	CalculateLOCRemovedForAccounts(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time) (*float64, error)
	CalculatePRsMergedForAccounts(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time) (*float64, error)
	CalculatePRsReviewedForAccounts(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time) (*float64, error)
	CalculateTimeToMergeForAccounts(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time) (*float64, error)
	CalculatePRReviewComplexityForAccounts(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time) (*float64, error)
}

type SourceControlDB struct {
	db *gorm.DB
}

// NewSourceControlDB creates a new instance of the source control database
func NewSourceControlDB(db *gorm.DB) *SourceControlDB {
	return &SourceControlDB{
		db: db,
	}
}

// GetSourceControlAccounts retrieves source control accounts
func (d *SourceControlDB) GetSourceControlAccounts(ctx context.Context, params *types.SourceControlAccountParams) ([]types.SourceControlAccount, error) {
	var accounts []types.SourceControlAccount
	query := d.db.WithContext(ctx).Model(&types.SourceControlAccount{})

	// Query by source control account IDs if provided
	if len(params.SourceControlAccountIDs) > 0 {
		query = query.Where("id IN ?", params.SourceControlAccountIDs)
	}

	// Query by usernames if provided
	if len(params.Usernames) > 0 {
		query = query.Where("username IN ?", params.Usernames)
	}

	// Filter by organization ID if provided
	if params.OrganizationID != "" {
		query = query.Where("organization_id = ?", params.OrganizationID)
	}

	// Filter by member IDs if provided
	if len(params.MemberIDs) > 0 {
		query = query.Where("member_id IN ?", params.MemberIDs)
	}

	if err := query.Find(&accounts).Error; err != nil {
		return nil, err
	}

	return accounts, nil
}

// CreateSourceControlAccounts creates multiple source control accounts
func (d *SourceControlDB) CreateSourceControlAccounts(ctx context.Context, accounts []*types.SourceControlAccount) error {
	return d.db.WithContext(ctx).Create(accounts).Error
}

// GetPullRequests retrieves pull requests based on the given parameters
func (d *SourceControlDB) GetPullRequests(ctx context.Context, params *types.PullRequestParams) ([]*types.PullRequest, error) {
	var prs []types.PullRequest
	query := d.db.WithContext(ctx).Model(&types.PullRequest{})

	// Add JOINs if needed for filtering
	if params.OrganizationID != nil || len(params.UserIDs) > 0 {
		query = query.Joins(`
			JOIN source_control_accounts sca ON pull_requests.source_control_account_id = sca.id
		`)
	}

	// Add organization_members JOIN if user IDs filtering is needed
	if len(params.UserIDs) > 0 {
		query = query.Joins(`
			LEFT JOIN organization_members om ON sca.member_id = om.id
		`)
	}

	if len(params.ProviderIDs) > 0 {
		query = query.Where("pull_requests.provider_id IN ?", params.ProviderIDs)
	}
	if params.OrganizationID != nil {
		query = query.Where("sca.organization_id = ?", *params.OrganizationID)
	}
	if params.RepositoryName != "" {
		query = query.Where("repository_name = ?", params.RepositoryName)
	}
	// Add user IDs filter if provided - convert to member IDs
	if len(params.UserIDs) > 0 {
		query = query.Where("om.user_id IN ?", params.UserIDs)
	}
	// Add date range filters if provided
	if params.StartDate != nil {
		query = query.Where("created_at >= ?", params.StartDate)
	}
	if params.EndDate != nil {
		query = query.Where("created_at <= ?", params.EndDate)
	}

	if err := query.Find(&prs).Error; err != nil {
		return nil, err
	}

	result := make([]*types.PullRequest, len(prs))
	for i := range prs {
		result[i] = &prs[i]
	}
	return result, nil
}

// CreatePRComments creates multiple PR comments
func (d *SourceControlDB) CreatePRComments(ctx context.Context, comments []*types.PRComment) error {
	return d.db.WithContext(ctx).Create(comments).Error
}

// UpdatePullRequest updates an existing pull request
func (d *SourceControlDB) UpdatePullRequest(ctx context.Context, pr *types.PullRequest) error {
	return d.db.WithContext(ctx).Model(pr).Updates(pr).Error
}

// UpdateSourceControlAccount updates an existing source control account
func (d *SourceControlDB) UpdateSourceControlAccount(ctx context.Context, account *types.SourceControlAccount) error {
	// Use explicit field updates to handle nil values properly
	updates := map[string]interface{}{
		"member_id":       account.MemberID,
		"organization_id": account.OrganizationID,
		"provider_name":   account.ProviderName,
		"provider_id":     account.ProviderID,
		"username":        account.Username,
		"metadata":        account.Metadata,
		"last_synced_at":  account.LastSyncedAt,
	}

	return d.db.WithContext(ctx).Model(account).Updates(updates).Error
}

// GetSourceControlAccount retrieves a source control account by ID
func (d *SourceControlDB) GetSourceControlAccount(ctx context.Context, id string) (*types.SourceControlAccount, error) {
	var account types.SourceControlAccount
	err := d.db.WithContext(ctx).First(&account, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetMemberActivity retrieves a timeline of source control activities for a specific member
func (d *SourceControlDB) GetMemberActivity(ctx context.Context, params *types.MemberActivityParams) ([]*types.MemberActivity, error) {
	var activities []*types.MemberActivity

	// Build date filter conditions for PRs
	var prDateFilter string
	var commentDateFilter string
	var args []interface{}
	args = append(args, params.MemberID)

	if params.StartDate != nil {
		prDateFilter += " AND pr.created_at >= ?"
		commentDateFilter += " AND pc.created_at >= ?"
		args = append(args, params.StartDate)
	}
	if params.EndDate != nil {
		prDateFilter += " AND pr.created_at <= ?"
		commentDateFilter += " AND pc.created_at <= ?"
		args = append(args, params.EndDate)
	}

	// Get pull requests created by the member
	prQuery := `
		SELECT 
			pr.id,
			'pull_request' as type,
			pr.title,
			pr.metrics,
			COALESCE(pr.description, '') as description,
			COALESCE(pr.url, '') as url,
			COALESCE(pr.repository_name, '') as repository,
			pr.created_at,
			COALESCE(pr.metadata, '{}'::jsonb) as metadata,
			sca.username as author_username,
			NULL as pr_title,
			NULL as pr_author_username,
			pr.metrics as pr_metrics
		FROM pull_requests pr
		JOIN source_control_accounts sca ON pr.source_control_account_id = sca.id
		WHERE sca.member_id = ? AND pr.merged_at is not null` + prDateFilter + `
		ORDER BY pr.created_at DESC
	`

	var prs []types.MemberActivity
	if err := d.db.WithContext(ctx).Raw(prQuery, args...).Scan(&prs).Error; err != nil {
		return nil, err
	}

	// Get comments and reviews by the member on other people's PRs
	commentQuery := `
		SELECT 
			pc.id,
			CASE 
				WHEN pc.type = 'REVIEW' THEN 'pr_review'
				ELSE 'pr_comment'
			END as type,
			pc.body as description,
			pr.url as url,
			pr.repository_name as repository,
			pc.created_at,
			'{}'::jsonb as metadata,
			sca.username as author_username,
			pr.title as pr_title,
			pr_author.username as pr_author_username,
			NULL as pr_metrics
		FROM pr_comments pc
		JOIN pull_requests pr ON pc.pr_id = pr.id
		JOIN source_control_accounts sca ON pc.source_control_account_id = sca.id
		JOIN source_control_accounts pr_author ON pr.source_control_account_id = pr_author.id
		WHERE sca.member_id = ? 
		AND sca.id != pr.source_control_account_id` + commentDateFilter + `
		ORDER BY pc.created_at DESC
	`

	var comments []types.MemberActivity
	if err := d.db.WithContext(ctx).Raw(commentQuery, args...).Scan(&comments).Error; err != nil {
		return nil, err
	}

	// Convert to pointers and combine
	for i := range prs {
		activities = append(activities, &prs[i])
	}
	for i := range comments {
		activities = append(activities, &comments[i])
	}

	return activities, nil
}

// GetMemberPullRequests retrieves pull requests for a specific member
func (d *SourceControlDB) GetMemberPullRequests(ctx context.Context, params *types.MemberPullRequestParams) ([]*types.PullRequest, error) {
	var prs []types.PullRequest
	query := d.db.WithContext(ctx).Model(&types.PullRequest{})

	// Join with source_control_accounts to filter by member_id
	query = query.Joins(`
		JOIN source_control_accounts sca ON pull_requests.source_control_account_id = sca.id
	`)

	// Filter by member ID
	query = query.Where("sca.member_id = ?", params.MemberID)

	// Add date range filters if provided
	if params.StartDate != nil {
		query = query.Where("pull_requests.created_at >= ?", params.StartDate)
	}
	if params.EndDate != nil {
		query = query.Where("pull_requests.created_at <= ?", params.EndDate)
	}

	// Order by created_at descending
	query = query.Order("pull_requests.created_at DESC")

	if err := query.Find(&prs).Error; err != nil {
		return nil, err
	}

	result := make([]*types.PullRequest, len(prs))
	for i := range prs {
		result[i] = &prs[i]
	}
	return result, nil
}

// GetMemberPullRequestReviews retrieves pull request reviews for a specific member
func (d *SourceControlDB) GetMemberPullRequestReviews(ctx context.Context, params *types.MemberPullRequestReviewsParams) ([]*types.MemberActivity, error) {
	var activities []*types.MemberActivity

	// Build date filter conditions for comments
	var commentDateFilter string
	var args []interface{}
	args = append(args, params.MemberID)

	if params.StartDate != nil {
		commentDateFilter += " AND pc.created_at >= ?"
		args = append(args, params.StartDate)
	}
	if params.EndDate != nil {
		commentDateFilter += " AND pc.created_at <= ?"
		args = append(args, params.EndDate)
	}

	// Get comments and reviews by the member on other people's PRs
	commentQuery := `
		SELECT 
			pc.id,
			CASE 
				WHEN pc.type = 'REVIEW' THEN 'pr_review'
				ELSE 'pr_comment'
			END as type,
			pc.body as description,
			pr.url as url,
			pr.repository_name as repository,
			pc.created_at,
			'{}'::jsonb as metadata,
			sca.username as author_username,
			pr.title as pr_title,
			pr_author.username as pr_author_username,
			NULL as pr_metrics
		FROM pr_comments pc
		JOIN pull_requests pr ON pc.pr_id = pr.id
		JOIN source_control_accounts sca ON pc.source_control_account_id = sca.id
		JOIN source_control_accounts pr_author ON pr.source_control_account_id = pr_author.id
		WHERE sca.member_id = ? 
		AND sca.id != pr.source_control_account_id` + commentDateFilter + `
		ORDER BY pc.created_at DESC
	`

	var comments []types.MemberActivity
	if err := d.db.WithContext(ctx).Raw(commentQuery, args...).Scan(&comments).Error; err != nil {
		return nil, err
	}

	// Convert to pointers
	for i := range comments {
		activities = append(activities, &comments[i])
	}

	return activities, nil
}

// GetPullRequestComments retrieves all comments for a specific pull request
func (d *SourceControlDB) GetPullRequestComments(ctx context.Context, prID string) ([]*types.PRComment, error) {
	var comments []types.PRComment
	err := d.db.WithContext(ctx).Where("pr_id = ?", prID).Find(&comments).Error
	if err != nil {
		return nil, err
	}

	result := make([]*types.PRComment, len(comments))
	for i := range comments {
		result[i] = &comments[i]
	}
	return result, nil
}

// CreatePullRequest creates a single pull request and returns it
func (d *SourceControlDB) CreatePullRequest(ctx context.Context, pr *types.PullRequest) (*types.PullRequest, error) {
	if err := d.db.WithContext(ctx).Create(pr).Error; err != nil {
		return nil, err
	}
	return pr, nil
}

// CalculateTimeToMerge calculates the time to merge metric
func (d *SourceControlDB) CalculateTimeToMerge(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation) (*int, error) {
	selectStatement := ""
	switch metricOperation {
	case metrictypes.MetricOperationMedian:
		selectStatement = "PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY EXTRACT(EPOCH FROM (pr.merged_at - pr.created_at)))"
	case metrictypes.MetricOperationAverage:
		selectStatement = "AVG(EXTRACT(EPOCH FROM (pr.merged_at - pr.created_at)))"
	default:
		return nil, fmt.Errorf("invalid metric operation: %s", metricOperation)
	}

	query := `
		SELECT ` + selectStatement + ` as time_to_merge_seconds
		FROM pull_requests pr
		JOIN source_control_accounts sca ON pr.source_control_account_id = sca.id
		WHERE sca.organization_id = ?
		AND pr.created_at >= ?
		AND pr.created_at <= ?
		AND pr.merged_at IS NOT NULL
		AND pr.status = 'closed'
	`

	var args []any
	args = append(args, organizationID, startDate, endDate)

	if len(sourceControlAccountIDs) > 0 {
		query += " AND pr.source_control_account_id IN ?"
		args = append(args, sourceControlAccountIDs)
	}

	var result struct {
		TimeToMergeSeconds *float64 `json:"time_to_merge_seconds"`
	}

	if err := d.db.WithContext(ctx).Raw(query, args...).Scan(&result).Error; err != nil {
		return nil, err
	}

	// If no results found, return 0
	value := 0
	if result.TimeToMergeSeconds != nil {
		value = int(*result.TimeToMergeSeconds)
	}

	return &value, nil
}

// CalculateTimeToMergeGraph calculates the time to merge metric for a graph
func (d *SourceControlDB) CalculateTimeToMergeGraph(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation, metricLabel string, interval string) ([]types.TimeSeriesEntry, error) {
	selectStatement := ""
	switch metricOperation {
	case metrictypes.MetricOperationMedian:
		selectStatement = "PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY EXTRACT(EPOCH FROM (pr.merged_at - pr.created_at)))"
	case metrictypes.MetricOperationAverage:
		selectStatement = "AVG(EXTRACT(EPOCH FROM (pr.merged_at - pr.created_at)))"
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
			DATE_TRUNC('` + postgresInterval + `', pr.merged_at) as date,
			` + selectStatement + ` as time_to_merge_seconds
		FROM pull_requests pr
		JOIN source_control_accounts sca ON pr.source_control_account_id = sca.id
		WHERE sca.organization_id = ?
		AND pr.created_at >= ?
		AND pr.created_at <= ?
		AND pr.merged_at IS NOT NULL
		AND pr.status = 'closed'
	`

	var args []any
	args = append(args, organizationID, startDate, endDate)

	if len(sourceControlAccountIDs) > 0 {
		query += " AND pr.source_control_account_id IN ?"
		args = append(args, sourceControlAccountIDs)
	}

	// Add GROUP BY clause for the DATE_TRUNC grouping
	query += " GROUP BY DATE_TRUNC('" + postgresInterval + "', pr.merged_at)"

	// Add ORDER BY to ensure consistent results
	query += " ORDER BY date"

	var result struct {
		Date             time.Time `json:"date"`
		TimeToMergeValue float64   `json:"time_to_merge_seconds"`
	}

	rows, err := d.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dataPoints := []types.TimeSeriesEntry{}
	for rows.Next() {
		if err := rows.Scan(&result.Date, &result.TimeToMergeValue); err != nil {
			return nil, err
		}
		dataPoints = append(dataPoints, types.TimeSeriesEntry{
			Date: result.Date.Format("2006-01-02"),
			Data: []types.TimeSeriesDataPoint{
				{
					Key:   metricLabel,
					Value: result.TimeToMergeValue,
				},
			},
		})
	}

	return dataPoints, nil
}

// CalculatePRsMerged calculates the PRs merged metric
func (d *SourceControlDB) CalculatePRsMerged(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation) (*int, error) {
	selectStatement := ""
	switch metricOperation {
	case metrictypes.MetricOperationCount:
		selectStatement = "COUNT(*)"
	default:
		return nil, fmt.Errorf("invalid metric operation for PRs merged: %s", metricOperation)
	}

	query := `
		SELECT ` + selectStatement + ` as prs_merged_count
		FROM pull_requests pr
		JOIN source_control_accounts sca ON pr.source_control_account_id = sca.id
		WHERE sca.organization_id = ?
		AND pr.created_at >= ?
		AND pr.created_at <= ?
		AND pr.merged_at IS NOT NULL
		AND pr.status = 'closed'
	`

	var args []any
	args = append(args, organizationID, startDate, endDate)

	if len(sourceControlAccountIDs) > 0 {
		query += " AND pr.source_control_account_id IN ?"
		args = append(args, sourceControlAccountIDs)
	}

	var count int64
	if err := d.db.WithContext(ctx).Raw(query, args...).Scan(&count).Error; err != nil {
		return nil, err
	}

	// Convert int64 to int
	value := int(count)

	return &value, nil
}

// CalculatePRsMergedGraph calculates the PRs merged metric for a graph
func (d *SourceControlDB) CalculatePRsMergedGraph(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation, metricLabel string, interval string) ([]types.TimeSeriesEntry, error) {
	selectStatement := ""
	switch metricOperation {
	case metrictypes.MetricOperationCount:
		selectStatement = "COUNT(*)"
	default:
		return nil, fmt.Errorf("invalid metric operation for PRs merged: %s", metricOperation)
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
			DATE_TRUNC('` + postgresInterval + `', pr.merged_at) as date,
			` + selectStatement + ` as prs_merged_count
		FROM pull_requests pr
		JOIN source_control_accounts sca ON pr.source_control_account_id = sca.id
		WHERE sca.organization_id = ?
		AND pr.created_at >= ?
		AND pr.created_at <= ?
		AND pr.merged_at IS NOT NULL
		AND pr.status = 'closed'
	`

	var args []any
	args = append(args, organizationID, startDate, endDate)

	if len(sourceControlAccountIDs) > 0 {
		query += " AND pr.source_control_account_id IN ?"
		args = append(args, sourceControlAccountIDs)
	}

	// Add GROUP BY clause for the DATE_TRUNC grouping
	query += " GROUP BY DATE_TRUNC('" + postgresInterval + "', pr.merged_at)"

	// Add ORDER BY to ensure consistent results
	query += " ORDER BY date"

	var result struct {
		Date           time.Time `json:"date"`
		PRsMergedCount float64   `json:"prs_merged_count"`
	}

	rows, err := d.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dataPoints := []types.TimeSeriesEntry{}
	for rows.Next() {
		if err := rows.Scan(&result.Date, &result.PRsMergedCount); err != nil {
			return nil, err
		}
		dataPoints = append(dataPoints, types.TimeSeriesEntry{
			Date: result.Date.Format("2006-01-02"),
			Data: []types.TimeSeriesDataPoint{
				{
					Key:   metricLabel,
					Value: result.PRsMergedCount,
				},
			},
		})
	}

	return dataPoints, nil
}

// CalculatePRsReviewed calculates the PRs reviewed metric
func (d *SourceControlDB) CalculatePRsReviewed(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation) (*int, error) {
	selectStatement := ""
	switch metricOperation {
	case metrictypes.MetricOperationCount:
		selectStatement = "COUNT(DISTINCT pr.id)"
	default:
		return nil, fmt.Errorf("invalid metric operation for PRs reviewed: %s", metricOperation)
	}

	query := `
		SELECT ` + selectStatement + ` as prs_reviewed_count
		FROM pull_requests pr
		WHERE pr.created_at >= ?
		AND pr.created_at <= ?
		AND EXISTS (
			SELECT 1 FROM pr_comments pc 
			JOIN source_control_accounts sca ON pc.source_control_account_id = sca.id
			WHERE pc.pr_id = pr.id 
			AND sca.organization_id = ?
			AND sca.id IN ?
		)
	`

	var args []any
	args = append(args, startDate, endDate, organizationID, sourceControlAccountIDs)

	var count int64
	if err := d.db.WithContext(ctx).Raw(query, args...).Scan(&count).Error; err != nil {
		return nil, err
	}

	// Convert int64 to int
	value := int(count)

	return &value, nil
}

// CalculatePRsReviewedGraph calculates the PRs reviewed metric for a graph
func (d *SourceControlDB) CalculatePRsReviewedGraph(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation, metricLabel string, interval string) ([]types.TimeSeriesEntry, error) {
	selectStatement := ""
	switch metricOperation {
	case metrictypes.MetricOperationCount:
		selectStatement = "COUNT(DISTINCT pr.id)"
	default:
		return nil, fmt.Errorf("invalid metric operation for PRs reviewed: %s", metricOperation)
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
			DATE_TRUNC('` + postgresInterval + `', pr.created_at) as date,
			` + selectStatement + ` as prs_reviewed_count
		FROM pull_requests pr
		WHERE pr.created_at >= ?
		AND pr.created_at <= ?
		AND EXISTS (
			SELECT 1 FROM pr_comments pc 
			JOIN source_control_accounts sca ON pc.source_control_account_id = sca.id
			WHERE pc.pr_id = pr.id 
			AND sca.organization_id = ?
			AND sca.id IN ?
		)
		GROUP BY DATE_TRUNC('` + postgresInterval + `', pr.created_at)
		ORDER BY date
	`

	var args []any
	args = append(args, startDate, endDate, organizationID, sourceControlAccountIDs)

	var result struct {
		Date             time.Time `json:"date"`
		PRsReviewedCount float64   `json:"prs_reviewed_count"`
	}

	rows, err := d.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dataPoints := []types.TimeSeriesEntry{}
	for rows.Next() {
		if err := rows.Scan(&result.Date, &result.PRsReviewedCount); err != nil {
			return nil, err
		}
		dataPoints = append(dataPoints, types.TimeSeriesEntry{
			Date: result.Date.Format("2006-01-02"),
			Data: []types.TimeSeriesDataPoint{
				{
					Key:   metricLabel,
					Value: result.PRsReviewedCount,
				},
			},
		})
	}

	return dataPoints, nil
}

// CalculateLOCAdded calculates the lines of code added metric
func (d *SourceControlDB) CalculateLOCAdded(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation) (*int, error) {
	selectStatement := ""
	switch metricOperation {
	case metrictypes.MetricOperationCount:
		selectStatement = "COALESCE(SUM(CAST(pr.metadata->>'additions' AS BIGINT)), 0)"
	default:
		return nil, fmt.Errorf("invalid metric operation for LOC added: %s", metricOperation)
	}

	query := `
		SELECT ` + selectStatement + ` as loc_added_count
		FROM pull_requests pr
		JOIN source_control_accounts sca ON pr.source_control_account_id = sca.id
		WHERE sca.organization_id = ?
		AND pr.created_at >= ?
		AND pr.created_at <= ?
		AND pr.merged_at IS NOT NULL
		AND pr.status = 'closed'
	`

	var args []any
	args = append(args, organizationID, startDate, endDate)

	if len(sourceControlAccountIDs) > 0 {
		query += " AND pr.source_control_account_id IN ?"
		args = append(args, sourceControlAccountIDs)
	}

	var count int64
	if err := d.db.WithContext(ctx).Raw(query, args...).Scan(&count).Error; err != nil {
		return nil, err
	}

	// Convert int64 to int
	value := int(count)

	return &value, nil
}

// CalculateLOCAddedGraph calculates the lines of code added metric for a graph
func (d *SourceControlDB) CalculateLOCAddedGraph(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation, metricLabel string, interval string) ([]types.TimeSeriesEntry, error) {
	selectStatement := ""
	switch metricOperation {
	case metrictypes.MetricOperationCount:
		selectStatement = "COALESCE(SUM(CAST(pr.metadata->>'additions' AS BIGINT)), 0)"
	default:
		return nil, fmt.Errorf("invalid metric operation for LOC added: %s", metricOperation)
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
			DATE_TRUNC('` + postgresInterval + `', pr.merged_at) as date,
			` + selectStatement + ` as loc_added_count
		FROM pull_requests pr
		JOIN source_control_accounts sca ON pr.source_control_account_id = sca.id
		WHERE sca.organization_id = ?
		AND pr.created_at >= ?
		AND pr.created_at <= ?
		AND pr.merged_at IS NOT NULL
		AND pr.status = 'closed'
	`

	var args []any
	args = append(args, organizationID, startDate, endDate)

	if len(sourceControlAccountIDs) > 0 {
		query += " AND pr.source_control_account_id IN ?"
		args = append(args, sourceControlAccountIDs)
	}

	// Add GROUP BY clause for the DATE_TRUNC grouping
	query += " GROUP BY DATE_TRUNC('" + postgresInterval + "', pr.merged_at)"

	// Add ORDER BY to ensure consistent results
	query += " ORDER BY date"

	var result struct {
		Date          time.Time `json:"date"`
		LOCAddedCount float64   `json:"loc_added_count"`
	}

	rows, err := d.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dataPoints := []types.TimeSeriesEntry{}
	for rows.Next() {
		if err := rows.Scan(&result.Date, &result.LOCAddedCount); err != nil {
			return nil, err
		}
		dataPoints = append(dataPoints, types.TimeSeriesEntry{
			Date: result.Date.Format("2006-01-02"),
			Data: []types.TimeSeriesDataPoint{
				{
					Key:   metricLabel,
					Value: result.LOCAddedCount,
				},
			},
		})
	}

	return dataPoints, nil
}

// CalculateLOCRemoved calculates the lines of code removed metric
func (d *SourceControlDB) CalculateLOCRemoved(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation) (*int, error) {
	selectStatement := ""
	switch metricOperation {
	case metrictypes.MetricOperationCount:
		selectStatement = "COALESCE(SUM(CAST(pr.metadata->>'deletions' AS BIGINT)), 0)"
	default:
		return nil, fmt.Errorf("invalid metric operation for LOC removed: %s", metricOperation)
	}

	query := `
		SELECT ` + selectStatement + ` as loc_removed_count
		FROM pull_requests pr
		JOIN source_control_accounts sca ON pr.source_control_account_id = sca.id
		WHERE sca.organization_id = ?
		AND pr.created_at >= ?
		AND pr.created_at <= ?
		AND pr.merged_at IS NOT NULL
		AND pr.status = 'closed'
	`

	var args []any
	args = append(args, organizationID, startDate, endDate)

	if len(sourceControlAccountIDs) > 0 {
		query += " AND pr.source_control_account_id IN ?"
		args = append(args, sourceControlAccountIDs)
	}

	var count int64
	if err := d.db.WithContext(ctx).Raw(query, args...).Scan(&count).Error; err != nil {
		return nil, err
	}

	// Convert int64 to int
	value := int(count)

	return &value, nil
}

// CalculateLOCRemovedGraph calculates the lines of code removed metric for a graph
func (d *SourceControlDB) CalculateLOCRemovedGraph(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation, metricLabel string, interval string) ([]types.TimeSeriesEntry, error) {
	selectStatement := ""
	switch metricOperation {
	case metrictypes.MetricOperationCount:
		selectStatement = "COALESCE(SUM(CAST(pr.metadata->>'deletions' AS BIGINT)), 0)"
	default:
		return nil, fmt.Errorf("invalid metric operation for LOC removed: %s", metricOperation)
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
			DATE_TRUNC('` + postgresInterval + `', pr.merged_at) as date,
			` + selectStatement + ` as loc_removed_count
		FROM pull_requests pr
		JOIN source_control_accounts sca ON pr.source_control_account_id = sca.id
		WHERE sca.organization_id = ?
		AND pr.created_at >= ?
		AND pr.created_at <= ?
		AND pr.merged_at IS NOT NULL
		AND pr.status = 'closed'
	`

	var args []any
	args = append(args, organizationID, startDate, endDate)

	if len(sourceControlAccountIDs) > 0 {
		query += " AND pr.source_control_account_id IN ?"
		args = append(args, sourceControlAccountIDs)
	}

	// Add GROUP BY clause for the DATE_TRUNC grouping
	query += " GROUP BY DATE_TRUNC('" + postgresInterval + "', pr.merged_at)"

	// Add ORDER BY to ensure consistent results
	query += " ORDER BY date"

	var result struct {
		Date            time.Time `json:"date"`
		LOCRemovedCount float64   `json:"loc_removed_count"`
	}

	rows, err := d.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dataPoints := []types.TimeSeriesEntry{}
	for rows.Next() {
		if err := rows.Scan(&result.Date, &result.LOCRemovedCount); err != nil {
			return nil, err
		}
		dataPoints = append(dataPoints, types.TimeSeriesEntry{
			Date: result.Date.Format("2006-01-02"),
			Data: []types.TimeSeriesDataPoint{
				{
					Key:   metricLabel,
					Value: result.LOCRemovedCount,
				},
			},
		})
	}

	return dataPoints, nil
}

// CalculatePRReviewComplexity calculates the PR review complexity metric
func (d *SourceControlDB) CalculatePRReviewComplexity(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation) (*float64, error) {
	selectStatement := ""
	switch metricOperation {
	case metrictypes.MetricOperationAverage:
		selectStatement = "AVG(pr.additions + pr.deletions)"
	default:
		return nil, fmt.Errorf("invalid metric operation for PR review complexity: %s", metricOperation)
	}

	// First, let's debug by checking what PRs exist
	debugQuery := `
		SELECT COUNT(*) as total_prs
		FROM pull_requests pr
		WHERE pr.created_at >= ?
		AND pr.created_at <= ?
		AND pr.merged_at IS NOT NULL
		AND pr.status = 'closed'
	`

	var debugArgs []any
	debugArgs = append(debugArgs, startDate, endDate)

	var totalPRs int64
	if err := d.db.WithContext(ctx).Raw(debugQuery, debugArgs...).Scan(&totalPRs).Error; err != nil {
		return nil, fmt.Errorf("debug query failed: %w", err)
	}

	fmt.Printf("DEBUG - Total PRs in date range: %d\n", totalPRs)

	// Let's also check what source control accounts we're looking for
	fmt.Printf("DEBUG - Looking for reviews from source control accounts: %v\n", sourceControlAccountIDs)

	// Simple test: just count PRs with any comments
	simpleTestQuery := `
		SELECT COUNT(*) as total_commented_prs
		FROM pull_requests pr
		WHERE pr.created_at >= ?
		AND pr.created_at <= ?
		AND pr.merged_at IS NOT NULL
		AND pr.status = 'closed'
		AND EXISTS (
			SELECT 1 FROM pr_comments pc 
			WHERE pc.pr_id = pr.id 
		)
	`

	var simpleArgs []any
	simpleArgs = append(simpleArgs, startDate, endDate)

	var totalCommentedPRs int64
	if err := d.db.WithContext(ctx).Raw(simpleTestQuery, simpleArgs...).Scan(&totalCommentedPRs).Error; err != nil {
		return nil, fmt.Errorf("simple test query failed: %w", err)
	}

	fmt.Printf("DEBUG - PRs with any comments: %d\n", totalCommentedPRs)

	// Very simple test: just count PRs commented on by the member
	verySimpleQuery := `
		SELECT COUNT(*) as member_commented_prs
		FROM pull_requests pr
		WHERE pr.created_at >= ?
		AND pr.created_at <= ?
		AND pr.merged_at IS NOT NULL
		AND pr.status = 'closed'
		AND EXISTS (
			SELECT 1 FROM pr_comments pc 
			JOIN source_control_accounts sca ON pc.source_control_account_id = sca.id
			WHERE pc.pr_id = pr.id 
			AND sca.organization_id = ?
			AND sca.id IN ?
		)
	`

	var verySimpleArgs []any
	verySimpleArgs = append(verySimpleArgs, startDate, endDate, organizationID, sourceControlAccountIDs)

	var memberCommentedPRs int64
	if err := d.db.WithContext(ctx).Raw(verySimpleQuery, verySimpleArgs...).Scan(&memberCommentedPRs).Error; err != nil {
		return nil, fmt.Errorf("very simple test query failed: %w", err)
	}

	fmt.Printf("DEBUG - PRs commented on by member: %d\n", memberCommentedPRs)

	// Now run the actual query - temporarily simplified for debugging
	query := `
		SELECT ` + selectStatement + ` as avg_review_complexity
		FROM pull_requests pr
		WHERE pr.created_at >= ?
		AND pr.created_at <= ?
		AND pr.merged_at IS NOT NULL
		AND pr.status = 'closed'
		AND EXISTS (
			SELECT 1 FROM pr_comments pc 
			JOIN source_control_accounts sca ON pc.source_control_account_id = sca.id
			WHERE pc.pr_id = pr.id 
			AND sca.organization_id = ?
			AND sca.id IN ?
		)
		-- Temporarily commented out for debugging
		-- AND pr.source_control_account_id NOT IN (
		-- 	SELECT sca.id FROM source_control_accounts sca 
		-- 	WHERE sca.organization_id = ?
		-- )
	`

	var args []any
	args = append(args, startDate, endDate, organizationID, sourceControlAccountIDs)
	// args = append(args, startDate, endDate, organizationID, sourceControlAccountIDs, organizationID)

	var result struct {
		AvgReviewComplexity *float64 `json:"avg_review_complexity"`
	}

	if err := d.db.WithContext(ctx).Raw(query, args...).Scan(&result).Error; err != nil {
		return nil, err
	}

	// If no results found, return 0
	value := 0.0
	if result.AvgReviewComplexity != nil {
		value = *result.AvgReviewComplexity
	}

	fmt.Printf("DEBUG - Final result: %f\n", value)
	return &value, nil
}

// CalculatePRReviewComplexityGraph calculates the PR review complexity metric for a graph
func (d *SourceControlDB) CalculatePRReviewComplexityGraph(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation, metricLabel string, interval string) ([]types.TimeSeriesEntry, error) {
	selectStatement := ""
	switch metricOperation {
	case metrictypes.MetricOperationAverage:
		selectStatement = "AVG(pr.additions + pr.deletions)"
	default:
		return nil, fmt.Errorf("invalid metric operation for PR review complexity: %s", metricOperation)
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
			DATE_TRUNC('` + postgresInterval + `', pr.created_at) as date,
			` + selectStatement + ` as avg_review_complexity
		FROM pull_requests pr
		WHERE pr.created_at >= ?
		AND pr.created_at <= ?
		AND pr.merged_at IS NOT NULL
		AND pr.status = 'closed'
		AND EXISTS (
			SELECT 1 FROM pr_comments pc 
			JOIN source_control_accounts sca ON pc.source_control_account_id = sca.id
			WHERE pc.pr_id = pr.id 
			AND sca.organization_id = ?
			AND sca.id IN ?
		)
		AND pr.source_control_account_id NOT IN (
			SELECT sca.id FROM source_control_accounts sca 
			WHERE sca.organization_id = ?
		)
		GROUP BY DATE_TRUNC('` + postgresInterval + `', pr.created_at)
		ORDER BY date
	`

	var args []any
	args = append(args, startDate, endDate, organizationID, sourceControlAccountIDs, organizationID)

	var result struct {
		Date                time.Time `json:"date"`
		AvgReviewComplexity float64   `json:"avg_review_complexity"`
	}

	rows, err := d.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dataPoints := []types.TimeSeriesEntry{}
	for rows.Next() {
		if err := rows.Scan(&result.Date, &result.AvgReviewComplexity); err != nil {
			return nil, err
		}
		dataPoints = append(dataPoints, types.TimeSeriesEntry{
			Date: result.Date.Format("2006-01-02"),
			Data: []types.TimeSeriesDataPoint{
				{
					Key:   metricLabel,
					Value: result.AvgReviewComplexity,
				},
			},
		})
	}

	return dataPoints, nil
}

// CalculateLOCAddedForAccounts calculates the median LOC added across accounts
func (d *SourceControlDB) CalculateLOCAddedForAccounts(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time) (*float64, error) {
	query := `
		SELECT PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY member_total) as peer_loc_added
		FROM (
			SELECT sca.member_id, COALESCE(SUM(CAST(pr.metadata->>'additions' AS BIGINT)), 0) as member_total
			FROM pull_requests pr
			JOIN source_control_accounts sca ON pr.source_control_account_id = sca.id
			WHERE sca.organization_id = ?
			AND pr.created_at >= ?
			AND pr.created_at <= ?
			AND pr.merged_at IS NOT NULL
			AND pr.status = 'closed'
			AND sca.id IN ?
			GROUP BY sca.member_id
		) member_totals
	`

	var args []any
	args = append(args, organizationID, startDate, endDate, sourceControlAccountIDs)

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

// CalculateLOCRemovedForAccounts calculates the median LOC removed across accounts
func (d *SourceControlDB) CalculateLOCRemovedForAccounts(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time) (*float64, error) {
	query := `
		SELECT PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY member_total) as peer_loc_removed
		FROM (
			SELECT sca.member_id, COALESCE(SUM(CAST(pr.metadata->>'deletions' AS BIGINT)), 0) as member_total
			FROM pull_requests pr
			JOIN source_control_accounts sca ON pr.source_control_account_id = sca.id
			WHERE sca.organization_id = ?
			AND pr.created_at >= ?
			AND pr.created_at <= ?
			AND pr.merged_at IS NOT NULL
			AND pr.status = 'closed'
			AND sca.id IN ?
			GROUP BY sca.member_id
		) member_totals
	`

	var args []any
	args = append(args, organizationID, startDate, endDate, sourceControlAccountIDs)

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

// CalculatePRsMergedForAccounts calculates the median PRs merged across accounts
func (d *SourceControlDB) CalculatePRsMergedForAccounts(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time) (*float64, error) {
	query := `
		SELECT PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY member_total) as peer_prs_merged
		FROM (
			SELECT sca.member_id, COUNT(*) as member_total
			FROM pull_requests pr
			JOIN source_control_accounts sca ON pr.source_control_account_id = sca.id
			WHERE sca.organization_id = ?
			AND pr.created_at >= ?
			AND pr.created_at <= ?
			AND pr.merged_at IS NOT NULL
			AND pr.status = 'closed'
			AND sca.id IN ?
			GROUP BY sca.member_id
		) member_totals
	`

	var args []any
	args = append(args, organizationID, startDate, endDate, sourceControlAccountIDs)

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

// CalculatePRsReviewedForAccounts calculates the median PRs reviewed across accounts
func (d *SourceControlDB) CalculatePRsReviewedForAccounts(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time) (*float64, error) {
	query := `
		SELECT PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY member_total) as peer_prs_reviewed
		FROM (
			SELECT sca.member_id, COUNT(DISTINCT pr.id) as member_total
			FROM pull_requests pr
			JOIN pr_comments pc ON pc.pr_id = pr.id
			JOIN source_control_accounts sca ON pc.source_control_account_id = sca.id
			WHERE sca.organization_id = ?
			AND pr.created_at >= ?
			AND pr.created_at <= ?
			AND sca.id IN ?
			GROUP BY sca.member_id
		) member_totals
	`

	var args []any
	args = append(args, organizationID, startDate, endDate, sourceControlAccountIDs)

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

// CalculateTimeToMergeForAccounts calculates the median time to merge across accounts
func (d *SourceControlDB) CalculateTimeToMergeForAccounts(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time) (*float64, error) {
	query := `
		SELECT PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY member_avg) as peer_time_to_merge
		FROM (
			SELECT sca.member_id, AVG(EXTRACT(EPOCH FROM (pr.merged_at - pr.created_at))) as member_avg
			FROM pull_requests pr
			JOIN source_control_accounts sca ON pr.source_control_account_id = sca.id
			WHERE sca.organization_id = ?
			AND pr.created_at >= ?
			AND pr.created_at <= ?
			AND pr.merged_at IS NOT NULL
			AND pr.status = 'closed'
			AND sca.id IN ?
			GROUP BY sca.member_id
		) member_averages
	`

	var args []any
	args = append(args, organizationID, startDate, endDate, sourceControlAccountIDs)

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

// CalculatePRReviewComplexityForAccounts calculates the median PR review complexity across accounts
func (d *SourceControlDB) CalculatePRReviewComplexityForAccounts(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time) (*float64, error) {
	query := `
		SELECT PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY member_avg) as peer_pr_review_complexity
		FROM (
			SELECT sca.member_id, AVG(pr.additions + pr.deletions) as member_avg
			FROM pull_requests pr
			JOIN pr_comments pc ON pc.pr_id = pr.id
			JOIN source_control_accounts sca ON pc.source_control_account_id = sca.id
			WHERE sca.organization_id = ?
			AND pr.created_at >= ?
			AND pr.created_at <= ?
			AND pc.type = 'REVIEW'
			AND sca.id IN ?
			GROUP BY sca.member_id
		) member_averages
	`

	var args []any
	args = append(args, organizationID, startDate, endDate, sourceControlAccountIDs)

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
