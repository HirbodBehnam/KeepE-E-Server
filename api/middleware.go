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
