package api

import (
	"chat-app/server/api/messages"
	"chat-app/server/api/users"
	"net/http"

	"github.com/gin-gonic/gin"
)


func ApiRoutes(router *gin.RouterGroup, authMiddleWare gin.HandlerFunc) {
	v1 := router.Group("/v1")
	v1.GET("/health", health)

	users.UserRoutes(v1.Group("/users"))

	protecterRoutes := v1.Group("/", authMiddleWare)
	messagesRouter := protecterRoutes.Group("/messages")
	messages.MessageRoutes(messagesRouter)
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{ "status": "Ok", "message": "Server is running :) " })
}
