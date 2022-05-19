// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,    //check origin will check the cross region source (note : please not using in production)
	CheckOrigin: func(r *http.Request) bool {
        // TODO: Restrict the origin to our addresses
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan interface{}
}

type WsMessage struct {
	Owner *Client `json:"owner"`
	RoomId string `json:"room_id"`
	Type string `json:"type"`
	Payload interface{} `json:"payload"`
}

const JOIN_ROOM_REQUEST = "join_room"
const LEAVE_ROOM_REQUEST = "leave_room"
const CHAT_MESSAGE = "chat_message"
const LEAVE_SERVER_REQUEST = "leave_server"

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		fmt.Println("[client]: readPump exit")
		c.hub.unregisterClient <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		jsonMessage := WsMessage{}

		if jsonErr := json.Unmarshal(message, &jsonMessage); jsonErr != nil {
			log.Println(jsonErr)
			continue
		}
		jsonMessage.Owner = c

		switch (jsonMessage.Type) {
		case JOIN_ROOM_REQUEST:
			fmt.Println("[client]: Join room request")
			c.hub.registerRoom <- &ConnectionRequest{
				RoomId: jsonMessage.RoomId,
				Client: c,
			}
		case LEAVE_ROOM_REQUEST:
			fmt.Println("[client]: Leave room request")
			c.hub.unregisterRoom <- &ConnectionRequest{
				RoomId: jsonMessage.RoomId,
				Client: c,
			}
		case CHAT_MESSAGE:
			fmt.Println("[client]: Chat message")
			c.hub.broadcast <- &BroadcastRequest{
					RoomId: jsonMessage.RoomId,
					Message: &jsonMessage,
			}
		case LEAVE_SERVER_REQUEST:
			fmt.Println("[client]: Leave server request")
			c.hub.unregisterClient <- c
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		fmt.Println("[client]: writePump exit")
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			fmt.Println("[client]: Sending message: ", message)
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				fmt.Println("[client]: Not ok, The hub closed the channel")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.conn.WriteJSON(message)
			if err != nil {
				fmt.Println("[client]: Error sending message: ", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(c *gin.Context, hub *Hub) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan interface{}, 256)}
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	client.readPump()
	defer conn.Close()
}
