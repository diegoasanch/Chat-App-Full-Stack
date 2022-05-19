package websocket

import "fmt"

type Room struct {
	RoomId string `json:"room_id"`

	hub *Hub

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan *WsMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewRoom(hub *Hub) *Room {
	return &Room{
		broadcast:  make(chan *WsMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		hub:        hub,
	}
}

func (r *Room) Run() {
	fmt.Printf("[room] Starting Room go routine() id: %s goroutine\n", r.RoomId)
	defer fmt.Printf("[room] Stopping Room go routine() id: %s goroutine\n", r.RoomId)
	for {
		select {
		case client := <-r.register:
			r.Register(client)
		case client := <-r.unregister:
			r.Unregister(client)
		case message := <-r.broadcast:
			r.Broadcast(message)
		}
	}
}

func (r *Room) Register(client *Client) {
	fmt.Println("[room]: Register client")
	r.clients[client] = true
	client.send <- "Joined room"
	r.hub.broadcast <- &BroadcastRequest{ RoomId: r.RoomId, Message: &WsMessage{ Type: "chat_message", Payload: "client has joined the room", RoomId: r.RoomId }}
}

func (r *Room) Unregister(client *Client) {
	fmt.Println("[room]: Unregister client")
	delete(r.clients, client)
	client.send <- "You have left the room"
	r.hub.broadcast <- &BroadcastRequest{ RoomId: r.RoomId, Message: &WsMessage{ Type: "chat_message", Payload: "client has left the room", RoomId: r.RoomId }}
}

func (r *Room) Broadcast(message *WsMessage) {
	fmt.Println("[room]: Broadcast message")
	for client := range r.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(r.clients, client)
		}
	}
}
