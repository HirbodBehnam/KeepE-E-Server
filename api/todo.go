package api

import (
	"KeepExpandedAndEnhanced/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (api *API) AddTodo(c *gin.Context) {
	userID := c.MustGet(userIDAuthedContext).(uint64)
	var todo database.Todo
	err := c.BindJSON(&todo)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{bodyParseErr + err.Error()})
		return
	}
	todo, err = api.Database.AddTodo(c.Request.Context(), userID, todo)
	if err != nil {
		log.WithError(err).WithField("userID", userID).WithField("todo", todo).Error("cannot insert todo of user")
		c.JSON(http.StatusInternalServerError, errorResponse{internalServerErr})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (api *API) GetTodos(c *gin.Context) {
	userID := c.MustGet(userIDAuthedContext).(uint64)
	todos, err := api.Database.GetTodos(c.Request.Context(), userID)
	if err != nil {
		log.WithError(err).WithField("userID", userID).Error("cannot get todos of user")
		c.JSON(http.StatusInternalServerError, errorResponse{internalServerErr})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (api *API) GetTodo(c *gin.Context) {
	userID := c.MustGet(userIDAuthedContext).(uint64)
	var requestData getDeleteEntityRequest
	if err := c.BindQuery(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{bodyParseErr + err.Error()})
		return
	}
	todo, err := api.Database.GetTodo(c.Request.Context(), userID, requestData.ID)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	if err != nil {
		log.WithError(err).WithField("userID", userID).WithField("todoID", requestData.ID).Error("cannot get todo of user")
		c.JSON(http.StatusInternalServerError, errorResponse{internalServerErr})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (api *API) EditTodo(c *gin.Context) {
	userID := c.MustGet(userIDAuthedContext).(uint64)
	var todo database.Todo
	if err := c.BindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{bodyParseErr + err.Error()})
		return
	}
	if todo.ID == 0 {
		c.JSON(http.StatusBadRequest, errorResponse{bodyParseErr + "note id must not be zero"})
		return
	}
	// Edit the note
	err := api.Database.EditTodo(c.Request.Context(), userID, todo)
	if err != nil {
		log.WithError(err).WithField("userID", userID).WithField("todo", todo).Error("cannot edit todo of user")
		c.JSON(http.StatusInternalServerError, errorResponse{internalServerErr})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (api *API) DeleteTodo(c *gin.Context) {
	userID := c.MustGet(userIDAuthedContext).(uint64)
	var requestData getDeleteEntityRequest
	if err := c.BindQuery(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{bodyParseErr + err.Error()})
		return
	}
	err := api.Database.DeleteTodo(c.Request.Context(), userID, requestData.ID)
	if err != nil {
		log.WithError(err).WithField("userID", userID).WithField("todoID", requestData.ID).Error("cannot delete todo of user")
		c.JSON(http.StatusInternalServerError, errorResponse{internalServerErr})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (api *API) ReorderTodo(c *gin.Context) {
	userID := c.MustGet(userIDAuthedContext).(uint64)
	var requestData reOrderRequest
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{bodyParseErr + err.Error()})
		return
	}
	err := api.Database.ReorderTodo(c.Request.Context(), userID, requestData.ID1, requestData.ID2)
	if err != nil {
		log.WithError(err).WithField("userID", userID).WithField("request", requestData).Error("cannot reorder todos of user")
		c.JSON(http.StatusInternalServerError, errorResponse{internalServerErr})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
