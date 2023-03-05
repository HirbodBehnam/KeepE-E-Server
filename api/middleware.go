package api

import (
	"KeepExpandedAndEnhanced/pkg/session"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const userIDAuthedContext = "userID"

func (api *API) AuthorizeUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the header
		const prefix = "Bearer "
		header := c.Request.Header.Get("Authorization")
		if !strings.HasPrefix(header, prefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{"empty auth"})
			return
		}
		header = header[len(prefix):]
		// Authorize
		userID, err := api.SessionStorage.Authorize(header)
		if err == session.UnauthorizedTokenErr {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{"bad auth"})
			return
		} else if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{internalServerErr})
			return
		}
		// Set in map
		c.Set(userIDAuthedContext, userID)
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
