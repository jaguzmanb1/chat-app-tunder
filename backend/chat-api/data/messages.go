package data

import (
	"context"
	"reflect"

	"github.com/hashicorp/go-hclog"
	elastic "github.com/olivere/elastic/v7"
)

// Message structure that describes a message on the chat
type Message struct {
	To      string `json:"to"`
	From    string `json:"from"`
	Message string `json:"message"`
}

// Messages it's a collection of messages
type Messages struct {
	Messages []Message `json:"messages"`
}

//MessageService message service to handle messages persistence
type MessageService struct {
	cl *elastic.Client
	l  hclog.Logger
}

// New creates a new instance of a message service
func New(cl *elastic.Client, l hclog.Logger) *MessageService {
	return &MessageService{cl, l}
}

func returnMessagesArray(sr *elastic.SearchResult) []Message {
	var msgs []Message
	var msg Message

	for _, item := range sr.Each(reflect.TypeOf(msg)) {
		t := item.(Message)
		msgs = append(msgs, t)
	}

	return msgs
}

// GetMessages returns cluster info data
func (m *MessageService) GetMessages(from string) ([]Message, error) {
	ctx := context.Background()
	termQuery := elastic.NewTermQuery("from", from)
	rs, err := m.cl.Search().
		Index("test-messages"). // search in index "twitter"
		Query(termQuery).       // specify the query
		From(0).Size(10).       // take documents 0-9
		Pretty(true).           // pretty print request and response JSON
		Do(ctx)

	if err != nil {
		return nil, err
	}

	ms := returnMessagesArray(rs)

	return ms, nil
}
