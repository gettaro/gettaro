package middleware

import (
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	issuerURL := os.Getenv("AUTH0_ISSUER_URL")
	audience := os.Getenv("AUTH0_CLIENT_ID")

	issuer, err := url.Parse(issuerURL)
	if err != nil {
		panic(err)
	}

	provider := jwks.NewCachingProvider(issuer, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL,
		[]string{audience},
	)
	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Remove "Bearer " prefix if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		claims, err := jwtValidator.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
