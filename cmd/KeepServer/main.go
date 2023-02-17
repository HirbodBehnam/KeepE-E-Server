package main

import (
	"KeepExpandedAndEnhanced/api"
	"KeepExpandedAndEnhanced/internal/database"
	"KeepExpandedAndEnhanced/pkg/session"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	apiData := new(api.API)
	apiData.SessionStorage = session.NewInMemorySession()
	// Connect to database
	db, err := database.NewPostgresDatabase(getDatabaseConnectionURL())
	if err != nil {
		log.WithError(err).Fatalln("cannot initialize database")
	}
	apiData.Database = database.NewDatabase(db)
	// Make gin
	r := gin.Default()
	// Login
	users := r.Group("/users")
	{
		users.POST("/login", apiData.UserLogin)
		users.POST("/signup", apiData.UserSignup)
		users.POST("/logout", apiData.UserLogout)
	}
	// General stuff
	general := r.Group("/general")
	general.Use(apiData.AuthorizeUserMiddleware())
	{
		general.PUT("/note/simple")
		general.PUT("/note/todo")
		general.GET("/notes")
		general.DELETE("/note")
		general.GET("/note")
		general.POST("/reorder")
		general.PUT("/deadline")
	}
	// Listen
	err = r.Run(getListenAddress())
	if err != nil {
		log.WithError(err).Fatalln("cannot serve http server")
	}
}
