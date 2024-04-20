package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Hub struct {
	// Manage connections using the Pool from the WebsocketController
	clients map[string]*websocket.Conn

	// Channel to handle new connections
	connectCh chan *websocket.Conn

	// Channel to handle disconnections
	disconnectCh chan string
}

func NewHub(connectCh chan *websocket.Conn, disconnectCh chan string) *Hub {
	return &Hub{
		clients:      make(map[string]*websocket.Conn),
		connectCh:    connectCh,
		disconnectCh: disconnectCh,
	}
}

// Run starts the hub to process incoming connect and disconnect requests
func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.connectCh:
			// Assuming UserID can be fetched from WebSocket connection,
			// This might require additional logic to correctly associate UserID with connection
			userID := conn.RemoteAddr().String() // Placeholder for actual user ID fetch mechanism
			h.clients[userID] = conn
			log.Printf("Client connected: %s", userID)

		case userID := <-h.disconnectCh:
			if conn, ok := h.clients[userID]; ok {
				delete(h.clients, userID)
				conn.Close()
				log.Printf("Client disconnected: %s", userID)
			}
		}
	}
}
