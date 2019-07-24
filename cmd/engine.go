package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"

	"github.com/DataDog/datadog-go/statsd"
	bugsnag "github.com/bugsnag/bugsnag-go"
	"github.com/october93/engine/coinmanager"
	"github.com/october93/engine/debug"
	"github.com/october93/engine/gql"
	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/rpc"
	"github.com/october93/engine/rpc/notifications"
	"github.com/october93/engine/rpc/protocol"
	"github.com/october93/engine/rpc/push"
	"github.com/october93/engine/rpc/server"
	"github.com/october93/engine/search"
	"github.com/october93/engine/store"
	"github.com/october93/engine/store/datastore"
	"github.com/october93/engine/worker"
	"github.com/october93/engine/worker/activityrecorder"
	"github.com/october93/engine/worker/emailsender"
)

type Engine struct {
	sessionCleaner     *worker.SessionCleaner
	databaseMonitor    *worker.DatabaseMonitor
	reengagementWorker *worker.ReengagementWorker

	server interface {
		Open() error
		Close() error
	}

	Closed chan struct{}

	bp     BuildParameters
	config *Config
	log    log.Logger
}

func NewEngine(bp BuildParameters, config *Config, log log.Logger) *Engine {
	return &Engine{
		Closed: make(chan struct{}),
		bp:     bp,
		config: config,
		log:    log,
	}
}

func (e *Engine) Run() error {
	ctx := context.Background()

	e.log.Info(fmt.Sprintf("engine starting, version %s, branch %s, commit %s", e.bp.Version, e.bp.Branch, e.bp.Commit))
	e.log.Info(fmt.Sprintf("Go version %s, GOMAXPROCS set to %d", runtime.Version(), runtime.GOMAXPROCS(0)))

	var err error
	if isDeployed(e.bp.Version) {
		bugsnag.Configure(bugsnag.Configuration{
			APIKey:          e.config.BugsnagAPIKey,
			ReleaseStage:    e.config.Environment,
			ProjectPackages: []string{"main", "github.com/october93/engine/**"},
		})
	}

	if e.config.Debug.Profile {
		profiler := debug.NewProfiler(e.config.Debug, e.log)
		go func() {
			err = profiler.Start()
			if err != nil {
				e.log.Error(err)
			}
		}()
	}

	err = datastore.SetupDatabase(e.config.Store.Datastore)
	if err != nil {
		return err
	}
	store, err := store.NewStore(&e.config.Store, e.log)
	if err != nil {
		return err
	}
	err = store.EnsureRootUser()
	if err != nil {
		return err
	}
	settings, err := store.EnsureSettings()
	if err != nil {
		return err
	}

	// set up search indexing
	indexer, err := search.NewIndexer(store, e.log, &e.config.Search)
	if err != nil {
		return err
	}

	var statsdClient *statsd.Client
	if e.config.Metrics.Enabled {
		ec2InstanceID := "unknown"
		var resp *http.Response
		resp, err = http.Get("http://169.254.169.254/latest/meta-data/instance-id")
		if err != nil {
			e.log.Error(err)
		}
		var result []byte
		result, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			e.log.Error(err)
		} else {
			ec2InstanceID = string(result)
			e.log = e.log.With("instanceID", ec2InstanceID)
		}
		statsdClient, err = statsd.New(e.config.Metrics.DatadogAgentURL)
		if err != nil {
			return err
		}
		statsdClient.Namespace = "engine."
		statsdClient.Tags = append(statsdClient.Tags, fmt.Sprintf("env:%s", e.config.Environment))
		statsdClient.Tags = append(statsdClient.Tags, fmt.Sprintf("instance-id:%s", ec2InstanceID))
	}

	// set up Router
	connections := protocol.NewConnections(statsdClient, e.log)

	// set up notifier
	nw, err := worker.NewNotifier(store, &e.config.Worker, e.log)
	if err != nil {
		return err
	}

	// set up session cleaner worker
	e.sessionCleaner = worker.NewSessionCleaner(store, e.log)
	e.sessionCleaner.Start()

	// set up worker for sending out emails
	ew, err := emailsender.NewEmailSender(&e.config.Worker, e.log)
	if err != nil {
		return err
	}

	e.reengagementWorker = worker.NewReengagementWorker(store, nw, e.log)

	// set up activity recorder
	ar, err := activityrecorder.NewActivityRecorder(&e.config.Worker, &e.config.Store, e.log)
	if err != nil {
		return err
	}

	// bundle all workers together
	wk := struct {
		*worker.Notifier
		*worker.SessionCleaner
		*emailsender.EmailSender
	}{
		Notifier:       nw,
		SessionCleaner: e.sessionCleaner,
		EmailSender:    ew,
	}

	// set up email templates
	templates, err := emailsender.NewMailTemplates()
	if err != nil {
		return err
	}

	// set up image processor
	ip, err := rpc.NewImageProcessor(&e.config.Server.RPC, e.log)
	if err != nil {
		return err
	}

	// Populate datastore with necessary entities
	err = store.Populate()
	if err != nil {
		return err
	}

	// set up Facebook OAuth2
	oauth2 := rpc.NewOAuth2(e.config.Server.RPC.FacebookAppID, e.config.Server.RPC.FacebookAppSecret)

	router := protocol.NewRouter(store, connections, settings, &e.config.Server.Protocol, e.log)

	// set up Pusher
	p, err := push.NewPusher(connections, store, &e.config.Push, e.log)
	if err != nil {
		return err
	}

	// create notifications builder
	imagePath := fmt.Sprintf("https://s3-%s.amazonaws.com/%s/%s", e.config.Server.RPC.S3Region, e.config.Server.RPC.S3Bucket, e.config.Server.RPC.SystemIconPath)
	ns := notifications.NewNotifications(store, imagePath, e.config.Server.RPC.UnitsPerCoin)

	// Coin manager
	cm := coinmanager.NewCoinManager(store, &e.config.CoinManager)

	resps := rpc.NewResponses(store)

	// set up RPC API and wrap into analytics middleware
	r := rpc.NewRPC(store, wk, oauth2, ip, &e.config.Server.RPC, e.log, templates, nw, p, settings, ns, indexer, cm, resps)

	// set up GraphQL API
	graphql := gql.NewGraphQL(store, r, router, ip, settings, e.log, nw, p, indexer, &e.config.GraphQL)

	// wrap RPC API in instrumenting middleware when enabled
	if e.config.Metrics.Enabled {
		r = rpc.NewInstrumentingMiddleware(r, statsdClient, e.log)
		e.databaseMonitor = worker.NewDatabaseMonitor(store.Store, statsdClient, e.log)
		e.databaseMonitor.Start()
	}
	r = rpc.NewLoggingMiddleware(r, e.log)
	r = rpc.NewRecordingMiddleware(r, ar, e.log)

	e.log.Info("engine started")

	// start server
	e.server = server.NewServer(store, graphql, r, router, statsdClient, &e.config.Server, e.log, ctx)
	return e.server.Open()
}

func (e *Engine) Close() {
	if e.sessionCleaner != nil {
		e.sessionCleaner.Stop()
	}
	if e.databaseMonitor != nil {
		e.databaseMonitor.Stop()
	}
	if err := e.server.Close(); err != nil {
		e.log.Error(err)
	}
	close(e.Closed)
}
