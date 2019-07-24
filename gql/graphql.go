//go:generate gorunpkg github.com/vektah/gqlgen -typemap types.json -schema ./schema.graphql
package gql

import (
	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/model"
	"github.com/october93/engine/rpc"
	"github.com/october93/engine/rpc/notifications"
	"github.com/october93/engine/rpc/protocol"
	"github.com/october93/engine/rpc/push"
	"github.com/october93/engine/search"
	"github.com/october93/engine/store"
	"github.com/october93/engine/worker"
)

// GraphQL gives access to a GraphQL handler in order to setup an endpoint for
// processing GraphQL queries.
type GraphQL struct {
	Store          *store.Store
	RPC            rpc.RPC
	router         *protocol.Router
	Notifier       *worker.Notifier
	Notifications  *notifications.Notifications
	Pusher         *push.Pusher
	ImageProcessor *rpc.ImageProcessor
	Indexer        *search.Indexer

	settings *model.Settings
	config   *Config
	log      log.Logger
}

// NewGraphQL returns a new instance of GraphQL.
func NewGraphQL(s *store.Store, r rpc.RPC, router *protocol.Router, ip *rpc.ImageProcessor, settings *model.Settings, l log.Logger, n *worker.Notifier, p *push.Pusher, i *search.Indexer, c *Config) *GraphQL {
	return &GraphQL{
		Store:          s,
		RPC:            r,
		router:         router,
		ImageProcessor: ip,
		settings:       settings,
		log:            l,
		Notifier:       n,
		Notifications:  notifications.NewNotifications(s, "", 10000),
		Pusher:         p,
		Indexer:        i,
		config:         c,
	}
}
