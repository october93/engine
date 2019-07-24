package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gorilla/mux"
	"github.com/october93/engine/gql"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/rpc"
	"github.com/october93/engine/rpc/protocol"
	"github.com/october93/engine/store"
	"go.uber.org/ratelimit"
)

// Server is the HTTP server which is used to upgrade connections to the
// WebSocket protocol for the RPC protocol as well as accepting GraphQL
// queries.
type Server struct {
	Store        *store.Store
	GraphQL      *gql.GraphQL
	RPC          rpc.RPC
	router       *protocol.Router
	statsdClient *statsd.Client
	config       *Config
	Log          log.Logger
	server       http.Server
	ctx          context.Context
}

// NewServer returns a new instance of Server.
func NewServer(s *store.Store, g *gql.GraphQL, r rpc.RPC, rt *protocol.Router, sc *statsd.Client, c *Config, l log.Logger, ctx context.Context) *Server {
	return &Server{Store: s, GraphQL: g, RPC: r, router: rt, statsdClient: sc, config: c, Log: l, ctx: ctx}
}

// ListenAndServe listens on the configured TCP network address and handles
// incoming HTTP requests which are translated to Remote Procedure Calls (RPC).
// ListenAndServe configures the specific set of RPCs made available at the
// endpoint.
func (s *Server) Open() error {
	s.registerRPC()
	addr := net.JoinHostPort(s.config.Host, strconv.Itoa(s.config.Port))
	s.Log.Info("starting rpc service", "host", s.config.Host, "port", s.config.Port, "address", addr)

	httpRouter := s.registerHTTP(s.router)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.server = http.Server{
		Handler: httpRouter,
	}
	go func() {
		err := s.server.Serve(listener)
		if err != http.ErrServerClosed {
			s.Log.Error(err)
		} else {
			s.Log.Info("rpc http server gracefully stopped")
		}
	}()
	return nil
}

func (s *Server) Close() error {
	s.Log.Info("shutting down rpc server")
	return s.server.Shutdown(s.ctx)
}

func (s *Server) registerHTTP(router *protocol.Router) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.Handle("/graphql", accessControl(s.requireAdmin(s.GraphQL.SetupHandler(gql.MakeExecutableSchema(s.GraphQL)))))

	// register static file server for serving card images
	fs := http.FileServer(http.Dir(fmt.Sprintf("%s/", s.config.PublicPath)))
	tpl := fmt.Sprintf("/%s/", s.config.PublicPath)
	r.PathPrefix(tpl).Handler(http.StripPrefix(tpl, fs))

	// ping endpoint is used by monitoring services
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Pong!"))
		if err != nil {
			s.Log.Error(err, "path", "/ping")
		}
	})

	// register WebSocket endpoint
	r.Handle("/deck_endpoint/", accessControl(router.UpgradeConnection()))
	r.Handle("/api/", accessControl(router.UpgradeConnection()))
	return r
}

