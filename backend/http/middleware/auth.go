package middleware

import (
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	authTypes "ems.dev/backend/http/types/auth"
	userapi "ems.dev/backend/services/user/api"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(userApi *userapi.Api) gin.HandlerFunc {
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
		validator.WithCustomClaims(func() validator.CustomClaims {
			return &authTypes.CustomClaims{}
		}),
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

		// Extract user claims from the validated token
		customClaims, ok := claims.(*validator.ValidatedClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Get the custom claims
		cc, ok := customClaims.CustomClaims.(*authTypes.CustomClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid custom claims format"})
			return
		}

		// Create user claims from the token
		userClaims := &authTypes.UserClaims{
			Sub:      customClaims.RegisteredClaims.Subject,
			Email:    cc.Email,
			Name:     cc.Name,
			Provider: "auth0", // Auth0 is our provider
		}

		// Only perform user creation/retrieval for /me endpoint
		if strings.HasSuffix(c.Request.URL.Path, "/api/users/me") && userApi != nil {
			dbUser, err := userApi.GetOrCreateUserFromAuthProvider(
				userClaims.Provider,
				userClaims.Sub,
				userClaims.Email,
				userClaims.Name,
			)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get or create user: " + err.Error()})
				return
			}
			c.Set("user", dbUser)
		}

		// Store the claims in the context
		c.Set("user_claims", userClaims)
		c.Next()
	}
}
