package main

import (
	"chat-app/server/api"
	"chat-app/server/db"
	ws "chat-app/server/webSocket"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	configureEnvironment()
	db.Initialize()
	hub := ws.NewHub()
	go hub.Run()

	router := gin.Default()
	ws.WebSocketConnections(router.Group("/ws"), hub)
	api.ApiRoutes(router.Group("/api"))

	router.Run(fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT")))

	defer cleanup()
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

func cleanup() {
	if r := recover(); r != nil {
		fmt.Printf("\n\n---- Web server paniced ----\n\nPanic:\n")
		fmt.Println(r)
		fmt.Println("---- ---- ---- ----")
		panic(r)
	}
}
