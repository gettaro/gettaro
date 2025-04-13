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
	protected.Use(middleware.AuthMiddleware(s.userApi))
	{
		// User routes
		users := protected.Group("/users")
		{
			users.GET("", handlers.ListUsers)
			users.POST("", handlers.CreateUser)
			users.GET("/:id", handlers.GetUser)
			users.PUT("/:id", handlers.UpdateUser)
			users.DELETE("/:id", handlers.DeleteUser)
			users.GET("/me", handlers.NewUserHandler(s.userApi).GetMe)
		}

		// Organization routes
		orgs := protected.Group("/organizations")
		{
			orgs.POST("", handlers.CreateOrganization)
			orgs.GET("", handlers.ListOrganizations)
		}

		// Integration config routes
		integrations := protected.Group("/integration-configs")
		{
			integrations.POST("", handlers.CreateIntegrationConfig)
		}

		// Team routes
		teams := protected.Group("/teams")
		{
			teams.GET("", handlers.ListTeams)
			teams.POST("", handlers.CreateTeam)
			teams.POST("/:id/members", handlers.AddTeamMember)
		}

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
		pm := protected.Group("/project-management-accounts")
		{
			pm.GET("", handlers.ListPMAccounts)
			pm.POST("", handlers.CreatePMAccount)
		}

		// PM tickets routes
		tickets := protected.Group("/pm-tickets")
		{
			tickets.GET("", handlers.ListPMTickets)
			tickets.POST("", handlers.CreatePMTicket)
		}
	}
}
