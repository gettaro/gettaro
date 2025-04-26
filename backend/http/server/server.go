package server

import (
	"os"

	authapi "ems.dev/backend/services/auth/api"
	integrationapi "ems.dev/backend/services/integration/api"
	orgapi "ems.dev/backend/services/organization/api"
	teamapi "ems.dev/backend/services/team/api"
	userapi "ems.dev/backend/services/user/api"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	router         *gin.Engine
	db             *gorm.DB
	userApi        userapi.UserAPI
	orgApi         orgapi.OrganizationAPI
	teamApi        teamapi.TeamAPI
	authApi        authapi.AuthAPI
	integrationApi integrationapi.IntegrationAPI
}

func New(db *gorm.DB, userApi userapi.UserAPI, orgApi orgapi.OrganizationAPI, teamApi teamapi.TeamAPI, authApi authapi.AuthAPI, integrationApi integrationapi.IntegrationAPI) *Server {
	s := &Server{
		router:         gin.Default(),
		db:             db,
		userApi:        userApi,
		orgApi:         orgApi,
		teamApi:        teamApi,
		authApi:        authApi,
		integrationApi: integrationApi,
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s
}

func (s *Server) setupMiddleware() {
	s.router.Use(gin.Logger())
	s.router.Use(gin.Recovery())

	// CORS middleware
	s.router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
