package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"api.backend.xjco2913/controller/dto"
	"api.backend.xjco2913/util/zlog"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var (
	ConnectCh    = make(chan dto.ConnectionEvent, 1)
	DisconnectCh = make(chan dto.ConnectionEvent, 1)
	ServicesCh   = make(chan dto.Msg, 1)

	Pool = make(map[string]*Client)
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin: func(r *http.Request) bool {
	// 	return true
	// },
}

type Client struct {
	UserID    string
	Conn      *websocket.Conn
	LastHeart time.Time
	IsAdmin   bool
}

type WebsocketController struct{}

func NewWebsocketController() *WebsocketController {
	return &WebsocketController{}
}

func (wsc *WebsocketController) HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Set a persistent connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		zlog.Error("Upgrade failed", zap.Error(err))
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			zlog.Error("Read failed", zap.Error(err))
			break
		}

		fmt.Println(string(message))

		var m dto.Msg
		if err = json.Unmarshal(message, &m); err != nil {
			zlog.Error("JSON unmarshal failed", zap.String("message", string(message)), zap.Error(err))
			continue
		}

		switch m.Type {
		case "connect", "disconnect":
			fmt.Println("here2222")
			event := dto.ConnectionEvent{
				UserID: m.SenderID,
				Conn:   conn,
			}
			if m.Type == "connect" {
				ConnectCh <- event
			} else {
				DisconnectCh <- event
			}
		default:
			fmt.Println("here service")
			ServicesCh <- m
		}
	}
}
