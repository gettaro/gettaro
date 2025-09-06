package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"ems.dev/backend/http/types/sourcecontrol"
	"ems.dev/backend/http/utils"
	organizationapi "ems.dev/backend/services/organization/api"
	sourcecontrolapi "ems.dev/backend/services/sourcecontrol/api"
	servicetypes "ems.dev/backend/services/sourcecontrol/types"
	"github.com/gin-gonic/gin"
)

type SourceControlHandler struct {
	scApi  sourcecontrolapi.SourceControlAPI
	orgApi organizationapi.OrganizationAPI
}

func NewSourceControlHandler(scApi sourcecontrolapi.SourceControlAPI, orgApi organizationapi.OrganizationAPI) *SourceControlHandler {
	return &SourceControlHandler{
		scApi:  scApi,
		orgApi: orgApi,
	}
}

// ListOrganizationPullRequests handles retrieving pull requests for an organization
// Params:
// - c: The Gin context containing request and response
// Query Parameters:
// - userIds: Optional list of user IDs to filter pull requests by
// - repositoryName: Optional repository name to filter pull requests by
// - startDate: Optional start date in format "2006-01-02" to filter pull requests by
// - endDate: Optional end date in format "2006-01-02" to filter pull requests by
// - status: Optional status to filter pull requests by (one of: "open", "closed", "merged")
// Returns:
// - 200: Success response with list of pull requests in ListOrganizationPullRequestsResponse format
// - 400: Bad request if organization ID is missing or query parameters are invalid
// - 401: Unauthorized if user is not authenticated
// - 403: Forbidden if user does not have access to the organization
// - 500: Internal server error if service layer fails
// Side Effects:
// - Makes a database query to fetch pull requests
// - Performs organization ownership check
// Errors:
// - ErrMissingOrganizationID: When organization ID is missing from the request
// - ErrInvalidDateFormat: When startDate or endDate has invalid format
// - ErrInvalidStatus: When status parameter has invalid value
// - ErrDatabaseQuery: When database query fails
// - ErrUnauthorized: When user is not authenticated
// - ErrForbidden: When user does not have access to the organization
func (h *SourceControlHandler) ListOrganizationPullRequests(c *gin.Context) {
	orgID, err := utils.GetOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user has access to the organization
	if !utils.CheckOrganizationOwnership(c, h.orgApi, orgID) {
		return
	}

	var query sourcecontrol.ListOrganizationPullRequestsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse dates if provided
	var startDate, endDate *time.Time
	if query.StartDate != "" {
		parsed, err := time.Parse("2006-01-02", query.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid startDate format"})
			return
		}
		startDate = &parsed
	}
	if query.EndDate != "" {
		parsed, err := time.Parse("2006-01-02", query.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid endDate format"})
			return
		}
		endDate = &parsed
	}

	// Get pull requests from service
	prs, err := h.scApi.GetPullRequests(c.Request.Context(), &servicetypes.PullRequestParams{
		OrganizationID: &orgID,
		UserIDs:        query.UserIDs,
		RepositoryName: query.RepositoryName,
		StartDate:      startDate,
		EndDate:        endDate,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := sourcecontrol.ListOrganizationPullRequestsResponse{
		PullRequests: make([]sourcecontrol.PullRequest, len(prs)),
	}

	for i, pr := range prs {
		response.PullRequests[i] = sourcecontrol.PullRequest{
			ID:             pr.ID,
			Title:          pr.Title,
			Description:    pr.Description,
			URL:            pr.URL,
			Status:         pr.Status,
			CreatedAt:      pr.CreatedAt,
			MergedAt:       pr.MergedAt,
			Comments:       pr.Comments,
			ReviewComments: pr.ReviewComments,
			Additions:      pr.Additions,
			Deletions:      pr.Deletions,
			ChangedFiles:   pr.ChangedFiles,
		}
	}

	c.JSON(http.StatusOK, response)
}

// ListOrganizationSourceControlAccounts handles retrieving source control accounts for an organization
// Params:
// - c: The Gin context containing request and response
// Returns:
// - 200: Success response with list of source control accounts
// - 400: Bad request if organization ID is missing
// - 401: Unauthorized if user is not authenticated
// - 403: Forbidden if user does not have access to the organization
// - 500: Internal server error if service layer fails
// Side Effects:
// - Makes a database query to fetch source control accounts
// - Performs organization membership check
// Errors:
// - ErrMissingOrganizationID: When organization ID is missing from the request
// - ErrDatabaseQuery: When database query fails
// - ErrUnauthorized: When user is not authenticated
// - ErrForbidden: When user does not have access to the organization
func (h *SourceControlHandler) ListOrganizationSourceControlAccounts(c *gin.Context) {
	orgID, err := utils.GetOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user has access to the organization
	if !utils.CheckOrganizationMembership(c, h.orgApi, &orgID) {
		return
	}

	// Get source control accounts from service
	accounts, err := h.scApi.GetSourceControlAccounts(c.Request.Context(), &servicetypes.SourceControlAccountParams{
		OrganizationID: orgID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := sourcecontrol.ListOrganizationSourceControlAccountsResponse{
		SourceControlAccounts: make([]servicetypes.SourceControlAccount, len(accounts)),
	}

	// Copy accounts to response
	copy(response.SourceControlAccounts, accounts)

	c.JSON(http.StatusOK, response)
}

// GetMemberPullRequests handles retrieving pull requests for a specific member
// Params:
// - c: The Gin context containing request and response
// Path Parameters:
// - id: Organization ID
// - memberId: Member ID to get pull requests for
// Query Parameters:
// - startDate: Optional start date in format "2006-01-02" to filter pull requests by
// - endDate: Optional end date in format "2006-01-02" to filter pull requests by
// Returns:
// - 200: Success response with list of pull requests
// - 400: Bad request if organization ID or member ID is missing
// - 401: Unauthorized if user is not authenticated
// - 403: Forbidden if user does not have access to the organization
// - 500: Internal server error if service layer fails
// Side Effects:
// - Makes a database query to fetch member pull requests
// - Performs organization membership check
// Errors:
// - ErrMissingOrganizationID: When organization ID is missing from the request
// - ErrInvalidDateFormat: When startDate or endDate has invalid format
// - ErrDatabaseQuery: When database query fails
// - ErrUnauthorized: When user is not authenticated
// - ErrForbidden: When user does not have access to the organization
func (h *SourceControlHandler) GetMemberPullRequests(c *gin.Context) {
	orgID, err := utils.GetOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user has access to the organization
	if !utils.CheckOrganizationMembership(c, h.orgApi, &orgID) {
		return
	}

	memberID := c.Param("memberId")
	if memberID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "member ID is required"})
		return
	}

	var query sourcecontrol.GetMemberPullRequestsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse dates if provided
	var startDate, endDate *time.Time
	if query.StartDate != "" {
		parsed, err := time.Parse("2006-01-02", query.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid startDate format, expected YYYY-MM-DD"})
			return
		}
		startDate = &parsed
	}
	if query.EndDate != "" {
		parsed, err := time.Parse("2006-01-02", query.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid endDate format, expected YYYY-MM-DD"})
			return
		}
		endDate = &parsed
	}

	// Get member pull requests from service
	prs, err := h.scApi.GetMemberPullRequests(c.Request.Context(), &servicetypes.MemberPullRequestParams{
		MemberID:  memberID,
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert service types to HTTP types
	response := sourcecontrol.GetMemberPullRequestsResponse{
		PullRequests: make([]sourcecontrol.PullRequest, len(prs)),
	}

	for i, pr := range prs {
		response.PullRequests[i] = sourcecontrol.PullRequest{
			ID:             pr.ID,
			Title:          pr.Title,
			Description:    pr.Description,
			URL:            pr.URL,
			Status:         pr.Status,
			CreatedAt:      pr.CreatedAt,
			MergedAt:       pr.MergedAt,
			Comments:       pr.Comments,
			ReviewComments: pr.ReviewComments,
			Additions:      pr.Additions,
			Deletions:      pr.Deletions,
			ChangedFiles:   pr.ChangedFiles,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetMemberPullRequestReviews handles retrieving pull request reviews for a specific member
// Params:
// - c: The Gin context containing request and response
// Path Parameters:
// - id: Organization ID
// - memberId: Member ID to get pull request reviews for
// Query Parameters:
// - startDate: Optional start date in format "2006-01-02" to filter reviews by
// - endDate: Optional end date in format "2006-01-02" to filter reviews by
// Returns:
// - 200: Success response with list of pull request reviews
// - 400: Bad request if organization ID or member ID is missing
// - 401: Unauthorized if user is not authenticated
// - 403: Forbidden if user does not have access to the organization
// - 500: Internal server error if service layer fails
// Side Effects:
// - Makes a database query to fetch member pull request reviews
// - Performs organization membership check
// Errors:
// - ErrMissingOrganizationID: When organization ID is missing from the request
// - ErrInvalidDateFormat: When startDate or endDate has invalid format
// - ErrDatabaseQuery: When database query fails
// - ErrUnauthorized: When user is not authenticated
// - ErrForbidden: When user does not have access to the organization
func (h *SourceControlHandler) GetMemberPullRequestReviews(c *gin.Context) {
	orgID, err := utils.GetOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user has access to the organization
	if !utils.CheckOrganizationMembership(c, h.orgApi, &orgID) {
		return
	}

	memberID := c.Param("memberId")
	if memberID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "member ID is required"})
		return
	}

	var query sourcecontrol.GetMemberPullRequestReviewsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse dates if provided
	var startDate, endDate *time.Time
	if query.StartDate != "" {
		parsed, err := time.Parse("2006-01-02", query.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid startDate format, expected YYYY-MM-DD"})
			return
		}
		startDate = &parsed
	}
	if query.EndDate != "" {
		parsed, err := time.Parse("2006-01-02", query.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid endDate format, expected YYYY-MM-DD"})
			return
		}
		endDate = &parsed
	}

	// Get member pull request reviews from service
	reviews, err := h.scApi.GetMemberPullRequestReviews(c.Request.Context(), &servicetypes.MemberPullRequestReviewsParams{
		MemberID:  memberID,
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert service types to HTTP types
	response := sourcecontrol.GetMemberPullRequestReviewsResponse{
		Reviews: make([]sourcecontrol.MemberActivity, len(reviews)),
	}

	for i, review := range reviews {
		// Convert metadata from datatypes.JSON to map[string]interface{}
		var metadata map[string]interface{}
		if review.Metadata != nil {
			if err := json.Unmarshal(review.Metadata, &metadata); err != nil {
				metadata = make(map[string]interface{})
			}
		}

		response.Reviews[i] = sourcecontrol.MemberActivity{
			ID:               review.ID,
			Type:             review.Type,
			Title:            review.Title,
			Description:      review.Description,
			URL:              review.URL,
			Repository:       review.Repository,
			CreatedAt:        review.CreatedAt,
			Metadata:         metadata,
			AuthorUsername:   review.AuthorUsername,
			PRTitle:          review.PRTitle,
			PRAuthorUsername: review.PRAuthorUsername,
			PRMetrics:        review.PRMetrics,
		}
	}

	c.JSON(http.StatusOK, response)
}

// RegisterRoutes registers all source control-related routes
func (h *SourceControlHandler) RegisterRoutes(api *gin.RouterGroup) {
	sourceControl := api.Group("/organizations/:id")
	{
		sourceControl.GET("/source-control-accounts", h.ListOrganizationSourceControlAccounts)
		sourceControl.GET("/pull-requests", h.ListOrganizationPullRequests)
	}

	// Member-specific routes
	members := api.Group("/organizations/:id/members")
	{
		members.GET("/:memberId/pull-requests", h.GetMemberPullRequests)
		members.GET("/:memberId/pull-request-reviews", h.GetMemberPullRequestReviews)
	}
}
