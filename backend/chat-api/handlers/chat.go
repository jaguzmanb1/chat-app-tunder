package handlers

import (
	"go-chat/data"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/hashicorp/go-hclog"
)

// Message structure that describes a message on the chat
type Message struct {
	To      int    `json:"to"`
	From    int    `json:"from"`
	Message string `json:"message"`
}

// Chat handler for chatting with users
type Chat struct {
	l         hclog.Logger
	clients   map[int]data.User
	broadcast chan Message
}

var upgrader = websocket.Upgrader{}

// New creates a new instance of a chat server
func New(l hclog.Logger) *Chat {
	c := make(map[int]data.User) // Connected clients
	m := make(chan Message)
	return &Chat{l, c, m}
}

// HandleMessages recieve messages and send them to a user
func (h *Chat) HandleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-h.broadcast
		h.l.Info("Recieving new message from client")

		// Send it out to the client with the id specified on the to parameter
		for i, u := range h.clients {
			if i == msg.To {
				h.l.Debug("User found to send message", "id", i)

				err := u.WS.WriteJSON(msg)
				if err != nil {
					h.l.Error("Error sending a message to client", "error", err)
					u.WS.Close()
				}
			}

		}
	}
}

//HandleConnections Handle new clined connections to server
func (h *Chat) HandleConnections(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ui, err := strconv.Atoi(vars["id"])

	h.l.Info("Recieving new connection from", "id", ui)

	if err != nil {
		h.l.Error("Could not parse id header, invalid id")
	}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.l.Error("Could not connect new client: ", "error", err)
	}

	defer ws.Close()
	user := data.User{ID: ui, WS: ws}
	h.clients[ui] = user
	for {
		var msg Message

		err := ws.ReadJSON(&msg)
		if err != nil {
			h.l.Error("Unable to read JSON message fom client", "error", err)
			break
		}

		// Send the newly received message to the broadcast channel
		h.broadcast <- msg
	}

}

//Test Handle a test for the webpage
func (h *Chat) Test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Test example"))
}
