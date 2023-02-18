package api

import (
	"KeepExpandedAndEnhanced/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func (api *API) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{"cannot get file: " + err.Error()})
		return
	}
	// Save the file
	fileID := uuid.NewString()
	rawFilename := filepath.Base(file.Filename)
	savePath := path.Join(api.UploadPath, fileID, rawFilename)
	// Create the folder before saving
	err = os.Mkdir(path.Join(api.UploadPath, fileID), 0666)
	if err != nil {
		log.WithError(err).WithField("path", path.Join(api.UploadPath, fileID)).Error("cannot create directory for file upload")
		return
	}
	// Save the file
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{"cannot upload file: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, database.SavedImage{
		ID:       fileID,
		Filename: rawFilename,
	})
}
