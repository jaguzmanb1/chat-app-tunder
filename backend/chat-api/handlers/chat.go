package handlers

import (
	"go-chat/client"
	"go-chat/data"
	"go-chat/server"
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/websocket"
	"github.com/hashicorp/go-hclog"
)

// Chat handler for chatting with users
type Chat struct {
	l hclog.Logger
	s server.Server
	b chan data.Message
}

// New creates a new instance of a chat server
func New(l hclog.Logger, s server.Server, b chan data.Message) *Chat {
	return &Chat{l, s, b}
}

// HandleConnections Handle new client connections to server
func (h *Chat) HandleConnections(w http.ResponseWriter, r *http.Request) {
	h.l.Info("[HandleConnections] Recieving new connection")
	var ws = (context.Get(r, "ws")).(*websocket.Conn)
	var us = (context.Get(r, "us")).(data.User)
	ws.SetReadDeadline(time.Time{})

	c := client.New(h.l, ws, h.b, us)
	h.s.AddCurrentClient(c)

	c.HandleConnection()
}

//Test Handle a test for the webpage
func (h *Chat) Test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Test example"))
}
