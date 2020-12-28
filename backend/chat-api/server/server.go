package server

import (
	"go-chat/client"
	"go-chat/data"

	"github.com/hashicorp/go-hclog"
)

// Server define a server object
type Server struct {
	l hclog.Logger
	b chan data.Message
	c map[string]*client.Client
}

// New creates a new server instance
func New(l hclog.Logger, b chan data.Message) *Server {
	c := make(map[string]*client.Client)
	return &Server{l, b, c}
}

// AddCurrentClient adds a client to server clients map
func (s *Server) AddCurrentClient(client *client.Client) {
	s.c[client.US.Phone] = client
}

// HandleMessages recieve messages from the clients
func (s *Server) HandleMessages() {
	for {
		msg := <-s.b
		if s.c[msg.To] != nil {
			s.c[msg.To].SendMessage(msg.Message)
		} else {

		}
	}
}