func (s *Server) registerRPC() {
	requireUser := RequireUser(s.Store)
	requireAdmin := RequireAdmin(s.Store)

	auth := Auth(s.router.Connections())
	deauth := Deauth(s.router.Connections())

	// 10 requests per second
	rt := ratelimit.New(10)
	rateLimit := RateLimit(rt)

	// Unauthenticated
	s.router.RegisterRPC("signup", auth(AuthEndpoint(s.RPC)))
	s.router.RegisterRPC("login", auth(AuthEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.Auth, auth(AuthEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.ResetPassword, ResetPasswordEndpoint(s.RPC))
	s.router.RegisterRPC(rpc.ValidateInviteCode, rateLimit(ValidateInviteCodeEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.ValidateUsername, ValidateUsernameEndpoint(s.RPC))
	s.router.RegisterRPC(rpc.AddToWaitlist, AddToWaitlistEndpoint(s.RPC))
	// Authenticated
	s.router.RegisterRPC(rpc.Logout, deauth(requireUser(LogoutEndpoint(s.RPC))))
	s.router.RegisterRPC(rpc.GetCards, requireUser(GetCardsEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.ReactToCard, requireUser(ReactToCardEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.VoteOnCard, requireUser(VoteOnCardEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.PostCard, requireUser(PostCardEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.NewInvite, requireUser(NewInviteEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.GetCard, GetCardEndpoint(s.RPC))
	s.router.RegisterRPC(rpc.GetThread, GetThreadEndpoint(s.RPC))
	s.router.RegisterRPC(rpc.RegisterDevice, requireUser(RegisterDeviceEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.UnregisterDevice, requireUser(UnregisterDeviceEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.UpdateSettings, requireUser(UpdateSettingsEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.GetUser, GetUserEndpoint(s.RPC))
	s.router.RegisterRPC(rpc.GetNotifications, requireUser(GetNotificationsEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.UpdateNotifications, requireUser(UpdateNotificationsEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.GetAnonymousHandle, requireUser(GetAnonymousHandleEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.DeleteCard, requireUser(DeleteCardEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.FollowUser, requireUser(FollowUserEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.UnfollowUser, requireUser(UnfollowUserEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.GetFollowingUsers, requireUser(GetFollowingUsersEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.GetPostsForUser, GetPostsForUserEndpoint(s.RPC))
	s.router.RegisterRPC(rpc.GetTags, GetTagsEndpoint(s.RPC))
	s.router.RegisterRPC(rpc.GetFeaturesForUser, GetFeaturesForUserEndpoint(s.RPC))
	s.router.RegisterRPC(rpc.PreviewContent, requireUser(PreviewContentEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.UploadImage, requireUser(UploadImageEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.GetTaggableUsers, requireUser(GetTaggableUsersEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.ModifyCardScore, requireAdmin(ModifyCardScoreEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.GetInvites, requireUser(GetInvitesEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.GetOnboardingData, requireUser(GetOnboardingDataEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.GetMyNetwork, requireUser(GetMyNetworkEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.UnsubscribeFromCard, requireUser(UnsubscribeFromCardEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.SubscribeToCard, requireUser(SubscribeToCardEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.GroupInvites, requireAdmin(GroupInvitesEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.ReportCard, requireUser(ReportCardEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.BlockUser, requireUser(BlockUserEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.GetChannels, requireUser(GetChannelsEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.GetCardsForChannel, requireUser(GetCardsForChannelEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.UpdateChannelSubscription, requireUser(UpdateChannelSubscriptionEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.JoinChannel, requireUser(JoinChannelEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.LeaveChannel, requireUser(LeaveChannelEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.MuteChannel, requireUser(MuteChannelEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.UnmuteChannel, requireUser(UnmuteChannelEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.MuteUser, requireUser(MuteUserEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.UnmuteUser, requireUser(UnmuteUserEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.MuteThread, requireUser(MuteThreadEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.UnmuteThread, requireUser(UnmuteThreadEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.CreateChannel, requireUser(CreateChannelEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.GetPopularCards, requireUser(GetPopularCardsEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.ValidateChannelName, requireUser(ValidateChannelNameEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.GetChannel, requireUser(GetChannelEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.GetActionCosts, requireUser(GetActionCostsEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.UseInviteCode, requireUser(UseInviteCodeEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.RequestValidation, requireUser(RequestValidationEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.ConfirmValidation, requireUser(ConfirmValidationEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.CanAffordAnonymousPost, requireUser(CanAffordAnonymousPostEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.GetLeaderboard, requireUser(GetLeaderboardEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.SubmitFeedback, requireUser(SubmitFeedbackEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.TipCard, requireUser(TipCardEndpoint(s.RPC)))

	// Admin panel
	s.router.RegisterRPC(rpc.NewUser, requireAdmin(NewUserEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.GetUsers, requireAdmin(GetUsersEndpoint(s.RPC)))
	s.router.RegisterRPC(rpc.ConnectUsers, requireAdmin(ConnectUsersEndpoint(s.RPC)))
}

// Authenticate provides an authentication middleware to secure handlers.
func (s *Server) requireAdmin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokens, ok := r.Header["Authorization"]
		if !ok || len(tokens) == 0 {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		id, err := globalid.Parse(tokens[0])
		if err != nil {
			http.Error(w, "error parsing session ID", http.StatusUnauthorized)
			return
		}
		session, err := s.Store.GetSession(id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		user, err := s.Store.GetUser(session.UserID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if !user.Admin {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// Middleware to activate CORS
func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}
