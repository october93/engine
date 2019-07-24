package protocol

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	bugsnag "github.com/bugsnag/bugsnag-go"
	"github.com/gorilla/websocket"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/model"
	"github.com/october93/engine/rpc"
	"github.com/october93/engine/store"
	"go.uber.org/ratelimit"
)

const (
	// RequestID is a unique identifier for every request in order to
	// implement distributed request tracing.
	RequestID = "RequstID"
	// SessionID is used in order to retrieve the current session for the given
	// context.
	SessionID = "SessionID"
	// Callback is used to identify and match up the RPC of the original
	// request. It is helpful for logging and the client.
	Callback = "Callback"
	// IPAddress is the context key used to identify the IP address of a
	// connection.
	IPAddress = "IPAddress"
	// IPAddress is the context key used to identify the user agent of a
	// connection.
	UserAgent = "UserAgent"
)

var connectionUpgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

// Router registers RPCs to be matched and dispatches a handler.
type Router struct {
	endpoints map[string]MessageEndpoint
	store     *store.Store
	settings  *model.Settings
	config    *Config
	log       log.Logger

	conns *Connections
}

// NewRouter returns a new Router instance.
func NewRouter(s *store.Store, conns *Connections, settings *model.Settings, c *Config, l log.Logger) *Router {
	return &Router{
		endpoints: make(map[string]MessageEndpoint),
		conns:     conns,
		store:     s,
		settings:  settings,
		config:    c,
		log:       l,
	}
}

// RegisterRPC adds a new endpoint to the list of available routes. It ensures
// some resilience to using clients by lower casing the RPC name.
func (r *Router) RegisterRPC(rpcName string, e MessageEndpoint) {
	rpcName = strings.ToLower(rpcName)
	r.endpoints[rpcName] = e
}

// UpgradeConnection upgrades the HTTP server connection to the WebSocket
// protocol and handles incoming messages in accordance to the RPC protocol.
func (r *Router) UpgradeConnection() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if r.config.EnableBugSnag {
			defer bugsnag.Recover()
		}
		// upgrades HTTP requests to WebSocket protocol
		r.log.Info("client connected", "session_id", req.FormValue("session"), "ip_address", getIP(req), "user_agent", req.UserAgent())
		conn, err := connectionUpgrader.Upgrade(w, req, nil)
		if err != nil {
			r.log.Error(err)
			if _, err = w.Write([]byte(err.Error())); err != nil {
				r.log.Error(err)
			}
			return
		}
		defer r.Close(conn)

		ctx := context.Background()
		ctx = context.WithValue(ctx, UserAgent, req.UserAgent())
		ctx = context.WithValue(ctx, IPAddress, getIP(req))
		err = r.HandleConnection(ctx, conn, req)
		if err != nil {
			r.log.Error(err)
		}
	}
}

// HandleConnection handles incoming messages in accordance to the RPC
// protocol. Each RPC request is handled in its own go routine.
func (r *Router) HandleConnection(ctx context.Context, conn *websocket.Conn, req *http.Request) error {
	writer, err := r.conns.Register(ctx, conn)
	if err != nil {
		return err
	}
	defer r.conns.Deregister(writer)
	session, err := r.handleSession(ctx, writer, req)
	if err != nil {
		return err
	}
	rt := ratelimit.New(500)

	// incoming requests loop
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		rt.Take()

		settings, err := r.store.GetSettings()
		if err != nil {
			r.log.Error(err)
			settings = &model.Settings{}
		}
		if settings.MaintenanceMode {
			err = r.handleMaintenanceMode(ctx, conn, writer, string(message))
			if err != nil {
				return err
			}
		}

		go r.handleRequest(ctx, session, writer, string(message))
	}
	return nil
}

func (r *Router) handleSession(ctx context.Context, writer *PushWriter, req *http.Request) (*model.Session, error) {
	sessionID := req.FormValue("session")
	adminPanel := req.FormValue("adminpanel")
	if sessionID != "" {
		id, err := globalid.Parse(sessionID)
		if err != nil {
			return nil, err
		}
		session, err := r.store.GetSession(id)
		if err != nil && errors.Cause(err) == sql.ErrNoRows {
			m := NewMessage(rpc.Logout)
			r.log.Info("unknown session, logging out client", "session_id", id)
			return model.NewSession(nil), m.Encode(ctx, writer)
		} else if err != nil {
			return nil, err
		}
		r.conns.Authenticate(writer, session)
		return session, nil
	}
	if adminPanel != "" {
		writer.conn.AdminPanel = true
	}
	return model.NewSession(nil), nil
}

func (r *Router) handleRequest(ctx context.Context, session *model.Session, writer *PushWriter, request string) {
	if r.config.EnableBugSnag {
		defer bugsnag.Recover()
	}
	err := r.Route(ctx, session, writer, request)
	if err != nil {
		r.log.Error(err)
	}
}

func (r *Router) handleMaintenanceMode(ctx context.Context, conn io.Closer, writer io.Writer, request string) error {
	m := &Message{}
	err := json.Unmarshal([]byte(request), &m)
	if err != nil {
		return err
	}
	if m.RPC == "login" {
		return nil
	}

	// try getting the session and ignore maint. mode if it's an admin session
	session, err := r.store.GetSession(m.SessionID)
	if err == nil {
		user, uerr := r.store.GetUser(session.UserID)
		if uerr != nil {
			return uerr
		}
		// maintenance mode does not apply to admins
		if user.Admin {
			return nil
		}
	}

	type maintenanceModeData struct {
		Status bool `json:"status"`
	}

	m = NewMessage("maintainanceMode")
	err = m.EncodePayload(maintenanceModeData{Status: true})
	if err != nil {
		return nil
	}
	err = m.Encode(ctx, writer)
	if err != nil {
		return err
	}
	return conn.Close()
}

// Route is responsible for parsing incoming messages and invokes the endpoint
// after determening it based on the message header.
func (r *Router) Route(ctx context.Context, session *model.Session, writer *PushWriter, request string) error {
	var m Message
	err := json.Unmarshal([]byte(request), &m)
	if err != nil {
		return DefaultEncoder(ctx, nil, err, writer)
	}

	ctx = context.WithValue(ctx, RequestID, m.RequestID)
	ctx = context.WithValue(ctx, SessionID, m.SessionID)
	ctx = context.WithValue(ctx, Callback, m.Callback)

	rpc := strings.ToLower(m.RPC)
	endpoint := r.endpoints[rpc]
	if endpoint == nil {
		return DefaultEncoder(ctx, nil, ErrRPCNotFound(m.RPC), writer)
	}
	err = endpoint(ctx, session, writer, m)
	if err != nil {
		return DefaultEncoder(ctx, nil, err, writer)
	}
	return nil
}

// Close helper to ensure the returned error is checked.
func (r *Router) Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		r.log.Error(err)
	}
}

func (r *Router) Connections() *Connections {
	return r.conns
}

func getIP(req *http.Request) string {
	xff := req.Header.Get("X-Forwarded-For")
	if xff != "" {
		return strings.Split(xff, ",")[0]
	}
	remoteAddr := req.RemoteAddr
	return strings.Split(remoteAddr, ":")[0]
}
