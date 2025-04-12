package handlers

import (
	"net/http"

	"ems.dev/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectHandler struct {
	db *gorm.DB
}

func NewProjectHandler(db *gorm.DB) *ProjectHandler {
	return &ProjectHandler{db: db}
}

// GetProjects returns all projects
func (h *ProjectHandler) GetProjects(c *gin.Context) {
	var projects []models.Project
	if err := h.db.Preload("Team").Preload("Team.Users").Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projects)
}

// GetProject returns a single project by ID
func (h *ProjectHandler) GetProject(c *gin.Context) {
	id := c.Param("id")
	var project models.Project
	if err := h.db.Preload("Team").Preload("Team.Users").First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}
	c.JSON(http.StatusOK, project)
}

// CreateProject creates a new project
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate team exists
	var team models.Team
	if err := h.db.First(&team, project.TeamID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team not found"})
		return
	}

	if err := h.db.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// UpdateProject updates an existing project
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id := c.Param("id")
	var project models.Project
	if err := h.db.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate team exists if being updated
	if project.TeamID != 0 {
		var team models.Team
		if err := h.db.First(&team, project.TeamID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Team not found"})
			return
		}
	}

	if err := h.db.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, project)
}

// DeleteProject deletes a project
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&models.Project{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// GetProjectMetrics returns metrics for a project
func (h *ProjectHandler) GetProjectMetrics(c *gin.Context) {
	id := c.Param("id")
	var metrics []models.ProjectMetric
	if err := h.db.Where("project_id = ?", id).Find(&metrics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, metrics)
}

// GetProjectTasks returns all tasks for a project
func (h *ProjectHandler) GetProjectTasks(c *gin.Context) {
	id := c.Param("id")
	var tasks []models.Task
	if err := h.db.Where("project_id = ?", id).Preload("Assignee").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// GetProjectTeam returns the team associated with a project
func (h *ProjectHandler) GetProjectTeam(c *gin.Context) {
	id := c.Param("id")
	var project models.Project
	if err := h.db.Preload("Team").Preload("Team.Users").First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}
	c.JSON(http.StatusOK, project.Team)
}
