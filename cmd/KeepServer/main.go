package main

import (
	"KeepExpandedAndEnhanced/api"
	"KeepExpandedAndEnhanced/pkg/session"
	"github.com/gin-gonic/gin"
)

func main() {
	apiData := new(api.API)
	apiData.SessionStorage = session.NewInMemorySession()
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
}
