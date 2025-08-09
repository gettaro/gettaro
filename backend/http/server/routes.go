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

		// Integration routes
		integrationHandler := handlers.NewIntegrationHandler(s.integrationApi, s.orgApi)
		integrationHandler.RegisterRoutes(protected)

		// Source control routes
		sourceControlHandler := handlers.NewSourceControlHandler(s.sourcecontrolApi, s.orgApi)
		sourceControlHandler.RegisterRoutes(protected)

		// Team routes
		teamHandler := handlers.NewTeamHandler(s.teamApi, s.orgApi)
		teamHandler.RegisterRoutes(protected)

		// Title routes
		titleHandler := handlers.NewTitleHandler(s.titleApi, s.orgApi)
		titleHandler.RegisterRoutes(protected)

		// Member routes
		memberHandler := handlers.NewMemberHandler(s.memberApi, s.orgApi)
		memberHandler.RegisterRoutes(protected)

		// Direct reports routes
		directs := protected.Group("/directs")
		{
			directs.POST("", handlers.CreateDirectReport)
		}

		// Project management routes
		pmHandler := handlers.NewProjectManagementHandler()
		pmHandler.RegisterRoutes(protected)
	}
}
