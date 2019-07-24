package protocol

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gorilla/websocket"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/metrics"
	"github.com/october93/engine/metrics/dogstatsd"
	"github.com/october93/engine/model"
)

// Connections manages PushWriter by making them accessible by their user ID if
// authenticated.
type Connections struct {
	log log.Logger

	mu           sync.RWMutex
	writers      map[globalid.ID]*PushWriter
	writerByUser map[globalid.ID]map[globalid.ID]*PushWriter

	totalMetric      metrics.Gauge
	uniqueUserMetric metrics.Gauge
	lengthMetric     metrics.TimeHistogram
	lengths          map[globalid.ID]time.Time
}

type Connection struct {
	IPAddress  string         `json:"ipAddress"`
	UserAgent  string         `json:"userAgent"`
	AdminPanel bool           `json:"adminPanel"`
	Session    *model.Session `json:"session"`
	CreatedAt  time.Time      `json:"createdAt"`

	conn *websocket.Conn
}

func NewConnection(ctx context.Context, conn *websocket.Conn) (*Connection, error) {
	userAgent, ok := ctx.Value(UserAgent).(string)
	if !ok {
		return nil, errors.New("unexpected value for user agent")
	}
	ipAddress, ok := ctx.Value(IPAddress).(string)
	if !ok {
		return nil, errors.New("unexpected value for IP address")
	}
	return &Connection{
		IPAddress: ipAddress,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
		conn:      conn,
	}, nil
}

func NewConnections(s *statsd.Client, l log.Logger) *Connections {
	var totalMetric metrics.Gauge
	var uniqueUserMetric metrics.Gauge
	var lengthMetric metrics.TimeHistogram

	if s != nil {
		uniqueUserMetric = dogstatsd.NewGauge("unique_user_connected", 1, s, l)
		totalMetric = dogstatsd.NewGauge("total_connections", 1, s, l)
		lengthMetric = dogstatsd.NewTimeHistogram("connection_length", 1, s, l)
	}
	return &Connections{
		writers:          make(map[globalid.ID]*PushWriter),
		writerByUser:     make(map[globalid.ID]map[globalid.ID]*PushWriter),
		lengths:          make(map[globalid.ID]time.Time),
		totalMetric:      totalMetric,
		uniqueUserMetric: uniqueUserMetric,
		lengthMetric:     lengthMetric,
		log:              l,
	}
}

func (c *Connections) Register(ctx context.Context, conn *websocket.Conn) (*PushWriter, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	connection, err := NewConnection(ctx, conn)
	if err != nil {
		return nil, err
	}
	writer := NewPushWriter(connection, nil, c.log)
	c.writers[writer.ID] = writer
	c.updateMetrics(writer)
	return writer, nil
}

func (c *Connections) Deregister(writer *PushWriter) {
	c.mu.Lock()
	defer c.mu.Unlock()

	writer = c.writers[writer.ID]
	delete(c.writers, writer.ID)
	// remove user mapping if this writer has an associated session
	if writer.Authenticated() {
		userID := writer.Session().UserID
		delete(c.writerByUser[userID], writer.ID)
		c.cleanup(userID)
	}
	c.updateMetrics(writer)
}

func (c *Connections) Authenticate(writer *PushWriter, session *model.Session) {
	c.mu.Lock()
	defer c.mu.Unlock()
	writer.SetSession(session)
	c.writerByUser[session.UserID] = make(map[globalid.ID]*PushWriter)
	c.writerByUser[session.UserID][writer.ID] = writer
}

func (c *Connections) Deauthenticate(writer *PushWriter) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !writer.Authenticated() {
		c.log.Error(errors.New("deauthenticate called on writer without a session"))
		return
	}
	userID := writer.Session().UserID
	delete(c.writerByUser[userID], writer.ID)
	writer.SetSession(nil)
	c.cleanup(userID)
}

func (c *Connections) Writers() []*PushWriter {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var writers []*PushWriter
	for _, writer := range c.writers {
		writers = append(writers, writer)
	}
	return writers
}

func (c *Connections) WritersByUser(userID globalid.ID) []*PushWriter {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var writers []*PushWriter
	for _, writer := range c.writerByUser[userID] {
		writers = append(writers, writer)
	}
	return writers
}

// updateMetrics is only to be called when there is a lock on connections.
func (c *Connections) updateMetrics(writer *PushWriter) {
	if c.totalMetric == nil || c.lengthMetric == nil || c.uniqueUserMetric == nil {
		return
	}
	c.totalMetric.Set(float64(len(c.writers)))
	c.uniqueUserMetric.Set(float64(len(c.writerByUser)))

	var user *model.User
	session := writer.Session()
	if session != nil {
		user = session.GetUser()
	}

	begin, ok := c.lengths[writer.ID]
	if !ok {
		c.lengths[writer.ID] = time.Now()
	} else if user != nil {
		tag := metrics.Tag{Key: "username", Value: user.Username}
		c.lengthMetric.With(tag).Observe(time.Since(begin))
		delete(c.lengths, writer.ID)
	} else {
		c.lengthMetric.Observe(time.Since(begin))
		delete(c.lengths, writer.ID)
	}
}

func (c *Connections) cleanup(userID globalid.ID) {
	if len(c.writerByUser[userID]) == 0 {
		delete(c.writerByUser, userID)
	}
}

func (c *Connections) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.writers)
}

func (c *Connections) EncodeTo(userID globalid.ID, m *Message) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return nil
}
