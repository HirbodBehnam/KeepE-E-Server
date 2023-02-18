package api

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// GetNotes will get all notes of a user
func (api *API) GetNotes(c *gin.Context) {
	userID := c.MustGet(userIDAuthedContext).(uint64)
	notes, err := api.Database.GetNotes(c.Request.Context(), userID)
	if err != nil {
		log.WithError(err).WithField("userID", userID).Error("cannot get notes of user")
		c.JSON(http.StatusInternalServerError, errorResponse{internalServerErr})
		return
	}
	c.JSON(http.StatusOK, notes)
}

// GetNote will get one note of a user
func (api *API) GetNote(c *gin.Context) {
	userID := c.MustGet(userIDAuthedContext).(uint64)
	var requestData getDeleteEntityRequest
	if err := c.BindQuery(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{bodyParseErr + err.Error()})
		return
	}
	note, err := api.Database.GetNote(c.Request.Context(), userID, requestData.ID)
	if err != nil {
		log.WithError(err).WithField("userID", userID).WithField("noteID", requestData.ID).Error("cannot get note of user")
		c.JSON(http.StatusInternalServerError, errorResponse{internalServerErr})
		return
	}
	c.JSON(http.StatusOK, note)
}

func (api *API) DeleteNote(c *gin.Context) {
	userID := c.MustGet(userIDAuthedContext).(uint64)
	var requestData getDeleteEntityRequest
	if err := c.BindQuery(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{bodyParseErr + err.Error()})
		return
	}
	err := api.Database.DeleteNote(c.Request.Context(), userID, requestData.ID)
	if err != nil {
		log.WithError(err).WithField("userID", userID).WithField("noteID", requestData.ID).Error("cannot delete note of user")
		c.JSON(http.StatusInternalServerError, errorResponse{internalServerErr})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (api *API) ReorderNote(c *gin.Context) {
	userID := c.MustGet(userIDAuthedContext).(uint64)
	var requestData reOrderRequest
	if err := c.BindQuery(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{bodyParseErr + err.Error()})
		return
	}
	err := api.Database.ReorderNote(c.Request.Context(), userID, requestData.ID1, requestData.ID2)
	if err != nil {
		log.WithError(err).WithField("userID", userID).WithField("request", requestData).Error("cannot reorder notes of user")
		c.JSON(http.StatusInternalServerError, errorResponse{internalServerErr})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (api *API) SetNoteDeadline(c *gin.Context) {
	userID := c.MustGet(userIDAuthedContext).(uint64)
	var requestData setDeadlineRequest
	if err := c.BindQuery(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{bodyParseErr + err.Error()})
		return
	}
	err := api.Database.SetNoteDeadline(c.Request.Context(), userID, requestData.ID, requestData.Deadline)
	if err != nil {
		log.WithError(err).WithField("userID", userID).WithField("request", requestData).Error("cannot set deadline of a note of user")
		c.JSON(http.StatusInternalServerError, errorResponse{internalServerErr})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
