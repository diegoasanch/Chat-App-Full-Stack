package main

import (
	"chat-app/api/db"
	"chat-app/api/messages"
	"chat-app/api/users"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	configureEnvironment()
	db.ConnectDB()

	router := gin.Default()
	v1 := router.Group("/api/v1")
	v1.GET("/health", health)

	messages.MessageRoutes(v1.Group("/messages"))
	users.UserRoutes(v1.Group("/users"))

	router.Run(fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT")))
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{ "status": "Ok", "message": "Server is running :) " })
}

func configureEnvironment() {
	err := godotenv.Load(".env")
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.DebugMode
	}

	gin.SetMode(ginMode)
	if err != nil {
		panic("Error loading .env:\n%s" + err.Error())
	}
}
