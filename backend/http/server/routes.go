package server

import (
	"ems.dev/backend/http/handlers"
	"ems.dev/backend/http/middleware"

	"github.com/gin-gonic/gin"
)

func (s *Server) setupRoutes() {
	// Public routes
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Protected routes
	protected := s.router.Group("/api")
	protected.Use(middleware.AuthMiddleware(s.userApi, s.authApi))
	{
		// User routes
		usersHandler := handlers.NewUsersHandler(s.userApi)
		usersHandler.RegisterRoutes(protected)

		// Organization routes
		orgHandler := handlers.NewOrganizationHandler(s.orgApi, s.userApi)
		orgHandler.RegisterRoutes(protected)

		// Integration config routes
		integrations := protected.Group("/integration-configs")
		{
			integrations.POST("", handlers.CreateIntegrationConfig)
		}

		// Team routes
		teamHandler := handlers.NewTeamHandler(s.teamApi, s.orgApi)
		teamHandler.RegisterRoutes(protected)

		// Direct reports routes
		directs := protected.Group("/directs")
		{
			directs.POST("", handlers.CreateDirectReport)
		}

		// Source control routes
		sourceControl := protected.Group("/source-control-accounts")
		{
			sourceControl.GET("", handlers.ListSourceControlAccounts)
			sourceControl.POST("", handlers.CreateSourceControlAccount)
		}

		// Pull request routes
		pullRequests := protected.Group("/pull-requests")
		{
			pullRequests.GET("", handlers.ListPullRequests)
			pullRequests.POST("", handlers.CreatePullRequest)
			pullRequests.POST("/:id/comments", handlers.CreatePRComment)
			pullRequests.POST("/:id/reviewers", handlers.AddPRReviewer)
		}

		// Project management routes
		pmHandler := handlers.NewProjectManagementHandler()
		pmHandler.RegisterRoutes(protected)
	}
}
