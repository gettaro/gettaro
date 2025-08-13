package database

import (
	"context"

	"ems.dev/backend/services/sourcecontrol/types"
	"gorm.io/gorm"
)

// SourceControlDB defines the interface for source control database operations
type DB interface {
	// Source Control Accounts
	GetSourceControlAccountsByUsernames(ctx context.Context, usernames []string) (map[string]*types.SourceControlAccount, error)
	CreateSourceControlAccounts(ctx context.Context, accounts []*types.SourceControlAccount) error
	GetSourceControlAccount(ctx context.Context, id string) (*types.SourceControlAccount, error)
	UpdateSourceControlAccount(ctx context.Context, account *types.SourceControlAccount) error
	GetSourceControlAccountsByOrganization(ctx context.Context, orgID string) ([]*types.SourceControlAccount, error)

	// Pull Requests
	GetPullRequests(ctx context.Context, params *types.PullRequestParams) ([]*types.PullRequest, error)
	CreatePullRequests(ctx context.Context, prs []*types.PullRequest) error
	UpdatePullRequest(ctx context.Context, pr *types.PullRequest) error

	// Comments
	CreatePRComments(ctx context.Context, comments []*types.PRComment) error
	GetPullRequestComments(ctx context.Context, prID string) ([]*types.PRComment, error)

	// Member Activity
	GetMemberActivity(ctx context.Context, params *types.MemberActivityParams) ([]*types.MemberActivity, error)
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

// GetSourceControlAccountsByUsernames retrieves source control accounts by their usernames
func (d *SourceControlDB) GetSourceControlAccountsByUsernames(ctx context.Context, usernames []string) (map[string]*types.SourceControlAccount, error) {
	var accounts []types.SourceControlAccount
	err := d.db.WithContext(ctx).
		Where("username IN ?", usernames).
		Find(&accounts).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string]*types.SourceControlAccount)
	for i := range accounts {
		result[accounts[i].Username] = &accounts[i]
	}

	return result, nil
}

// CreateSourceControlAccounts creates multiple source control accounts
func (d *SourceControlDB) CreateSourceControlAccounts(ctx context.Context, accounts []*types.SourceControlAccount) error {
	return d.db.WithContext(ctx).Create(accounts).Error
}

// GetPullRequests retrieves pull requests based on the given parameters
func (d *SourceControlDB) GetPullRequests(ctx context.Context, params *types.PullRequestParams) ([]*types.PullRequest, error) {
	var prs []types.PullRequest
	query := d.db.WithContext(ctx).Model(&types.PullRequest{})
	query = query.Joins(`
		JOIN source_control_accounts sca ON pull_requests.source_control_account_id = sca.id
		JOIN organization_members om ON sca.member_id = om.id
	`)

	if params.ProviderID != "" {
		query = query.Where("provider_id = ?", params.ProviderID)
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

// CreatePullRequests creates multiple pull requests
func (d *SourceControlDB) CreatePullRequests(ctx context.Context, prs []*types.PullRequest) error {
	return d.db.WithContext(ctx).Create(prs).Error
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
	return d.db.WithContext(ctx).Model(account).Updates(account).Error
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

// GetSourceControlAccountsByOrganization retrieves source control accounts for an organization
func (d *SourceControlDB) GetSourceControlAccountsByOrganization(ctx context.Context, orgID string) ([]*types.SourceControlAccount, error) {
	var accounts []types.SourceControlAccount
	err := d.db.WithContext(ctx).
		Where("organization_id = ?", orgID).
		Find(&accounts).Error
	if err != nil {
		return nil, err
	}

	result := make([]*types.SourceControlAccount, len(accounts))
	for i := range accounts {
		result[i] = &accounts[i]
	}
	return result, nil
}

// GetMemberActivity retrieves a timeline of source control activities for a specific member
func (d *SourceControlDB) GetMemberActivity(ctx context.Context, params *types.MemberActivityParams) ([]*types.MemberActivity, error) {
	var activities []types.MemberActivity

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

	// Get pull requests with their activities for the member
	query := `
		SELECT DISTINCT
			pr.id,
			'pull_request' as type,
			pr.title,
			COALESCE(pr.description, '') as description,
			COALESCE(pr.url, '') as url,
			COALESCE(pr.repository_name, '') as repository,
			pr.created_at,
			COALESCE(pr.metadata, '{}'::jsonb) as metadata,
			sca.username as author_username
		FROM pull_requests pr
		JOIN source_control_accounts sca ON pr.source_control_account_id = sca.id
		WHERE EXISTS (
			SELECT 1 FROM source_control_accounts sca2
			WHERE sca2.id = pr.source_control_account_id 
			AND sca2.member_id = ?
		)` + dateFilter + `
		ORDER BY pr.created_at DESC
	`

	var prActivities []types.MemberActivity
	if err := d.db.WithContext(ctx).Raw(query, args...).Scan(&prActivities).Error; err != nil {
		return nil, err
	}

	// For each PR, get its comments and reviews
	for _, pr := range prActivities {
		// Get comments for this PR
		commentQuery := `
			SELECT 
				pc.id,
				CASE 
					WHEN pc.type = 'REVIEW' THEN 'pr_review'
					ELSE 'pr_comment'
				END as type,
				pr.title as title,
				pc.body as description,
				pr.url as url,
				pr.repository_name as repository,
				pc.created_at,
				'{}'::jsonb as metadata,
				sca.username as author_username
			FROM pr_comments pc
			JOIN pull_requests pr ON pc.pr_id = pr.id
			JOIN source_control_accounts sca ON pc.source_control_account_id = sca.id
			WHERE pc.pr_id = ? AND EXISTS (
				SELECT 1 FROM source_control_accounts sca2
				WHERE sca2.id = pc.source_control_account_id 
				AND sca2.member_id = ?
			)
			ORDER BY pc.created_at DESC
		`

		var comments []types.MemberActivity
		if err := d.db.WithContext(ctx).Raw(commentQuery, pr.ID, params.MemberID).Scan(&comments).Error; err != nil {
			return nil, err
		}

		// Add comments to the activities
		activities = append(activities, pr)
		activities = append(activities, comments...)
	}

	// Sort all activities by creation date (newest first)
	for i := 0; i < len(activities)-1; i++ {
		for j := i + 1; j < len(activities); j++ {
			if activities[i].CreatedAt.Before(activities[j].CreatedAt) {
				activities[i], activities[j] = activities[j], activities[i]
			}
		}
	}

	result := make([]*types.MemberActivity, len(activities))
	for i := range activities {
		result[i] = &activities[i]
	}
	return result, nil
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
