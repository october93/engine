package protocol

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

// MessageEndpoint is the fundamental building block for the protocol on the
// receiving respectively serving end. It represents a single RPC method.
type MessageEndpoint func(ctx context.Context, s *model.Session, pw *PushWriter, m Message) error

// Middleware is a building block which takes an endpoint and returns an
// endpoint again.
type Middleware func(endpoint MessageEndpoint) MessageEndpoint

// Message is the basic building of the protocol. A message can be an action, a
// request and a corresponding response.
type Message struct {
	RPC       string          `json:"rpc,omitempty"`
	SessionID globalid.ID     `json:"sessionID,omitempty"`
	RequestID globalid.ID     `json:"requestID,omitempty"`
	Ack       globalid.ID     `json:"ack,omitempty"`
	Data      json.RawMessage `json:"data,omitempty"`
	Err       string          `json:"error,omitempty"`
	Callback  string          `json:"callback,omitempty"`
}

func NewMessage(rpc string) *Message {
	return &Message{
		RPC:       strings.ToLower(rpc),
		RequestID: globalid.Next(),
	}
}

// Error implements the error interface for Message.
func (m *Message) Error() string {
	return m.Err
}

// EncodePayload is a method for the client package to marshal the data part of
// the message.
func (m *Message) EncodePayload(data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	m.Data = b
	return nil
}

func (m *Message) Encode(ctx context.Context, w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(m)
}
