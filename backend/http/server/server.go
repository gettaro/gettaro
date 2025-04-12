package server

import (
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	router *gin.Engine
	db     *gorm.DB
}

func New(db *gorm.DB) *Server {
	s := &Server{
		router: gin.Default(),
		db:     db,
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s
}

func (s *Server) setupMiddleware() {
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
