package data

import "github.com/gorilla/websocket"

// User describes a user in the chat
type User struct {
	WS *websocket.Conn
	ID int
}
