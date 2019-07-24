package protocol

import (
	"io"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/model"
)

// PushWriter is a writer which is unique for every session, respectively for
// every unauthenticated request. Besides session it tracks callbacks which
// have been registered to the message bus.
type PushWriter struct {
	ID   globalid.ID
	conn *Connection
	log  log.Logger

	mu      sync.RWMutex
	session *model.Session
	writer  io.Writer
}

// NewPushWriter returns a new instance of PushWriter. An optional writer can
// be passed for testing purposes.
func NewPushWriter(conn *Connection, w io.Writer, log log.Logger) *PushWriter {
	if w == nil {
		w = WebSocketWriter{Conn: conn.conn}
	}
	return &PushWriter{
		ID:     globalid.Next(),
		conn:   conn,
		log:    log,
		writer: w,
	}
}

// Write is a write operation which can be used concurrently.
func (pw *PushWriter) Write(data []byte) (int, error) {
	pw.mu.Lock()
	defer pw.mu.Unlock()
	return pw.writer.Write(data)
}

// SetSession allows the session to be set via go routines
func (pw *PushWriter) SetSession(s *model.Session) {
	pw.mu.Lock()
	defer pw.mu.Unlock()
	pw.session = s
	pw.conn.Session = s
}

func (pw *PushWriter) Session() *model.Session {
	pw.mu.RLock()
	defer pw.mu.RUnlock()
	return pw.session
}

func (pw *PushWriter) Authenticated() bool {
	pw.mu.RLock()
	defer pw.mu.RUnlock()
	return pw.session != nil
}

func (pw *PushWriter) Connection() *Connection {
	return pw.conn
}

// WebSocketWriter is a writer wrapper for a given WebSocket connection.
type WebSocketWriter struct {
	Conn *websocket.Conn
}

// Write is the basic operation for sending messages on the WebSocket connection.
func (wsw WebSocketWriter) Write(data []byte) (int, error) {
	err := wsw.Conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return 0, err
	}
	return len(data), nil
}
