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
	notes := r.Group("/notes")
	notes.Use(apiData.AuthorizeUserMiddleware())
	{
		notes.GET("/notes", apiData.GetNotes)
		notes.PUT("/note")
		notes.GET("/note")
		notes.PATCH("/note")
		notes.DELETE("/note")
		notes.POST("/reorder")
		notes.PUT("/deadline")
	}
	todo := r.Group("/todo")
	notes.Use(apiData.AuthorizeUserMiddleware())
	{
		todo.GET("/todos")
		todo.PUT("/todo")
		todo.GET("/todo")
		todo.DELETE("/todo")
		todo.POST("/reorder")
		todo.PUT("/deadline")
		todo.POST("/toggle")
	}
	// Listen
	err = r.Run(getListenAddress())
	if err != nil {
		log.WithError(err).Fatalln("cannot serve http server")
	}
}
