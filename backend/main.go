package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"ems.dev/backend/database"
	"ems.dev/backend/routes"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	database.InitDB()

	// Initialize Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
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

	// Public routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Protected routes
	protected := r.Group("/api")
	protected.Use(authMiddleware())
	{
		routes.SetupRoutes(protected, database.DB)
	}

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func authMiddleware() gin.HandlerFunc {
	issuerURL := os.Getenv("AUTH0_ISSUER_URL")
	audience := os.Getenv("AUTH0_AUDIENCE")

	validateToken := func(ctx context.Context, token string) (interface{}, error) {
		validator, err := validator.New(
			func(ctx context.Context) (interface{}, error) {
				return []byte(os.Getenv("AUTH0_PUBLIC_KEY")), nil
			},
			validator.RS256,
			issuerURL,
			[]string{audience},
		)
		if err != nil {
			return nil, err
		}

		return validator.ValidateToken(ctx, token)
	}

	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		_, err := validateToken(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.Next()
	}
}
