package main

import (
	"fmt"
	"time"

	"api.backend.xjco2913/controller/ws"
	"api.backend.xjco2913/util/zlog"
	"go.uber.org/zap"
)

type Hub struct {
	Pool map[string]*ws.Client
}

var (
	localHub = NewHub()
)

func NewHub() *Hub {
	return &Hub{
		Pool: ws.Pool,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case event := <-ws.ConnectCh:
			fmt.Println("11111")
			if _, exists := h.Pool[event.UserID]; !exists {
				h.Pool[event.UserID] = &ws.Client{
					UserID:    event.UserID,
					Conn:      event.Conn,
					LastHeart: time.Now(),
					IsAdmin:   false,
				}
				zlog.Info("Client connected", zap.String("userID", event.UserID))
				h.broadcastToAdmins(`{"Type":"new_online", "userID":"` + event.UserID + `"}`)
			}

		case event := <-ws.DisconnectCh:
			zlog.Info("Client disconnected", zap.String("userID", event.UserID))
			h.broadcastToAdmins(`{"Type": "new_offline", "userID":"` + event.UserID + `"}`)

		case msg := <-ws.ServicesCh:
			fmt.Println(333333, msg.Type)
			switch msg.Type {
			case "user_status":
				fmt.Println(888998899)
				if client, ok := h.Pool[msg.SenderID]; ok {
					onlineUsers := []string{}
					for userID, client := range h.Pool {
						// Check if each user is online
						if client != nil && client.Conn != nil {
							onlineUsers = append(onlineUsers, userID)
						}
					}
					// An array of IDs of all online users
					fmt.Println(onlineUsers)
					client.Conn.WriteJSON(map[string]interface{}{
						"onlineUsers": onlineUsers,
					})
					fmt.Println("finish")
				}
				// Other service cases TBD
			}
		}
	}
}

func (h *Hub) broadcastToAdmins(message string) {
	for _, client := range h.Pool {
		if !client.IsAdmin {
			if client.Conn != nil {
				client.Conn.WriteMessage(1, []byte(message))
			}
		}
	}
}
