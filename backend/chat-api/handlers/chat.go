package handlers

import (
	"go-chat/client"
	"go-chat/data"
	"go-chat/server"
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/hashicorp/go-hclog"
)

// Chat handler for chatting with users
type Chat struct {
	l  hclog.Logger
	s  server.Server
	b  chan data.Message
	ms data.MessageService
}

// New creates a new instance of a chat server
func New(l hclog.Logger, s server.Server, b chan data.Message, ms data.MessageService) *Chat {
	l.Debug("[New] Creating new instance of chat handler")
	return &Chat{l, s, b, ms}
}

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}

// getPhoneURL returns the phone from the URL
func (h *Chat) getPhoneURL(r *http.Request) string {
	// parse the product id from the url
	h.l.Info("[getPhoneURL] Extracting phone from url", "url", r.URL.Path)

	vars := mux.Vars(r)
	return vars["phone"]
}

// HandleConnections Handle new client connections to server
func (h *Chat) HandleConnections(w http.ResponseWriter, r *http.Request) {
	h.l.Info("[HandleConnections] Recieving new connection")
	var ws = (context.Get(r, "ws")).(*websocket.Conn)
	var us = (context.Get(r, "us")).(data.User)
	ws.SetReadDeadline(time.Time{})

	c := client.New(h.l, ws, h.b, us)

	m, err := h.ms.GetMessages(us.Phone)
	if err != nil {
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		h.l.Error("[HandleConnections] Can't send persisted messages", "error", err)

	}

	c.SendJSON(m)
	h.s.AddCurrentClient(c)

	c.HandleConnection()
}

//GetMessagesFromUser brings all mesagges from a given user phone number
func (h *Chat) GetMessagesFromUser(w http.ResponseWriter, r *http.Request) {
	h.l.Info("[GetMessagesFromUser] Handling request to bring messages from", "host", r.Host, "url", r.URL.Path)
	messages, err := h.ms.GetMessages(h.getPhoneURL(r))
	if err != nil {
		data.ToJSON(&GenericError{Message: err.Error()}, w)
	}
	data.ToJSON(messages, w)
}
