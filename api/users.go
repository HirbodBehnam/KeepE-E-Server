package api

import (
	"KeepExpandedAndEnhanced/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5/pgconn"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func (api *API) UserLogin(c *gin.Context) {
	// Get the form
	var request signupLoginRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: bodyParseErr + err.Error()})
		return
	}
	// Check username and password
	userID, err := api.Database.UserLogin(request.Username, request.Password)
	if errors.Is(err, database.InvalidCredentialsErr) {
		c.JSON(http.StatusForbidden, errorResponse{Error: database.InvalidCredentialsErr.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{internalServerErr})
		log.WithError(err).WithField("request", request).Error("cannot authorize user")
		return
	}
	// Create redis entry
	accessToken, err := api.SessionStorage.Store(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{internalServerErr})
		log.WithError(err).WithField("request", request).Error("cannot store user auth tokens")
		return
	}
	// Done!
	c.JSON(http.StatusOK, loginResponse{
		Token: accessToken,
	})
}

func (api *API) UserSignup(c *gin.Context) {
	var request signupLoginRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{bodyParseErr + err.Error()})
		return
	}
	err := api.Database.SignupUser(request.Username, request.Password)
	// Check for duplicate
	if pgErr, ok := err.(*pgconn.PgError); ok {
		if pgErr.Code == "23505" {
			c.JSON(http.StatusConflict, errorResponse{userExistsErr})
			return
		}
	}
	// General error
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{internalServerErr})
		log.WithError(err).WithField("request", request).Error("cannot signup user")
		return
	}
	// Done
	c.JSON(http.StatusOK, gin.H{})
}

func (api *API) UserLogout(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		_ = api.SessionStorage.Delete(auth[len("Bearer "):])
	}
	c.JSON(http.StatusOK, gin.H{})
}
