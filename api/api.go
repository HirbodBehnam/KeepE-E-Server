package api

import (
	"KeepExpandedAndEnhanced/internal/database"
	"KeepExpandedAndEnhanced/pkg/session"
)

type API struct {
	SessionStorage session.Interface
	Database       database.Database
}
