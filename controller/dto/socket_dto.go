package dto

import "github.com/gorilla/websocket"

type Msg struct {
	Type     string
	SenderID string
	Data     map[string]interface{}
}

type ConnectionEvent struct {
	UserID string
	Conn   *websocket.Conn
}
