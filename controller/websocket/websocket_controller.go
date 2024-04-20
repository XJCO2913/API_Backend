package websocket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebsocketController struct {
	Pool         map[string]*websocket.Conn
	ConnectCh    chan *websocket.Conn
	DisconnectCh chan string
	Upgrader     websocket.Upgrader
}

func NewWebsocketController() *WebsocketController {
	return &WebsocketController{
		Pool:         make(map[string]*websocket.Conn),
		ConnectCh:    make(chan *websocket.Conn),
		DisconnectCh: make(chan string),
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

type msg struct {
	Type   string `json:"type"`
	UserID string `json:"userID"`
}

func (wsc *WebsocketController) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := wsc.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read failed:", err)
			break
		}

		m := msg{}
		err = json.Unmarshal(message, &m)
		if err != nil {
			log.Println("JSON unmarshal failed:", err)
			continue
		}

		switch m.Type {
		case "connect":
			wsc.ConnectCh <- conn
			wsc.Pool[m.UserID] = conn
		case "disconnect":
			wsc.DisconnectCh <- m.UserID
			delete(wsc.Pool, m.UserID)
		case "user_status":
			if _, ok := wsc.Pool[m.UserID]; ok {
				conn.WriteJSON(map[string]interface{}{
					"userStatus": "online",
				})
			} else {
				conn.WriteJSON(map[string]interface{}{
					"userStatus": "offline",
				})
			}
		}

		err = conn.WriteMessage(websocket.TextMessage, []byte("回复消息"))
		if err != nil {
			log.Println("Write failed:", err)
			break
		}
	}
}
