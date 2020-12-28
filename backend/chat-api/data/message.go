package data

// Message structure that describes a message on the chat
type Message struct {
	To      string `json:"to"`
	From    string `json:"from"`
	Message string `json:"message"`
}
