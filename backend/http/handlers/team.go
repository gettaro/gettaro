package handlers

import (
	"net/http"

	teamtypes "ems.dev/backend/http/types/team"
	"ems.dev/backend/http/utils"
	orgapi "ems.dev/backend/services/organization/api"
	teamapi "ems.dev/backend/services/team/api"
	"ems.dev/backend/services/team/types"
	"github.com/gin-gonic/gin"
)

// TeamHandler handles all team-related HTTP requests.
// It provides endpoints for team management including CRUD operations and member management.
type TeamHandler struct {
	teamApi teamapi.TeamAPI
	orgApi  orgapi.OrganizationAPI
}

// NewTeamHandler creates a new instance of TeamHandler.
// It initializes the handler with the provided TeamAPI.
func NewTeamHandler(teamApi teamapi.TeamAPI, orgApi orgapi.OrganizationAPI) *TeamHandler {
	return &TeamHandler{
		teamApi: teamApi,
		orgApi:  orgApi,
	}
}

// RegisterRoutes registers all team-related routes with the provided router group.
// It sets up the following routes:
// - POST /api/teams - Create a new team
// - GET /api/teams - List all teams
// - GET /api/teams/:id - Get a specific team
// - PUT /api/teams/:id - Update a team
// - DELETE /api/teams/:id - Delete a team
// - POST /api/teams/:id/members - Add a team member
// - DELETE /api/teams/:id/members/:memberId - Remove a team member
func (h *TeamHandler) RegisterRoutes(router *gin.RouterGroup) {
	teams := router.Group("/teams")
	{
		teams.POST("", h.CreateTeam)
		teams.GET("", h.ListTeams)
		teams.GET("/:id", h.GetTeam)
		teams.PUT("/:id", h.UpdateTeam)
		teams.DELETE("/:id", h.DeleteTeam)
		teams.POST("/:id/members", h.AddTeamMember)
		teams.DELETE("/:id/members/:memberId", h.RemoveTeamMember)
	}
}

// CreateTeam handles the POST /api/teams endpoint.
// It creates a new team with the provided information.
// Returns:
// - 201: The created team
// - 400: If the request body is invalid
// - 500: If there's a database error
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var req types.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user is an owner of the organization
	if !utils.CheckOrganizationOwnership(c, h.orgApi, req.OrganizationID) {
		return
	}

	team := &types.Team{
		Name:           req.Name,
		Description:    req.Description,
		OrganizationID: req.OrganizationID,
	}

	if err := h.teamApi.CreateTeam(c.Request.Context(), team); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, teamtypes.GetTeamResponse(team))
}

// GetTeam handles the GET /api/teams/:id endpoint.
// It retrieves a specific team by its ID.
// Returns:
// - 200: The team details
// - 400: If the team ID is missing
// - 404: If the team is not found
// - 500: If there's a database error
func (h *TeamHandler) GetTeam(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "team ID is required"})
		return
	}

	team, err := h.teamApi.GetTeam(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if team == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "team not found"})
		return
	}

	c.JSON(http.StatusOK, teamtypes.GetTeamResponse(team))
}

// ListTeams handles the GET /api/teams endpoint.
// It returns a list of teams, optionally filtered by organization ID or name.
// Returns:
// - 200: List of teams
// - 500: If there's a database error
func (h *TeamHandler) ListTeams(c *gin.Context) {
	params := types.TeamSearchParams{}

	if orgID := c.Query("organizationId"); orgID != "" {
		params.OrganizationID = &orgID
	}

	if name := c.Query("name"); name != "" {
		params.Name = &name
	}

	// Check if user is a member of the organization
	if !utils.CheckOrganizationMembership(c, h.orgApi, params.OrganizationID) {
		return
	}

	teams, err := h.teamApi.ListTeams(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := teamtypes.ListTeamsResponse{
		Teams: make([]teamtypes.TeamResponse, len(teams)),
	}
	for i, t := range teams {
		response.Teams[i] = teamtypes.GetTeamResponse(&t)
	}

	c.JSON(http.StatusOK, response)
}

// UpdateTeam handles the PUT /api/teams/:id endpoint.
// It updates an existing team's information.
// Returns:
// - 200: The updated team
// - 400: If the request body is invalid or team ID is missing
// - 404: If the team is not found
// - 500: If there's a database error
func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "team ID is required"})
		return
	}

	// Check if user is an owner of the organization
	team, err := h.teamApi.GetTeam(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !utils.CheckOrganizationOwnership(c, h.orgApi, team.OrganizationID) {
		return
	}

	var req types.UpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teamParams := &types.Team{
		Name:        *req.Name,
		Description: *req.Description,
	}

	if err := h.teamApi.UpdateTeam(c.Request.Context(), id, teamParams); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teamtypes.GetTeamResponse(teamParams))
}

// DeleteTeam handles the DELETE /api/teams/:id endpoint.
// It deletes a team from the system.
// Returns:
// - 204: If the team was successfully deleted
// - 400: If the team ID is missing
// - 500: If there's a database error
func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "team ID is required"})
		return
	}

	// Check if user is an owner of the organization
	team, err := h.teamApi.GetTeam(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !utils.CheckOrganizationOwnership(c, h.orgApi, team.OrganizationID) {
		return
	}

	if err := h.teamApi.DeleteTeam(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// AddTeamMember handles the POST /api/teams/:id/members endpoint.
// It adds a user as a member to a team.
// Returns:
// - 201: If the member was added successfully
// - 400: If the request body is invalid or team ID is missing
// - 500: If there's a database error
func (h *TeamHandler) AddTeamMember(c *gin.Context) {
	teamID := c.Param("id")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "team ID is required"})
		return
	}

	// Check if user is an owner of the organization
	team, err := h.teamApi.GetTeam(c.Request.Context(), teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !utils.CheckOrganizationOwnership(c, h.orgApi, team.OrganizationID) {
		return
	}

	var req types.AddTeamMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member := &types.TeamMember{
		TeamID:   teamID,
		MemberID: req.MemberID,
		Role:     req.Role,
	}

	if err := h.teamApi.AddTeamMember(c.Request.Context(), teamID, member); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// RemoveTeamMember handles the DELETE /api/teams/:id/members/:memberId endpoint.
// It removes a member from a team's members.
// Returns:
// - 204: If the member was successfully removed
// - 400: If the team ID or member ID is missing
// - 500: If there's a database error
func (h *TeamHandler) RemoveTeamMember(c *gin.Context) {
	teamID := c.Param("id")
	memberID := c.Param("memberId")
	if teamID == "" || memberID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "team ID and member ID are required"})
		return
	}

	// Check if user is an owner of the organization
	team, err := h.teamApi.GetTeam(c.Request.Context(), teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !utils.CheckOrganizationOwnership(c, h.orgApi, team.OrganizationID) {
		return
	}

	if err := h.teamApi.RemoveTeamMember(c.Request.Context(), teamID, memberID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
