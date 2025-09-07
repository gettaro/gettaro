package handlers

import (
	"net/http"

	directshttptypes "ems.dev/backend/http/types/directs"
	"ems.dev/backend/http/utils"
	directsapi "ems.dev/backend/services/directs/api"
	"ems.dev/backend/services/directs/types"
	membertypes "ems.dev/backend/services/member/types"
	orgapi "ems.dev/backend/services/organization/api"
	"github.com/gin-gonic/gin"
)

// DirectsHandler handles all direct reports-related HTTP requests.
type DirectsHandler struct {
	directsApi directsapi.DirectReportsAPI
	orgApi     orgapi.OrganizationAPI
}

// NewDirectsHandler creates a new instance of DirectsHandler.
func NewDirectsHandler(directsApi directsapi.DirectReportsAPI, orgApi orgapi.OrganizationAPI) *DirectsHandler {
	return &DirectsHandler{
		directsApi: directsApi,
		orgApi:     orgApi,
	}
}

// RegisterRoutes registers all direct reports-related routes.
func (h *DirectsHandler) RegisterRoutes(router *gin.RouterGroup) {
	organizations := router.Group("/organizations")
	{
		organizations.POST("/:id/directs", h.CreateDirectReport)
		organizations.GET("/:id/directs", h.ListDirectReports)
		organizations.GET("/:id/directs/:directId", h.GetDirectReport)
		organizations.PUT("/:id/directs/:directId", h.UpdateDirectReport)
		organizations.DELETE("/:id/directs/:directId", h.DeleteDirectReport)

		// Manager operations
		organizations.GET("/:id/managers/:managerId/directs", h.GetManagerDirectReports)
		organizations.GET("/:id/managers/:managerId/directs/tree", h.GetManagerTree)
		organizations.POST("/:id/managers/:managerId/directs", h.AddDirectReport)

		// Employee operations
		organizations.GET("/:id/members/:memberId/manager", h.GetMemberManager)
		organizations.GET("/:id/members/:memberId/management-chain", h.GetMemberManagementChain)

		// Organizational structure
		organizations.GET("/:id/org-chart", h.GetOrgChart)
		organizations.GET("/:id/org-chart/flat", h.GetOrgChartFlat)
	}
}

