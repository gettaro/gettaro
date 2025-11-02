package server

import (
	"os"

	aicodeassistantapi "ems.dev/backend/services/aicodeassistant/api"
	apiai "ems.dev/backend/services/ai/api"
	authapi "ems.dev/backend/services/auth/api"
	conversationapi "ems.dev/backend/services/conversation/api"
	conversationtemplateapi "ems.dev/backend/services/conversationtemplate/api"
	directsapi "ems.dev/backend/services/directs/api"
	integrationapi "ems.dev/backend/services/integration/api"
	memberapi "ems.dev/backend/services/member/api"
	metricsapi "ems.dev/backend/services/metrics/api"
	orgapi "ems.dev/backend/services/organization/api"
	sourcecontrolapi "ems.dev/backend/services/sourcecontrol/api"
	teamapi "ems.dev/backend/services/team/api"
	titleapi "ems.dev/backend/services/title/api"
	userapi "ems.dev/backend/services/user/api"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	router                  *gin.Engine
	db                      *gorm.DB
	userApi                 userapi.UserAPI
	orgApi                  orgapi.OrganizationAPI
	teamApi                 teamapi.TeamAPI
	titleApi                titleapi.TitleAPI
	authApi                 authapi.AuthAPI
	integrationApi          integrationapi.IntegrationAPI
	sourcecontrolApi        sourcecontrolapi.SourceControlAPI
	memberApi               memberapi.MemberAPI
	metricsApi              metricsapi.MetricsAPI
	directsApi              directsapi.DirectReportsAPI
	conversationTemplateApi conversationtemplateapi.ConversationTemplateAPIInterface
	conversationApi         conversationapi.ConversationAPIInterface
	aiApi                   apiai.AIServiceInterface
	aiCodeAssistantApi      aicodeassistantapi.AICodeAssistantAPI
}

func New(db *gorm.DB, userApi userapi.UserAPI, orgApi orgapi.OrganizationAPI, teamApi teamapi.TeamAPI, titleApi titleapi.TitleAPI, authApi authapi.AuthAPI, integrationApi integrationapi.IntegrationAPI, sourcecontrolApi sourcecontrolapi.SourceControlAPI, memberApi memberapi.MemberAPI, metricsApi metricsapi.MetricsAPI, directsApi directsapi.DirectReportsAPI, conversationTemplateApi conversationtemplateapi.ConversationTemplateAPIInterface, conversationApi conversationapi.ConversationAPIInterface, aiApi apiai.AIServiceInterface, aiCodeAssistantApi aicodeassistantapi.AICodeAssistantAPI) *Server {
	s := &Server{
		router:                  gin.Default(),
		db:                      db,
		userApi:                 userApi,
		orgApi:                  orgApi,
		teamApi:                 teamApi,
		titleApi:                titleApi,
		authApi:                 authApi,
		integrationApi:          integrationApi,
		sourcecontrolApi:        sourcecontrolApi,
		memberApi:               memberApi,
		metricsApi:              metricsApi,
		directsApi:              directsApi,
		conversationTemplateApi: conversationTemplateApi,
		conversationApi:         conversationApi,
		aiApi:                   aiApi,
		aiCodeAssistantApi:      aiCodeAssistantApi,
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
