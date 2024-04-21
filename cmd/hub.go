package main

import (
	"time"

	"api.backend.xjco2913/controller/ws"
	"api.backend.xjco2913/util/zlog"
	"go.uber.org/zap"
)

type Hub struct {
	Pool map[string]*ws.Client
}

func NewHub() *Hub {
	return &Hub{}
}

func (h *Hub) Run() {
	for {
		select {
		case event := <-ws.ConnectCh:
			if _, exists := h.Pool[event.UserID]; !exists {
				h.Pool[event.UserID] = &ws.Client{
					UserID:    event.UserID,
					Conn:      event.Conn,
					LastHeart: time.Now(),
					IsAdmin:   false,
				}
				zlog.Info("Client connected", zap.String("userID", event.UserID))
				h.broadcastToAdmins("User connected: " + event.UserID)
			}

		case event := <-ws.DisconnectCh:
			if _, ok := h.Pool[event.UserID]; ok {
				delete(h.Pool, event.UserID)
				zlog.Info("Client disconnected", zap.String("userID", event.UserID))
				h.broadcastToAdmins("User disconnected: " + event.UserID)
			}

		case msg := <-ws.ServicesCh:
			switch msg.Type {
			case "user_status":
				if client, ok := h.Pool[msg.SenderID]; ok {
					onlineUsers := []string{}
					for userID, client := range h.Pool {
						// Check if each user is online
						if client != nil && client.Conn != nil {
							onlineUsers = append(onlineUsers, userID)
						}
					}
					// An array of IDs of all online users
					client.Conn.WriteJSON(map[string]interface{}{
						"onlineUsers": onlineUsers,
					})
				}
				// Other service cases TBD
			}
		}
	}
}

func (h *Hub) broadcastToAdmins(message string) {
	for _, client := range h.Pool {
		if client.IsAdmin {
			if client.Conn != nil {
				client.Conn.WriteJSON(map[string]interface{}{
					"message": message,
				})
			}
		}
	}
}
