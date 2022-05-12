package main

import (
	"chat-app/api/db"
	"chat-app/api/messages"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	db.ConnectDB()

	router := gin.Default()
	v1 := router.Group("/api/v1")
	v1.GET("/health", health)

	messages.MessageRoutes(v1.Group("/messages"))

	router.Run("0.0.0.0:3001")
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{ "status": "Ok", "message": "Server is running" })
}
