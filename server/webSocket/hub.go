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

	hubCommand chan *HubCommand
}
type ConnectionRequest struct {
	RoomId string `json:"room_id"`
	Client *Client `json:"-"`
}
type BroadcastRequest struct {
	RoomId string `json:"room_id"`
	Message *WsMessage `json:"message"`
}

type HubCommand struct {
	Command string
	Params map[string]interface{}
}

const H_EXIT = "exit"


func NewHub() *Hub {
	return &Hub{
		broadcast:        make(chan *BroadcastRequest),
		registerRoom:     make(chan *ConnectionRequest),
		unregisterRoom:   make(chan *ConnectionRequest),
		unregisterClient: make(chan *Client),
		hubCommand: 	  make(chan *HubCommand),
		rooms:            make(map[string]*Room),
	}
}

func (h *Hub) Run() {
	fmt.Println("[hub]: Starting hub goroutine")
	defer func() {
		fmt.Println("[hub]: Stopping hub goroutine")
		close(h.broadcast)
		close(h.registerRoom)
		close(h.unregisterRoom)
		close(h.unregisterClient)
		for _, room := range h.rooms {
			room.unregister <- nil
		}
	}()

	for {
		select {
		case registerRequest := <-h.registerRoom:
			h.RegisterRoom(registerRequest)
		case unregisterRequest := <-h.unregisterRoom:
			h.UnregisterRoom(unregisterRequest)
		case unregisterRequest := <-h.unregisterClient:
			h.UnregisterClient(unregisterRequest)
		case broadcastRequest := <-h.broadcast:
			h.Broadcast(broadcastRequest.Message)
		case hubCommand := <-h.hubCommand:
			if (hubCommand.Command == H_EXIT) {
				fmt.Println("[hub]: Exiting hub")
				return
			}
		}
	}
}

func (h *Hub) Broadcast(message *WsMessage) {
	fmt.Println("[hub]: Broadcast message")
	if room, ok := h.rooms[message.RoomId]; ok {
		fmt.Printf("[hub]: Broadcasting message to room: %s\n", message.RoomId)
		room.broadcast <- message
	} else {
		fmt.Printf("[hub]: Room %s does not exist\n", message.RoomId)
		message.Owner.send <- &WsMessage{
			Type: "server_error",
			Payload: "Room does not exist",
			RoomId: message.RoomId,
		}
	}
}

func (h *Hub) RegisterRoom(registerRequest *ConnectionRequest) {
	fmt.Println("[hub]: Register Room")
	room, exists := h.rooms[registerRequest.RoomId]
	if !exists {
		fmt.Printf("[hub]: Creating room: %s\n", registerRequest.RoomId)
		room = NewRoom(h)
		room.RoomId = registerRequest.RoomId
		h.rooms[room.RoomId] = room
		go room.Run()
	}
	room.register <- registerRequest.Client
}

func (h *Hub) UnregisterRoom(unregisterRequest *ConnectionRequest) {
	fmt.Println("[hub]: Unregister Room")
	if room, exists := h.rooms[unregisterRequest.RoomId]; exists {
		room.unregister <- unregisterRequest.Client
	}
}

func (h *Hub) UnregisterClient(toUnregister *Client) {
	fmt.Println("[hub]: Unregister Client")
	for _, room := range h.rooms {
		room.unregister <- toUnregister
	}
	toUnregister.conn.Close()
	close(toUnregister.send)
}

func (h *Hub) DeleteRoom(roomId string) {
	fmt.Println("[hub]: Deleting room: ", roomId)
	delete(h.rooms, roomId)
}
