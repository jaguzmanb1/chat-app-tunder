package client

import (
	"go-chat/data"

	"github.com/gorilla/websocket"
	"github.com/hashicorp/go-hclog"
)

// Client define a object of a client connected to the chat server
type Client struct {
	l  hclog.Logger
	ws *websocket.Conn
	b  chan data.Message
	US data.User
}

// New creates a new client instance
func New(l hclog.Logger, ws *websocket.Conn, b chan data.Message, us data.User) *Client {
	l.Debug("[New] Creating new client to", "phone", us.Phone)
	return &Client{l, ws, b, us}
}

// SendMessage sends a message to this client
func (c *Client) SendMessage(m string) error {
	err := c.ws.WriteJSON(m)
	if err != nil {
		return err
	}
	return nil
}

// SendJSON sends a message to this client
func (c *Client) SendJSON(i interface{}) error {
	c.l.Info("[SendJSON] Sending messages data to user", "phone", c.US.Phone)

	err := c.ws.WriteJSON(i)
	if err != nil {
		return err
	}
	return nil
}

// HandleConnection handles a client connection to the server
func (c *Client) HandleConnection() error {
	for {
		var msg data.Message

		err := c.ws.ReadJSON(&msg)
		if err != nil {
			c.ws.Close()
			return err
		}
		c.b <- msg
	}
}
