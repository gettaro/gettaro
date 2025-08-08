package handlers

import (
	"fmt"
	"net/http"
	"time"

	"ems.dev/backend/http/types/sourcecontrol"
	types "ems.dev/backend/http/types/sourcecontrol"
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

// getOrganizationIDFromContext extracts the organization ID from the request context and returns it
func (h *SourceControlHandler) getOrganizationIDFromContext(c *gin.Context) (string, error) {
	id := c.Param("id")
	if id == "" {
		return "", fmt.Errorf("organization ID is required")
	}
	return id, nil
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
	orgID, err := h.getOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user has access to the organization
	if !utils.CheckOrganizationOwnership(c, h.orgApi, orgID) {
		return
	}

	var query types.ListOrganizationPullRequestsQuery
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

	response := types.ListOrganizationPullRequestsResponse{
		PullRequests: make([]types.PullRequest, len(prs)),
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

// ListOrganizationPullRequestsMetrics handles retrieving pull request metrics for an organization
// Params:
// - c: The Gin context containing request and response
// Query Parameters:
// - userIds: Optional list of user IDs to filter pull requests by
// - repositoryName: Optional repository name to filter pull requests by
// - startDate: Optional start date in format "2006-01-02" to filter pull requests by
// - endDate: Optional end date in format "2006-01-02" to filter pull requests by
// Returns:
// - 200: Success response with pull request metrics in PullRequestMetrics format
// - 400: Bad request if organization ID is missing or query parameters are invalid
// - 401: Unauthorized if user is not authenticated
// - 403: Forbidden if user does not have access to the organization
// - 500: Internal server error if service layer fails
// Side Effects:
// - Makes a database query to fetch pull requests
// - Performs organization ownership check
// - Calculates metrics from pull request data
// Errors:
// - ErrMissingOrganizationID: When organization ID is missing from the request
// - ErrInvalidDateFormat: When startDate or endDate has invalid format
// - ErrDatabaseQuery: When database query fails
// - ErrUnauthorized: When user is not authenticated
// - ErrForbidden: When user does not have access to the organization
// - ErrMetricsCalculation: When metrics calculation fails
func (h *SourceControlHandler) ListOrganizationPullRequestsMetrics(c *gin.Context) {
	orgID, err := h.getOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user has access to the organization
	if !utils.CheckOrganizationMembership(c, h.orgApi, &orgID) {
		return
	}

	var query types.ListOrganizationPullRequestsMetricsQuery
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

	// Get metrics from service
	metrics, err := h.scApi.GetPullRequestMetrics(c.Request.Context(), orgID, query.UserIDs, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// RegisterRoutes registers all source control-related routes
func (h *SourceControlHandler) RegisterRoutes(api *gin.RouterGroup) {
	sourceControl := api.Group("/organizations/:id")
	{
		sourceControl.GET("/pull-requests", h.ListOrganizationPullRequests)
		sourceControl.GET("/pull-requests/metrics", h.ListOrganizationPullRequestsMetrics)
	}
}
