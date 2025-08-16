package database

import (
	"context"
	"encoding/json"
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

	// GetMemberMetrics retrieves source control metrics for a specific member
	GetMemberMetrics(ctx context.Context, params *types.MemberMetricsParams) (*types.MetricsResponse, error)

	// Calculate time to merge metrics
	CalculateTimeToMerge(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation) (*int, error)
	CalculateTimeToMergeGraph(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation, metricLabel string, interval string) ([]types.TimeSeriesEntry, error)

	// Calculate PRs merged metrics
	CalculatePRsMerged(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation) (*int, error)
	CalculatePRsMergedGraph(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation, metricLabel string, interval string) ([]types.TimeSeriesEntry, error)

	// Calculate PRs reviewed metrics
	CalculatePRsReviewed(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation) (*int, error)
	CalculatePRsReviewedGraph(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation, metricLabel string, interval string) ([]types.TimeSeriesEntry, error)
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

// GetMemberMetrics retrieves source control metrics for a specific member
func (d *SourceControlDB) GetMemberMetrics(ctx context.Context, params *types.MemberMetricsParams) (*types.MetricsResponse, error) {
	// Build date filter conditions
	var dateFilter string
	var args []interface{}
	args = append(args, params.MemberID)

	if params.StartDate != nil {
		dateFilter += " AND pr.created_at >= ?"
		args = append(args, params.StartDate)
	}
	if params.EndDate != nil {
		dateFilter += " AND pr.created_at <= ?"
		args = append(args, params.EndDate)
	}

	// Get member's activities for metrics calculation
	activities, err := d.GetMemberActivity(ctx, &types.MemberActivityParams{
		MemberID:  params.MemberID,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
	})
	if err != nil {
		return nil, err
	}

	// Get organization ID for peer comparison
	var orgID string
	if err := d.db.WithContext(ctx).Raw(`
		SELECT om.organization_id 
		FROM organization_members om 
		WHERE om.id = ?
	`, params.MemberID).Scan(&orgID).Error; err != nil {
		return nil, err
	}

	// Get peer metrics for comparison
	peerMetrics, err := d.getPeerMetrics(ctx, orgID, params.MemberID, params.StartDate, params.EndDate)
	if err != nil {
		return nil, err
	}

	// Calculate snapshot metrics
	snapshotMetrics := d.calculateSnapshotMetrics(activities, peerMetrics)

	// Calculate graph metrics
	graphMetrics := d.calculateGraphMetrics(activities, peerMetrics, params.Interval)

	return &types.MetricsResponse{
		SnapshotMetrics: snapshotMetrics,
		GraphMetrics:    graphMetrics,
	}, nil
}

// getPeerMetrics calculates average metrics for other members in the same organization with the same title
func (d *SourceControlDB) getPeerMetrics(ctx context.Context, orgID, excludeMemberID string, startDate, endDate *time.Time) (map[string]float64, error) {
	// First, get the member's title
	var memberTitleID *string
	if err := d.db.WithContext(ctx).Raw(`
		SELECT title_id 
		FROM organization_members om 
		WHERE om.id = ?
	`, excludeMemberID).Scan(&memberTitleID).Error; err != nil {
		return nil, err
	}

	// If member has no title, return empty metrics (no peers to compare against)
	if memberTitleID == nil {
		return make(map[string]float64), nil
	}

	peerMetrics := make(map[string]float64)

	// Query pull request metrics for peers with same title
	prMetrics, err := d.getPeerPullRequestMetrics(ctx, orgID, *memberTitleID, excludeMemberID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Query review metrics for peers with same title
	reviewMetrics, err := d.getPeerReviewMetrics(ctx, orgID, *memberTitleID, excludeMemberID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Query comment metrics for peers with same title
	commentMetrics, err := d.getPeerCommentMetrics(ctx, orgID, *memberTitleID, excludeMemberID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Merge all metrics
	for k, v := range prMetrics {
		peerMetrics[k] = v
	}
	for k, v := range reviewMetrics {
		peerMetrics[k] = v
	}
	for k, v := range commentMetrics {
		peerMetrics[k] = v
	}

	return peerMetrics, nil
}

// getPeerPullRequestMetrics gets pull request metrics for peers
func (d *SourceControlDB) getPeerPullRequestMetrics(ctx context.Context, orgID, titleID, excludeMemberID string, startDate, endDate *time.Time) (map[string]float64, error) {
	query := `
		SELECT 
			COUNT(*) as pr_count,
			COALESCE(SUM(CAST(pr.metadata->>'additions' AS BIGINT)), 0) as total_additions,
			COALESCE(SUM(CAST(pr.metadata->>'deletions' AS BIGINT)), 0) as total_deletions,
			COALESCE(AVG(CAST(pr.metrics->>'time_to_merge_seconds' AS BIGINT)), 0) as avg_merge_time,
			COALESCE(AVG(CAST(pr.metrics->>'time_to_first_non_bot_review_seconds' AS BIGINT)), 0) as avg_review_time
		FROM pull_requests pr
		JOIN source_control_accounts sca ON pr.source_control_account_id = sca.id
		JOIN organization_members om ON sca.member_id = om.id
		WHERE om.organization_id = ? AND om.title_id = ? AND om.id != ? AND pr.status = 'merged'
	`

	var args []interface{}
	args = append(args, orgID, titleID, excludeMemberID)

	if startDate != nil {
		query += " AND pr.created_at >= ?"
		args = append(args, startDate)
	}
	if endDate != nil {
		query += " AND pr.created_at <= ?"
		args = append(args, endDate)
	}

	var result struct {
		PRCount        int64   `json:"pr_count"`
		TotalAdditions int64   `json:"total_additions"`
		TotalDeletions int64   `json:"total_deletions"`
		AvgMergeTime   float64 `json:"avg_merge_time"`
		AvgReviewTime  float64 `json:"avg_review_time"`
	}

	if err := d.db.WithContext(ctx).Raw(query, args...).Scan(&result).Error; err != nil {
		return nil, err
	}

	metrics := make(map[string]float64)
	if result.PRCount > 0 {
		metrics["merged_prs"] = float64(result.PRCount)
		metrics["loc_added"] = float64(result.TotalAdditions)
		metrics["loc_deleted"] = float64(result.TotalDeletions)
		metrics["mean_time_to_merge"] = result.AvgMergeTime
		metrics["mean_time_to_first_review"] = result.AvgReviewTime
	}

	return metrics, nil
}

// getPeerReviewMetrics gets review metrics for peers
func (d *SourceControlDB) getPeerReviewMetrics(ctx context.Context, orgID, titleID, excludeMemberID string, startDate, endDate *time.Time) (map[string]float64, error) {
	query := `
		SELECT COUNT(*) as review_count
		FROM pr_comments pc
		JOIN source_control_accounts sca ON pc.source_control_account_id = sca.id
		JOIN organization_members om ON sca.member_id = om.id
		WHERE om.organization_id = ? AND om.title_id = ? AND om.id != ? AND pc.type = 'REVIEW'
	`

	var args []interface{}
	args = append(args, orgID, titleID, excludeMemberID)

	if startDate != nil {
		query += " AND pc.created_at >= ?"
		args = append(args, startDate)
	}
	if endDate != nil {
		query += " AND pc.created_at <= ?"
		args = append(args, endDate)
	}

	var result struct {
		ReviewCount int `json:"review_count"`
	}

	if err := d.db.WithContext(ctx).Raw(query, args...).Scan(&result).Error; err != nil {
		return nil, err
	}

	metrics := make(map[string]float64)
	metrics["prs_reviewed"] = float64(result.ReviewCount)

	return metrics, nil
}

// getPeerCommentMetrics gets comment metrics for peers
func (d *SourceControlDB) getPeerCommentMetrics(ctx context.Context, orgID, titleID, excludeMemberID string, startDate, endDate *time.Time) (map[string]float64, error) {
	query := `
		SELECT COUNT(*) as comment_count
		FROM pr_comments pc
		JOIN source_control_accounts sca ON pc.source_control_account_id = sca.id
		JOIN organization_members om ON sca.member_id = om.id
		WHERE om.organization_id = ? AND om.title_id = ? AND om.id != ? AND pc.type != 'REVIEW'
	`

	var args []interface{}
	args = append(args, orgID, titleID, excludeMemberID)

	if startDate != nil {
		query += " AND pc.created_at >= ?"
		args = append(args, startDate)
	}
	if endDate != nil {
		query += " AND pc.created_at <= ?"
		args = append(args, endDate)
	}

	var result struct {
		CommentCount int `json:"comment_count"`
	}

	if err := d.db.WithContext(ctx).Raw(query, args...).Scan(&result).Error; err != nil {
		return nil, err
	}

	metrics := make(map[string]float64)
	metrics["pr_comments"] = float64(result.CommentCount)

	return metrics, nil
}

// calculateSnapshotMetrics calculates the snapshot metrics for the member
func (d *SourceControlDB) calculateSnapshotMetrics(activities []*types.MemberActivity, peerMetrics map[string]float64) []*types.SnapshotCategory {
	// Calculate member metrics
	memberMetrics := make(map[string]float64)

	prCount := 0
	reviewCount := 0
	commentCount := 0
	totalAdditions := 0
	totalDeletions := 0
	totalMergeTime := 0
	totalReviewTime := 0

	for _, activity := range activities {
		switch activity.Type {
		case "pull_request":
			prCount++
			// Parse metadata for LoC
			if activity.Metadata != nil {
				var metadata map[string]interface{}
				if err := json.Unmarshal(activity.Metadata, &metadata); err == nil {
					if additions, ok := metadata["additions"].(float64); ok {
						totalAdditions += int(additions)
					}
					if deletions, ok := metadata["deletions"].(float64); ok {
						totalDeletions += int(deletions)
					}
				}
			}
			// Parse metrics for time
			if activity.PRMetrics != nil {
				var metrics map[string]interface{}
				if err := json.Unmarshal(activity.PRMetrics, &metrics); err == nil {
					if mergeTime, ok := metrics["time_to_merge_seconds"].(float64); ok {
						totalMergeTime += int(mergeTime)
					}
					if reviewTime, ok := metrics["time_to_first_non_bot_review_seconds"].(float64); ok {
						totalReviewTime += int(reviewTime)
					}
				}
			}
		case "pr_review":
			reviewCount++
		case "pr_comment":
			commentCount++
		}
	}

	// Calculate averages
	if prCount > 0 {
		memberMetrics["mean_time_to_merge"] = float64(totalMergeTime) / float64(prCount)
		memberMetrics["mean_time_to_first_review"] = float64(totalReviewTime) / float64(prCount)
	}

	// Build snapshot categories
	activityCategory := types.SnapshotCategory{
		Category: "Activity",
		Metrics: []types.SnapshotMetric{
			{
				Label:      "Merged PRs",
				Value:      float64(prCount),
				PeersValue: peerMetrics["merged_prs"],
				Unit:       "count",
			},
			{
				Label:      "PRs Reviewed",
				Value:      float64(reviewCount),
				PeersValue: peerMetrics["prs_reviewed"],
				Unit:       "count",
			},
			{
				Label:      "LoC Added",
				Value:      float64(totalAdditions),
				PeersValue: peerMetrics["loc_added"],
				Unit:       "count",
			},
			{
				Label:      "LoC Deleted",
				Value:      float64(totalDeletions),
				PeersValue: peerMetrics["loc_deleted"],
				Unit:       "count",
			},
		},
	}

	efficiencyCategory := types.SnapshotCategory{
		Category: "Efficiency",
		Metrics: []types.SnapshotMetric{
			{
				Label:      "Mean time to merge",
				Value:      memberMetrics["mean_time_to_merge"],
				PeersValue: peerMetrics["mean_time_to_merge"],
				Unit:       "seconds",
			},
			{
				Label:      "Mean time to first review",
				Value:      memberMetrics["mean_time_to_first_review"],
				PeersValue: peerMetrics["mean_time_to_first_review"],
				Unit:       "seconds",
			},
			{
				Label:      "Time waiting on reviews",
				Value:      0, // Placeholder
				PeersValue: 0, // Placeholder
				Unit:       "seconds",
			},
		},
	}

	collaborationCategory := types.SnapshotCategory{
		Category: "Collaboration",
		Metrics: []types.SnapshotMetric{
			{
				Label:      "Mean response time to PRs",
				Value:      0, // Placeholder
				PeersValue: 0, // Placeholder
				Unit:       "seconds",
			},
		},
	}

	return []*types.SnapshotCategory{&activityCategory, &efficiencyCategory, &collaborationCategory}
}

// calculateGraphMetrics calculates the graph metrics for time series visualization
func (d *SourceControlDB) calculateGraphMetrics(activities []*types.MemberActivity, peerMetrics map[string]float64, interval string) []*types.GraphCategory {
	// For now, return a simple structure - this can be enhanced with actual time series data
	activityGraph := types.GraphCategory{
		Category: "Activity",
		Metrics: []types.GraphMetric{
			{
				Label: "Merged PRs",
				TimeSeries: []types.TimeSeriesEntry{
					{
						Date: time.Now().Format("2006-01-02"),
						Data: []types.TimeSeriesDataPoint{
							{
								Key:   "Merged PRs",
								Value: float64(len(activities)),
							},
							{
								Key:   "Peers merged PRs",
								Value: peerMetrics["merged_prs"],
							},
						},
					},
				},
			},
		},
	}

	return []*types.GraphCategory{&activityGraph}
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
		selectStatement = "COUNT(*)"
	default:
		return nil, fmt.Errorf("invalid metric operation for PRs reviewed: %s", metricOperation)
	}

	query := `
		SELECT ` + selectStatement + ` as prs_reviewed_count
		FROM pr_comments pc
		JOIN source_control_accounts sca ON pc.source_control_account_id = sca.id
		WHERE sca.organization_id = ?
		AND pc.created_at >= ?
		AND pc.created_at <= ?
		AND pc.type = 'REVIEW'
	`

	var args []any
	args = append(args, organizationID, startDate, endDate)

	if len(sourceControlAccountIDs) > 0 {
		query += " AND pc.source_control_account_id IN ?"
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

// CalculatePRsReviewedGraph calculates the PRs reviewed metric for a graph
func (d *SourceControlDB) CalculatePRsReviewedGraph(ctx context.Context, organizationID string, sourceControlAccountIDs []string, startDate, endDate time.Time, metricOperation metrictypes.MetricOperation, metricLabel string, interval string) ([]types.TimeSeriesEntry, error) {
	selectStatement := ""
	switch metricOperation {
	case metrictypes.MetricOperationCount:
		selectStatement = "COUNT(*)"
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
			DATE_TRUNC('` + postgresInterval + `', pc.created_at) as date,
			` + selectStatement + ` as prs_reviewed_count
		FROM pr_comments pc
		JOIN source_control_accounts sca ON pc.source_control_account_id = sca.id
		WHERE sca.organization_id = ?
		AND pc.created_at >= ?
		AND pc.created_at <= ?
		AND pc.type = 'REVIEW'
	`

	var args []any
	args = append(args, organizationID, startDate, endDate)

	if len(sourceControlAccountIDs) > 0 {
		query += " AND pc.source_control_account_id IN ?"
		args = append(args, sourceControlAccountIDs)
	}

	// Add GROUP BY clause for the DATE_TRUNC grouping
	query += " GROUP BY DATE_TRUNC('" + postgresInterval + "', pc.created_at)"

	// Add ORDER BY to ensure consistent results
	query += " ORDER BY date"

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