// CreateDirectReport creates a new direct report relationship.
func (h *DirectsHandler) CreateDirectReport(c *gin.Context) {
	orgID := c.Param("id")

	if !utils.ValidateOrganizationAccess(c, h.orgApi, orgID) {
		return
	}

	var req directshttptypes.CreateDirectReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := types.CreateDirectReportParams{
		ManagerMemberID: req.ManagerID,
		ReportMemberID:  req.ReportID,
		OrganizationID:  orgID,
		Depth:           req.Depth,
	}

	directReport, err := h.directsApi.CreateDirectReport(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := directshttptypes.CreateDirectReportResponse{
		DirectReport: h.mapDirectReportToResponse(directReport),
	}

	c.JSON(http.StatusCreated, response)
}

// ListDirectReports returns a list of direct report relationships.
func (h *DirectsHandler) ListDirectReports(c *gin.Context) {
	orgID := c.Param("id")

	if !utils.ValidateOrganizationAccess(c, h.orgApi, orgID) {
		return
	}

	var query directshttptypes.ListDirectReportsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := types.DirectReportSearchParams{
		OrganizationID:  &orgID,
		ManagerMemberID: query.ManagerID,
		ReportMemberID:  query.ReportID,
		Depth:           query.Depth,
	}

	directReports, err := h.directsApi.ListDirectReports(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := directshttptypes.ListDirectReportsResponse{
		DirectReports: h.mapDirectReportsToResponse(directReports),
	}

	c.JSON(http.StatusOK, response)
}

// GetDirectReport retrieves a specific direct report relationship.
func (h *DirectsHandler) GetDirectReport(c *gin.Context) {
	orgID := c.Param("id")
	directReportID := c.Param("directId")

	if !utils.ValidateOrganizationAccess(c, h.orgApi, orgID) {
		return
	}

	directReport, err := h.directsApi.GetDirectReport(c.Request.Context(), directReportID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if directReport == nil || directReport.OrganizationID != orgID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Direct report not found"})
		return
	}

	response := directshttptypes.GetDirectReportResponse{
		DirectReport: h.mapDirectReportToResponse(directReport),
	}

	c.JSON(http.StatusOK, response)
}

// UpdateDirectReport updates an existing direct report relationship.
func (h *DirectsHandler) UpdateDirectReport(c *gin.Context) {
	orgID := c.Param("id")
	directReportID := c.Param("directId")

	if !utils.ValidateOrganizationAccess(c, h.orgApi, orgID) {
		return
	}

	var req directshttptypes.UpdateDirectReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := types.UpdateDirectReportParams{
		Depth: req.Depth,
	}

	err := h.directsApi.UpdateDirectReport(c.Request.Context(), directReportID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	directReport, err := h.directsApi.GetDirectReport(c.Request.Context(), directReportID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := directshttptypes.UpdateDirectReportResponse{
		DirectReport: h.mapDirectReportToResponse(directReport),
	}

	c.JSON(http.StatusOK, response)
}

// DeleteDirectReport removes a direct report relationship.
func (h *DirectsHandler) DeleteDirectReport(c *gin.Context) {
	orgID := c.Param("id")
	directReportID := c.Param("directId")

	if !utils.ValidateOrganizationAccess(c, h.orgApi, orgID) {
		return
	}

	directReport, err := h.directsApi.GetDirectReport(c.Request.Context(), directReportID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if directReport == nil || directReport.OrganizationID != orgID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Direct report not found"})
		return
	}

	err = h.directsApi.DeleteDirectReport(c.Request.Context(), directReportID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetManagerDirectReports returns all direct reports for a specific manager.
func (h *DirectsHandler) GetManagerDirectReports(c *gin.Context) {
	orgID := c.Param("id")
	managerID := c.Param("managerId")

	if !utils.ValidateOrganizationAccess(c, h.orgApi, orgID) {
		return
	}

	directReports, err := h.directsApi.GetManagerDirectReports(c.Request.Context(), managerID, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := directshttptypes.GetManagerDirectReportsResponse{
		DirectReports: h.mapDirectReportsToResponse(directReports),
	}

	c.JSON(http.StatusOK, response)
}

// GetManagerTree returns the full management tree for a manager.
func (h *DirectsHandler) GetManagerTree(c *gin.Context) {
	orgID := c.Param("id")
	managerID := c.Param("managerId")

	if !utils.ValidateOrganizationAccess(c, h.orgApi, orgID) {
		return
	}

	orgChart, err := h.directsApi.GetManagerTree(c.Request.Context(), managerID, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := directshttptypes.GetManagerTreeResponse{
		OrgChart: h.mapOrgChartNodesToResponse(orgChart),
	}

	c.JSON(http.StatusOK, response)
}

// AddDirectReport adds a direct report to a manager.
func (h *DirectsHandler) AddDirectReport(c *gin.Context) {
	orgID := c.Param("id")
	managerID := c.Param("managerId")

	if !utils.ValidateOrganizationAccess(c, h.orgApi, orgID) {
		return
	}

	var req struct {
		ReportID string `json:"reportId" binding:"required"`
		Depth    int    `json:"depth" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := types.CreateDirectReportParams{
		ManagerMemberID: managerID,
		ReportMemberID:  req.ReportID,
		OrganizationID:  orgID,
		Depth:           req.Depth,
	}

	directReport, err := h.directsApi.CreateDirectReport(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := directshttptypes.CreateDirectReportResponse{
		DirectReport: h.mapDirectReportToResponse(directReport),
	}

	c.JSON(http.StatusCreated, response)
}

// GetMemberManager returns the manager of a specific member.
func (h *DirectsHandler) GetMemberManager(c *gin.Context) {
	orgID := c.Param("id")
	memberID := c.Param("memberId")

	if !utils.ValidateOrganizationAccess(c, h.orgApi, orgID) {
		return
	}

	directReport, err := h.directsApi.GetMemberManager(c.Request.Context(), memberID, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response directshttptypes.GetMemberManagerResponse
	if directReport != nil {
		response.Manager = &directshttptypes.DirectReportResponse{
			ID:             directReport.ID,
			ManagerID:      directReport.ManagerMemberID,
			ReportID:       directReport.ReportMemberID,
			OrganizationID: directReport.OrganizationID,
			Depth:          directReport.Depth,
			CreatedAt:      directReport.CreatedAt,
			UpdatedAt:      directReport.UpdatedAt,
			Manager:        h.mapUserToResponse(&directReport.Manager),
			Report:         h.mapUserToResponse(&directReport.Report),
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetMemberManagementChain returns the full management chain for a member.
func (h *DirectsHandler) GetMemberManagementChain(c *gin.Context) {
	orgID := c.Param("id")
	memberID := c.Param("memberId")

	if !utils.ValidateOrganizationAccess(c, h.orgApi, orgID) {
		return
	}

	managementChain, err := h.directsApi.GetMemberManagementChain(c.Request.Context(), memberID, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := directshttptypes.GetMemberManagementChainResponse{
		ManagementChain: h.mapManagementChainToResponse(managementChain),
	}

	c.JSON(http.StatusOK, response)
}

// GetOrgChart returns the complete organizational chart.
func (h *DirectsHandler) GetOrgChart(c *gin.Context) {
	orgID := c.Param("id")

	if !utils.ValidateOrganizationAccess(c, h.orgApi, orgID) {
		return
	}

	orgChart, err := h.directsApi.GetOrgChart(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := directshttptypes.GetOrgChartResponse{
		OrgChart: h.mapOrgChartNodesToResponse(orgChart),
	}

	c.JSON(http.StatusOK, response)
}

// GetOrgChartFlat returns a flat list of all manager-direct relationships.
func (h *DirectsHandler) GetOrgChartFlat(c *gin.Context) {
	orgID := c.Param("id")

	if !utils.ValidateOrganizationAccess(c, h.orgApi, orgID) {
		return
	}

	directReports, err := h.directsApi.GetOrgChartFlat(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := directshttptypes.GetOrgChartFlatResponse{
		DirectReports: h.mapDirectReportsToResponse(directReports),
	}

	c.JSON(http.StatusOK, response)
}

// Helper methods for mapping between service types and HTTP response types

func (h *DirectsHandler) mapDirectReportToResponse(dr *types.DirectReport) directshttptypes.DirectReportResponse {
	response := directshttptypes.DirectReportResponse{
		ID:             dr.ID,
		ManagerID:      dr.ManagerMemberID,
		ReportID:       dr.ReportMemberID,
		OrganizationID: dr.OrganizationID,
		Depth:          dr.Depth,
		CreatedAt:      dr.CreatedAt,
		UpdatedAt:      dr.UpdatedAt,
	}

	if dr.Manager.ID != "" {
		response.Manager = h.mapUserToResponse(&dr.Manager)
	}
	if dr.Report.ID != "" {
		response.Report = h.mapUserToResponse(&dr.Report)
	}

	return response
}

func (h *DirectsHandler) mapDirectReportsToResponse(drs []types.DirectReport) []directshttptypes.DirectReportResponse {
	responses := make([]directshttptypes.DirectReportResponse, len(drs))
	for i, dr := range drs {
		responses[i] = h.mapDirectReportToResponse(&dr)
	}
	return responses
}

func (h *DirectsHandler) mapUserToResponse(member *membertypes.OrganizationMember) *directshttptypes.MemberResponse {
	if member == nil || member.ID == "" {
		return nil
	}

	return &directshttptypes.MemberResponse{
		ID:        member.ID,
		Email:     member.Email,
		Username:  member.Username, // Also include username field
		TitleID:   *member.TitleID,
		CreatedAt: member.CreatedAt,
		UpdatedAt: member.UpdatedAt,
	}
}

func (h *DirectsHandler) mapOrgChartNodesToResponse(nodes []types.OrgChartNode) []directshttptypes.OrgChartNodeResponse {
	responses := make([]directshttptypes.OrgChartNodeResponse, len(nodes))
	for i, node := range nodes {
		responses[i] = directshttptypes.OrgChartNodeResponse{
			Member:        *h.mapUserToResponse(&node.Member),
			DirectReports: h.mapOrgChartNodesToResponse(node.DirectReports),
			Depth:         node.Depth,
		}
	}
	return responses
}

func (h *DirectsHandler) mapManagementChainToResponse(chain []types.ManagementChain) []directshttptypes.ManagementChainResponse {
	responses := make([]directshttptypes.ManagementChainResponse, len(chain))
	for i, mc := range chain {
		responses[i] = directshttptypes.ManagementChainResponse{
			Member:  *h.mapUserToResponse(&mc.Member),
			Manager: h.mapUserToResponse(mc.Manager),
			Depth:   mc.Depth,
		}
	}
	return responses
}
