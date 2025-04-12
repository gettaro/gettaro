package handlers

import (
	"net/http"
	"time"

	"ems.dev/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MetricsHandler struct {
	db *gorm.DB
}

func NewMetricsHandler(db *gorm.DB) *MetricsHandler {
	return &MetricsHandler{db: db}
}

// GetTeamMetrics returns metrics for a specific team
func (h *MetricsHandler) GetTeamMetrics(c *gin.Context) {
	teamID := c.Param("id")
	var metrics []models.TeamMetric
	query := h.db.Where("team_id = ?", teamID)

	// Apply date filters
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("start_date >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("end_date <= ?", endDate)
	}

	if err := query.Find(&metrics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// GetProjectMetrics returns metrics for a specific project
func (h *MetricsHandler) GetProjectMetrics(c *gin.Context) {
	projectID := c.Param("id")
	var metrics []models.ProjectMetric
	query := h.db.Where("project_id = ?", projectID)

	// Apply date filters
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("start_date >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("end_date <= ?", endDate)
	}

	if err := query.Find(&metrics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// GetUserMetrics returns metrics for a specific user
func (h *MetricsHandler) GetUserMetrics(c *gin.Context) {
	userID := c.Param("id")
	var metrics []models.PerformanceMetric
	query := h.db.Where("user_id = ?", userID)

	// Apply date filters
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("start_date >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("end_date <= ?", endDate)
	}

	if err := query.Find(&metrics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// CalculateMetrics calculates and stores metrics for a given period
func (h *MetricsHandler) CalculateMetrics(c *gin.Context) {
	period := c.Query("period") // weekly, monthly, quarterly
	if period == "" {
		period = "weekly"
	}

	// Calculate start and end dates based on period
	now := time.Now()
	var startDate, endDate time.Time

	switch period {
	case "weekly":
		startDate = now.AddDate(0, 0, -7)
		endDate = now
	case "monthly":
		startDate = now.AddDate(0, -1, 0)
		endDate = now
	case "quarterly":
		startDate = now.AddDate(0, -3, 0)
		endDate = now
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid period"})
		return
	}

	// Calculate team metrics
	var teams []models.Team
	if err := h.db.Find(&teams).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, team := range teams {
		// Calculate team metrics
		var totalTasks, completedTasks int
		var velocity float64

		// Get all projects for the team
		var projects []models.Project
		if err := h.db.Where("team_id = ?", team.ID).Find(&projects).Error; err != nil {
			continue
		}

		for _, project := range projects {
			// Get tasks for the project
			var tasks []models.Task
			if err := h.db.Where("project_id = ?", project.ID).Find(&tasks).Error; err != nil {
				continue
			}

			totalTasks += len(tasks)
			for _, task := range tasks {
				if task.Status == "done" {
					completedTasks++
				}
			}
		}

		// Calculate velocity (tasks completed per week)
		weeks := endDate.Sub(startDate).Hours() / (24 * 7)
		if weeks > 0 {
			velocity = float64(completedTasks) / weeks
		}

		// Create team metric
		metric := models.TeamMetric{
			TeamID:         team.ID,
			Period:         period,
			StartDate:      startDate,
			EndDate:        endDate,
			TotalTasks:     totalTasks,
			CompletedTasks: completedTasks,
			Velocity:       velocity,
		}

		if err := h.db.Create(&metric).Error; err != nil {
			continue
		}
	}

	c.Status(http.StatusOK)
}
