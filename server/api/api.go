package api

import (
	"chat-app/server/api/messages"
	"chat-app/server/api/users"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApiRoutes(router *gin.RouterGroup) {
	v1 := router.Group("/v1")
	v1.GET("/health", health)

	messages.MessageRoutes(v1.Group("/messages"))
	users.UserRoutes(v1.Group("/users"))
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{ "status": "Ok", "message": "Server is running :) " })
}
