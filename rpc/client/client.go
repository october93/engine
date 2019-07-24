package client

import (
	"encoding/json"
	"errors"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/rpc"
	"github.com/october93/engine/rpc/protocol"
)

// Client implements the RPC interface and provides a programmatic access to
// the backend API.
type Client struct {
	w      *WebSocketWriter
	config Config
	log    log.Logger

	sync.RWMutex
	requests map[globalid.ID]chan string
	sessions map[globalid.ID]globalid.ID
}

type WebSocketWriter struct {
	sync.Mutex
	conn *websocket.Conn
}

func (w *WebSocketWriter) WriteJSON(v interface{}) error {
	w.Lock()
	defer w.Unlock()
	return w.conn.WriteJSON(v)
}

func (w *WebSocketWriter) ReadMessage() (int, []byte, error) {
	return w.conn.ReadMessage()
}

// NewClient returns a new instance of Client.
func NewClient(config Config, log log.Logger) (rpc.RPC, error) {
	u, err := url.Parse(config.Address)
	if err != nil {
		return nil, err
	}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	client := &Client{
		config:   config,
		w:        &WebSocketWriter{conn: conn},
		requests: make(map[globalid.ID]chan string),
		sessions: make(map[globalid.ID]globalid.ID),
		log:      log}

	go client.messageLoop()
	return client, nil
}

func (c *Client) messageLoop() {
	for {
		_, message, err := c.w.ReadMessage()
		if err != nil {
			c.log.Error(err)
			break
		}
		var resp protocol.Message
		err = json.Unmarshal(message, &resp)
		if err != nil {
			c.log.Error(err)
		}
		// ignore push messages for now
		if resp.RequestID != globalid.Nil {
			continue
		}
		c.Lock()
		responseChannel := c.requests[resp.Ack]
		c.Unlock()
		if responseChannel == nil {
			c.log.Info("unmatched response",
				"ack", resp.Ack,
				"response", string(message))
			continue
		}
		// send matched response back to the method which is responsering on the
		// response channel
		responseChannel <- string(message)
	}
}

// NewMessage generates a new RPC message. NewMessage ensures that a unique
// request ID has been generated.
func NewMessage(rpc string) *protocol.Message {
	return &protocol.Message{RPC: rpc, RequestID: globalid.Next()}
}

func (c *Client) registerRequest(id globalid.ID) <-chan string {
	c.Lock()
	defer c.Unlock()
	c.requests[id] = make(chan string)
	return c.requests[id]
}

func (c *Client) unregisterRequest(id globalid.ID) {
	c.Lock()
	defer c.Unlock()
	delete(c.requests, id)
}

// DecodeGenericResponse unmarshales a response into the generic Message format
// in order to parse header information.
func (c *Client) DecodeGenericResponse(resp string) (*protocol.Message, error) {
	var m *protocol.Message
	err := json.Unmarshal([]byte(resp), &m)
	if err != nil {
		return nil, err
	}
	if m.Err != "" {
		return nil, errors.New(m.Err)
	}
	return m, nil
}
