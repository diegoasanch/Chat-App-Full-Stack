package websocket

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func WebSocketConnections(router *gin.RouterGroup, hub *Hub) {
	router.GET("/", func (c *gin.Context) {
		serveWs(c, hub)
	})
	router.GET("/test", wsTest)
}


func wsTest(c *gin.Context) {
	//upgrade get request to websocket protocol
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()
	for {
		//Read Message from client
		mt, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		//If client message is ping will return pong
		switch string(message) {
			case "ping":
				message = []byte("pong")
			case "pong":
				message = []byte("ping")
			case "hello":
				message = []byte("world")
			case "goodbye":
				message = []byte("bye.")
			default:
				message = []byte("unknown message")
		}

		//Response message to client
		err = ws.WriteMessage(mt, message)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

