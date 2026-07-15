package middleware

import (
	"net/http"

	"github.com/chai-rs/simple-bookstore/infrastructure/auth"
	"github.com/chai-rs/simple-bookstore/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// AuthMiddleware checks if the user is authenticated.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		metadata, err := auth.ExtractTokenMetadata(c.Request)
		if err != nil {
			utils.ResponseErrorWithStatus(c, http.StatusUnauthorized, "user hasn't logged in yet")
			c.Abort()
			return
		}

		setUserContext(c, metadata.UserID)
		c.Next()
	}
}

// Authorize checks if the user is authorized.
func Authorize(obj auth.AuthObject, act auth.AuthAction, enforcer auth.AuthEnforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		metadata, err := auth.ExtractTokenMetadata(c.Request)
		if err != nil {
			utils.ResponseErrorWithStatus(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		setUserContext(c, metadata.UserID)
		ok, err := enforcer.Enforce(metadata.UserID, obj, act)
		if err != nil {
			utils.ResponseErrorWithStatus(c, http.StatusUnauthorized, "error occurred while authorizing user")
			c.Abort()
			return
		}

		if !ok {
			utils.ResponseErrorWithStatus(c, http.StatusForbidden, "forbidden")
			c.Abort()
			return
		}

		c.Next()
	}
}

func setUserContext(c *gin.Context, userID string) {
	c.Set(UserIDKey, userID)

	requestLogger := log.Ctx(c.Request.Context()).With().Str(UserIDKey, userID).Logger()
	c.Request = c.Request.WithContext(requestLogger.WithContext(c.Request.Context()))
}
