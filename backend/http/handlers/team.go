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
// - POST /api/organizations/:id/teams - Create a new team
// - GET /api/organizations/:id/teams - List all teams for an organization
// - GET /api/organizations/:id/teams/:teamId - Get a specific team
// - PUT /api/organizations/:id/teams/:teamId - Update a team
// - DELETE /api/organizations/:id/teams/:teamId - Delete a team
// - POST /api/organizations/:id/teams/:teamId/members - Add a team member
// - DELETE /api/organizations/:id/teams/:teamId/members/:memberId - Remove a team member
func (h *TeamHandler) RegisterRoutes(router *gin.RouterGroup) {
	organizations := router.Group("/organizations")
	{
		organizations.POST("/:id/teams", h.CreateTeam)
		organizations.GET("/:id/teams", h.ListTeams)
		organizations.GET("/:id/teams/:teamId", h.GetTeam)
		organizations.PUT("/:id/teams/:teamId", h.UpdateTeam)
		organizations.DELETE("/:id/teams/:teamId", h.DeleteTeam)
		organizations.POST("/:id/teams/:teamId/members", h.AddTeamMember)
		organizations.DELETE("/:id/teams/:teamId/members/:memberId", h.RemoveTeamMember)
	}
}

// CreateTeam handles the POST /api/organizations/:id/teams endpoint.
// It creates a new team with the provided information.
// Returns:
// - 201: The created team
// - 400: If the request body is invalid
// - 500: If there's a database error
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	orgID := c.Param("id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "organization ID is required"})
		return
	}

	var req types.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user is an owner of the organization
	if !utils.CheckOrganizationOwnership(c, h.orgApi, orgID) {
		return
	}

	team := &types.Team{
		Name:           req.Name,
		Description:    req.Description,
		Type:           req.Type,
		PRPrefix:       req.PRPrefix,
		OrganizationID: orgID,
	}

	if err := h.teamApi.CreateTeam(c.Request.Context(), team); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, teamtypes.GetTeamResponseWrapped(team))
}

// GetTeam handles the GET /api/organizations/:id/teams/:teamId endpoint.
// It retrieves a specific team by its ID.
// Returns:
// - 200: The team details
// - 400: If the team ID is missing
// - 404: If the team is not found
// - 500: If there's a database error
func (h *TeamHandler) GetTeam(c *gin.Context) {
	orgID := c.Param("id")
	teamID := c.Param("teamId")
	if orgID == "" || teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "organization ID and team ID are required"})
		return
	}

	// Check if user is a member of the organization
	if !utils.CheckOrganizationMembership(c, h.orgApi, &orgID) {
		return
	}

	team, err := h.teamApi.GetTeamByOrganization(c.Request.Context(), teamID, orgID)
	if err != nil {
		if err.Error() == "team not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "team not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teamtypes.GetTeamResponseWrapped(team))
}

// ListTeams handles the GET /api/organizations/:id/teams endpoint.
// It returns a list of teams for the specified organization, optionally filtered by name.
// Returns:
// - 200: List of teams
// - 500: If there's a database error
func (h *TeamHandler) ListTeams(c *gin.Context) {
	orgID := c.Param("id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "organization ID is required"})
		return
	}

	// Check if user is a member of the organization
	if !utils.CheckOrganizationMembership(c, h.orgApi, &orgID) {
		return
	}

	params := types.TeamSearchParams{
		OrganizationID: &orgID,
	}

	if name := c.Query("name"); name != "" {
		params.Name = &name
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

// UpdateTeam handles the PUT /api/organizations/:id/teams/:teamId endpoint.
// It updates an existing team's information.
// Returns:
// - 200: The updated team
// - 400: If the request body is invalid or team ID is missing
// - 404: If the team is not found
// - 500: If there's a database error
func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	orgID := c.Param("id")
	teamID := c.Param("teamId")
	if orgID == "" || teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "organization ID and team ID are required"})
		return
	}

	// Check if user is an owner of the organization
	if !utils.CheckOrganizationOwnership(c, h.orgApi, orgID) {
		return
	}

	var req types.UpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teamParams := &types.Team{}
	if req.Name != nil {
		teamParams.Name = *req.Name
	}
	if req.Description != nil {
		teamParams.Description = *req.Description
	}
	if req.Type != nil {
		teamParams.Type = req.Type
	}
	if req.PRPrefix != nil {
		teamParams.PRPrefix = req.PRPrefix
	}

	if err := h.teamApi.UpdateTeam(c.Request.Context(), teamID, teamParams); err != nil {
		if err.Error() == "team not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "team not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch the updated team to return complete data including members
	updatedTeam, err := h.teamApi.GetTeamByOrganization(c.Request.Context(), teamID, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teamtypes.GetTeamResponseWrapped(updatedTeam))
}

// DeleteTeam handles the DELETE /api/organizations/:id/teams/:teamId endpoint.
// It deletes a team from the system.
// Returns:
// - 204: If the team was successfully deleted
// - 400: If the team ID is missing
// - 500: If there's a database error
func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	orgID := c.Param("id")
	teamID := c.Param("teamId")
	if orgID == "" || teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "organization ID and team ID are required"})
		return
	}

	// Check if user is an owner of the organization
	if !utils.CheckOrganizationOwnership(c, h.orgApi, orgID) {
		return
	}

	if err := h.teamApi.DeleteTeam(c.Request.Context(), teamID); err != nil {
		if err.Error() == "team not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "team not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// AddTeamMember handles the POST /api/organizations/:id/teams/:teamId/members endpoint.
// It adds a user as a member to a team.
// Returns:
// - 201: If the member was added successfully
// - 400: If the request body is invalid or team ID is missing
// - 500: If there's a database error
func (h *TeamHandler) AddTeamMember(c *gin.Context) {
	orgID := c.Param("id")
	teamID := c.Param("teamId")
	if orgID == "" || teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "organization ID and team ID are required"})
		return
	}

	// Check if user is an owner of the organization
	if !utils.CheckOrganizationOwnership(c, h.orgApi, orgID) {
		return
	}

	var req types.AddTeamMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member := &types.TeamMember{
		MemberID: req.MemberID,
	}

	if err := h.teamApi.AddTeamMember(c.Request.Context(), teamID, member); err != nil {
		if err.Error() == "team not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "team not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// RemoveTeamMember handles the DELETE /api/organizations/:id/teams/:teamId/members/:memberId endpoint.
// It removes a member from a team's members.
// Returns:
// - 204: If the member was successfully removed
// - 400: If the team ID or member ID is missing
// - 500: If there's a database error
func (h *TeamHandler) RemoveTeamMember(c *gin.Context) {
	orgID := c.Param("id")
	teamID := c.Param("teamId")
	memberID := c.Param("memberId")
	if orgID == "" || teamID == "" || memberID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "organization ID, team ID and member ID are required"})
		return
	}

	// Check if user is an owner of the organization
	if !utils.CheckOrganizationOwnership(c, h.orgApi, orgID) {
		return
	}

	if err := h.teamApi.RemoveTeamMember(c.Request.Context(), teamID, memberID); err != nil {
		if err.Error() == "team not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "team not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
