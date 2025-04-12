package handlers

import (
	"net/http"

	"ems.dev/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkLogHandler struct {
	db *gorm.DB
}

func NewWorkLogHandler(db *gorm.DB) *WorkLogHandler {
	return &WorkLogHandler{db: db}
}

// CreateWorkLog creates a new work log
func (h *WorkLogHandler) CreateWorkLog(c *gin.Context) {
	var workLog models.WorkLog
	if err := c.ShouldBindJSON(&workLog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate user exists
	var user models.User
	if err := h.db.First(&user, workLog.UserID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// Validate task exists
	var task models.Task
	if err := h.db.First(&task, workLog.TaskID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task not found"})
		return
	}

	// Calculate duration if not provided
	if workLog.Duration == 0 && !workLog.EndTime.IsZero() {
		workLog.Duration = int(workLog.EndTime.Sub(workLog.StartTime).Minutes())
	}

	if err := h.db.Create(&workLog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, workLog)
}

// GetWorkLogs returns all work logs with optional filters
func (h *WorkLogHandler) GetWorkLogs(c *gin.Context) {
	var workLogs []models.WorkLog
	query := h.db.Preload("User").Preload("Task")

	// Apply filters
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if taskID := c.Query("task_id"); taskID != "" {
		query = query.Where("task_id = ?", taskID)
	}
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("start_time >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("end_time <= ?", endDate)
	}

	if err := query.Find(&workLogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workLogs)
}

// GetUserWorkLogs returns all work logs for a specific user
func (h *WorkLogHandler) GetUserWorkLogs(c *gin.Context) {
	userID := c.Param("id")
	var workLogs []models.WorkLog
	if err := h.db.Where("user_id = ?", userID).Preload("Task").Find(&workLogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, workLogs)
}

// GetTaskWorkLogs returns all work logs for a specific task
func (h *WorkLogHandler) GetTaskWorkLogs(c *gin.Context) {
	taskID := c.Param("id")
	var workLogs []models.WorkLog
	if err := h.db.Where("task_id = ?", taskID).Preload("User").Find(&workLogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, workLogs)
}
