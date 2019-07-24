package server

import (
	"context"
	"errors"

	"github.com/october93/engine/model"
	"github.com/october93/engine/rpc/protocol"
	"github.com/october93/engine/store"

	"go.uber.org/ratelimit"
)

// RequireUser is a middleware for endpoints to ensure the request is
// authenticated.
func RequireUser(store *store.Store) protocol.Middleware {
	return func(endpoint protocol.MessageEndpoint) protocol.MessageEndpoint {
		return func(ctx context.Context, session *model.Session, pw *protocol.PushWriter, m protocol.Message) error {
			if session == nil {
				return errors.New("unauthenticated request")
			}
			user := session.GetUser()
			if user == nil {
				return errors.New("unauthenticated request")
			}
			ctx = context.WithValue(ctx, "nodeID", user.ID)
			ctx = context.WithValue(ctx, "username", user.Username)
			return endpoint(ctx, session, pw, m)
		}
	}
}

// RequireAdmin is a middleware to ensure the request is authorized for
// admin-only actions.
func RequireAdmin(store *store.Store) protocol.Middleware {
	return func(endpoint protocol.MessageEndpoint) protocol.MessageEndpoint {
		return func(ctx context.Context, session *model.Session, pw *protocol.PushWriter, m protocol.Message) error {
			if session == nil {
				return errors.New("unauthenticated request")
			}
			user := session.GetUser()
			if user == nil {
				return errors.New("unauthenticated request")
			}
			if !user.Admin {
				return errors.New("permission denied")
			}
			ctx = context.WithValue(ctx, "nodeID", user.ID)
			ctx = context.WithValue(ctx, "username", user.Username)
			return endpoint(ctx, session, pw, m)
		}
	}
}

func Auth(c *protocol.Connections) protocol.Middleware {
	return func(endpoint protocol.MessageEndpoint) protocol.MessageEndpoint {
		return func(ctx context.Context, session *model.Session, pw *protocol.PushWriter, m protocol.Message) error {
			err := endpoint(ctx, session, pw, m)
			if err != nil {
				return err
			}
			c.Authenticate(pw, session)
			return nil
		}
	}
}

func Deauth(c *protocol.Connections) protocol.Middleware {
	return func(endpoint protocol.MessageEndpoint) protocol.MessageEndpoint {
		return func(ctx context.Context, session *model.Session, pw *protocol.PushWriter, m protocol.Message) error {
			err := endpoint(ctx, session, pw, m)
			if err != nil {
				return err
			}
			c.Deauthenticate(pw)
			return nil
		}
	}
}

func RateLimit(rt ratelimit.Limiter) protocol.Middleware {
	return func(endpoint protocol.MessageEndpoint) protocol.MessageEndpoint {
		return func(ctx context.Context, session *model.Session, pw *protocol.PushWriter, m protocol.Message) error {
			rt.Take()
			return endpoint(ctx, session, pw, m)
		}
	}
}
