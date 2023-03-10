package main

import (
	"KeepExpandedAndEnhanced/api"
	"KeepExpandedAndEnhanced/internal/database"
	"KeepExpandedAndEnhanced/pkg/session"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.TraceLevel)
	apiData := new(api.API)
	apiData.UploadPath = getUploadDir()
	apiData.SessionStorage = session.NewInMemorySession()
	// Connect to database
	db, err := database.NewPostgresDatabase(getDatabaseConnectionURL())
	if err != nil {
		log.WithError(err).Fatalln("cannot initialize database")
	}
	apiData.Database = database.NewDatabase(db)
	// Make gin
	r := gin.Default()
	r.Use(api.CORS())
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
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
		notes.GET("/note", apiData.GetNote)
		notes.POST("/note", apiData.AddNote)
		notes.PATCH("/note", apiData.EditNote)
		notes.DELETE("/note", apiData.DeleteNote)
		notes.PATCH("/reorder", apiData.ReorderNote)
	}
	todo := r.Group("/todos")
	todo.Use(apiData.AuthorizeUserMiddleware())
	{
		todo.GET("/todos", apiData.GetTodos)
		todo.POST("/todo", apiData.AddTodo)
		todo.GET("/todo", apiData.GetTodo)
		todo.PATCH("/todo", apiData.EditTodo)
		todo.DELETE("/todo", apiData.DeleteTodo)
		todo.PATCH("/reorder", apiData.ReorderTodo)
	}
	// File upload
	r.POST("/upload", apiData.AuthorizeUserMiddleware(), apiData.UploadFile)
	r.Static("/files", apiData.UploadPath)
	// Listen
	err = r.Run(getListenAddress())
	if err != nil {
		log.WithError(err).Fatalln("cannot serve http server")
	}
}
