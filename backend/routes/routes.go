package routes

import (
	"ems.dev/backend/handlers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.RouterGroup, db *gorm.DB) {
	// Initialize handlers
	teamHandler := handlers.NewTeamHandler(db)
	projectHandler := handlers.NewProjectHandler(db)
	taskHandler := handlers.NewTaskHandler(db)
	workLogHandler := handlers.NewWorkLogHandler(db)
	metricsHandler := handlers.NewMetricsHandler(db)

	// Team routes
	teams := r.Group("/teams")
	{
		teams.GET("", teamHandler.GetTeams)
		teams.GET("/:id", teamHandler.GetTeam)
		teams.POST("", teamHandler.CreateTeam)
		teams.PUT("/:id", teamHandler.UpdateTeam)
		teams.DELETE("/:id", teamHandler.DeleteTeam)
		teams.POST("/:id/users/:userId", teamHandler.AddUserToTeam)
		teams.DELETE("/:id/users/:userId", teamHandler.RemoveUserFromTeam)
		teams.GET("/:id/metrics", metricsHandler.GetTeamMetrics)
	}

	// Project routes
	projects := r.Group("/projects")
	{
		projects.GET("", projectHandler.GetProjects)
		projects.GET("/:id", projectHandler.GetProject)
		projects.POST("", projectHandler.CreateProject)
		projects.PUT("/:id", projectHandler.UpdateProject)
		projects.DELETE("/:id", projectHandler.DeleteProject)
		projects.GET("/:id/metrics", projectHandler.GetProjectMetrics)
		projects.GET("/:id/tasks", projectHandler.GetProjectTasks)
		projects.GET("/:id/team", projectHandler.GetProjectTeam)
	}

	// Task routes
	tasks := r.Group("/tasks")
	{
		tasks.GET("", taskHandler.GetTasks)
		tasks.GET("/:id", taskHandler.GetTask)
		tasks.POST("", taskHandler.CreateTask)
		tasks.PUT("/:id", taskHandler.UpdateTask)
		tasks.DELETE("/:id", taskHandler.DeleteTask)
		tasks.GET("/:id/work-logs", taskHandler.GetTaskWorkLogs)
	}

	// Work log routes
	workLogs := r.Group("/work-logs")
	{
		workLogs.POST("", workLogHandler.CreateWorkLog)
		workLogs.GET("", workLogHandler.GetWorkLogs)
		workLogs.GET("/users/:id", workLogHandler.GetUserWorkLogs)
		workLogs.GET("/tasks/:id", workLogHandler.GetTaskWorkLogs)
	}

	// Metrics routes
	metrics := r.Group("/metrics")
	{
		metrics.GET("/teams/:id", metricsHandler.GetTeamMetrics)
		metrics.GET("/projects/:id", metricsHandler.GetProjectMetrics)
		metrics.GET("/users/:id", metricsHandler.GetUserMetrics)
		metrics.POST("/calculate", metricsHandler.CalculateMetrics)
	}
}
