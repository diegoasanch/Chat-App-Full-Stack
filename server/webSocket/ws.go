package ws

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func WebSocketConnections(router *gin.RouterGroup) {
	router.GET("/test", wsTest)
}

var upgrader = websocket.Upgrader{
    //check origin will check the cross region source (note : please not using in production)
	CheckOrigin: func(r *http.Request) bool {
        //Here we just allow the chrome extension client accessable (you should check this verify accourding your client source)
		return true
	},
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

