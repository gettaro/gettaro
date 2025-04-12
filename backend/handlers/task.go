package handlers

import (
	"net/http"

	"ems.dev/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskHandler struct {
	db *gorm.DB
}

func NewTaskHandler(db *gorm.DB) *TaskHandler {
	return &TaskHandler{db: db}
}

// GetTasks returns all tasks
func (h *TaskHandler) GetTasks(c *gin.Context) {
	var tasks []models.Task
	if err := h.db.Preload("Project").Preload("Assignee").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// GetTask returns a single task by ID
func (h *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	if err := h.db.Preload("Project").Preload("Assignee").First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// CreateTask creates a new task
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate project exists
	var project models.Project
	if err := h.db.First(&project, task.ProjectID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project not found"})
		return
	}

	// Validate assignee exists if specified
	if task.AssigneeID != 0 {
		var user models.User
		if err := h.db.First(&user, task.AssigneeID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Assignee not found"})
			return
		}
	}

	if err := h.db.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// UpdateTask updates an existing task
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	if err := h.db.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate project exists if being updated
	if task.ProjectID != 0 {
		var project models.Project
		if err := h.db.First(&project, task.ProjectID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Project not found"})
			return
		}
	}

	// Validate assignee exists if being updated
	if task.AssigneeID != 0 {
		var user models.User
		if err := h.db.First(&user, task.AssigneeID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Assignee not found"})
			return
		}
	}

	if err := h.db.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteTask deletes a task
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&models.Task{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// GetTaskWorkLogs returns all work logs for a task
func (h *TaskHandler) GetTaskWorkLogs(c *gin.Context) {
	id := c.Param("id")
	var workLogs []models.WorkLog
	if err := h.db.Where("task_id = ?", id).Preload("User").Find(&workLogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, workLogs)
}
