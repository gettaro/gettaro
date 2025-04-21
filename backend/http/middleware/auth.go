package middleware

import (
	"net/http"
	"net/url"
	"os"
	"time"

	authhttptypes "ems.dev/backend/http/types/auth"
	authapi "ems.dev/backend/services/auth/api"
	authtypes "ems.dev/backend/services/auth/types"
	userapi "ems.dev/backend/services/user/api"
	usertypes "ems.dev/backend/services/user/types"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(userApi userapi.UserAPI, authApi authapi.AuthAPI) gin.HandlerFunc {
	issuerURL := os.Getenv("AUTH0_AUTHORITY")
	audience := os.Getenv("AUTH0_AUDIENCE")

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
			return &authhttptypes.CustomClaims{}
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

		externalAuth, err := authApi.GetExternalAuth(c.Request.Context(), customClaims.RegisteredClaims.Subject)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get external auth: " + err.Error()})
			return
		}

		user := &usertypes.User{}
		if externalAuth == nil {
			// Get user info from Auth0
			userInfo, err := authApi.GetUserInfo(c.Request.Context(), token)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get userinfo from auth: " + err.Error()})
				return
			}

			// Create user
			user, err = userApi.CreateUser(&usertypes.User{
				Email: userInfo.Email,
				Name:  userInfo.Name,
			})

			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
				return
			}

			// Create external auth
			err = authApi.CreateExternalAuth(c.Request.Context(), &authtypes.AuthProvider{
				UserID:     user.ID,
				Provider:   "auth0",
				ProviderID: customClaims.RegisteredClaims.Subject,
			})

			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create external auth: " + err.Error()})
				return
			}

		} else {
			user, err = userApi.FindUser(usertypes.UserSearchParams{
				ID: &externalAuth.UserID,
			})

			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user: " + err.Error()})
				return
			}

			if user == nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
				return
			}
		}
		c.Set("user", user)
		c.Next()
	}
}
