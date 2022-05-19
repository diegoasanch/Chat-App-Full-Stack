// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import "fmt"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Live rooms, keyed by room id.
	rooms map[string]*Room

	// Inbound messages from the clients.
	broadcast chan *BroadcastRequest

	// Register requests from the clients.
	registerRoom chan *ConnectionRequest

	// Unregister requests from clients.
	unregisterRoom chan *ConnectionRequest

	// Unregister requests from clients.
	unregisterClient chan *Client

}

type ConnectionRequest struct {
	RoomId string `json:"room_id"`
	Client *Client `json:"-"`
}
type BroadcastRequest struct {
	RoomId string `json:"room_id"`
	Message *WsMessage `json:"message"`
}

func NewHub() *Hub {
	return &Hub{
		broadcast:        make(chan *BroadcastRequest),
		registerRoom:     make(chan *ConnectionRequest),
		unregisterRoom:   make(chan *ConnectionRequest),
		unregisterClient: make(chan *Client),
		rooms:            make(map[string]*Room),
	}
}

func (h *Hub) Run() {
	fmt.Println("[hub]: Starting hub goroutine")
	defer fmt.Println("[hub]: Stopping hub goroutine")
	for {
		select {
		case registerRequest := <-h.registerRoom:
			fmt.Println("[hub]: Register Room")
			room, exists := h.rooms[registerRequest.RoomId]
			if exists {
				fmt.Printf("[hub]: Room %s already exists\n", registerRequest.RoomId)
			} else {
				fmt.Printf("[hub]: Creating room: %s\n", registerRequest.RoomId)
				room = NewRoom(h)
				room.RoomId = registerRequest.RoomId
				h.rooms[room.RoomId] = room
				go room.Run()
			}
			room.register <- registerRequest.Client
		case unregisterRequest := <-h.unregisterRoom:
			fmt.Println("[hub]: Unregister Room")
			if room, ok := h.rooms[unregisterRequest.RoomId]; ok {
				room.unregister <- unregisterRequest.Client
			}
		case unregisterRequest := <-h.unregisterClient:
			fmt.Println("[hub]: Unregister Client")
			for _, room := range h.rooms {
				room.unregister <- unregisterRequest
			}
			unregisterRequest.conn.Close()
			close(unregisterRequest.send)
		case message := <-h.broadcast:
			fmt.Println("[hub]: Broadcast message")
			if room, ok := h.rooms[message.RoomId]; ok {
				fmt.Printf("[hub]: Broadcasting message to room: %s\n", message.RoomId)
				room.broadcast <- message.Message
			} else {
				fmt.Printf("[hub]: Room %s does not exist\n", message.RoomId)
			}
		}
	}
}
