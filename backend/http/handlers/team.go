package handlers

import (
	"net/http"

	teamtypes "ems.dev/backend/http/types/team"
	teamapi "ems.dev/backend/services/team/api"
	"ems.dev/backend/services/team/types"
	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	teamApi teamapi.TeamAPI
}

func NewTeamHandler(teamApi teamapi.TeamAPI) *TeamHandler {
	return &TeamHandler{
		teamApi: teamApi,
	}
}

func (h *TeamHandler) RegisterRoutes(router *gin.RouterGroup) {
	teams := router.Group("/teams")
	{
		teams.POST("", h.CreateTeam)
		teams.GET("", h.ListTeams)
		teams.GET("/:id", h.GetTeam)
		teams.PUT("/:id", h.UpdateTeam)
		teams.DELETE("/:id", h.DeleteTeam)
		teams.POST("/:id/members", h.AddTeamMember)
		teams.DELETE("/:id/members/:userId", h.RemoveTeamMember)
	}
}

func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var req types.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func (h *TeamHandler) ListTeams(c *gin.Context) {
	params := types.TeamSearchParams{}

	if orgID := c.Query("organizationId"); orgID != "" {
		params.OrganizationID = &orgID
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

func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "team ID is required"})
		return
	}

	var req types.UpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team := &types.Team{
		Name:        *req.Name,
		Description: *req.Description,
	}

	if err := h.teamApi.UpdateTeam(c.Request.Context(), id, team); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedTeam, err := h.teamApi.GetTeam(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teamtypes.GetTeamResponse(updatedTeam))
}

func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "team ID is required"})
		return
	}

	if err := h.teamApi.DeleteTeam(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TeamHandler) AddTeamMember(c *gin.Context) {
	teamID := c.Param("id")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "team ID is required"})
		return
	}

	var req types.AddTeamMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member := &types.TeamMember{
		TeamID: teamID,
		UserID: req.UserID,
		Role:   req.Role,
	}

	if err := h.teamApi.AddTeamMember(c.Request.Context(), teamID, member); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *TeamHandler) RemoveTeamMember(c *gin.Context) {
	teamID := c.Param("id")
	userID := c.Param("userId")
	if teamID == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "team ID and user ID are required"})
		return
	}

	if err := h.teamApi.RemoveTeamMember(c.Request.Context(), teamID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
