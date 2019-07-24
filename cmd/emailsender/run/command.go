package run

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	bugsnag "github.com/bugsnag/bugsnag-go"
	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/worker/emailsender"
)

// Command represents the command executed by "emailsender run".
type Command struct {
	Version   string
	Branch    string
	Commit    string
	BuildTime string

	// Channels to use when shutting down gracefully
	closing chan struct{}
	Closed  chan struct{}

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	ctx    context.Context
	Logger log.Logger

	worker interface {
		ConsumeJobs() error
		Shutdown() error
	}
}

// NewCommand return a new instance of Command.
func NewCommand() *Command {
	return &Command{
		ctx:     context.Background(),
		closing: make(chan struct{}),
		Closed:  make(chan struct{}),
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Logger:  log.NopLogger(),
	}
}

func (c *Command) production() bool {
	return c.Version != "" && c.Version != "unknown"
}

var environment = ""

func (cmd *Command) Run(args ...string) error {
	// Parse the command line flags.
	var options Options
	if o, err := cmd.ParseFlags(args...); err != nil {
		return err
	} else {
		options = o
	}

	var config *Config
	// Parse config
	if c, err := cmd.ParseConfig(options.GetConfigPath()); err != nil {
		return fmt.Errorf("parse config: %s", err)
	} else {
		config = c
	}

	// Validate the configuration.
	if err := config.Validate(); err != nil {
		return fmt.Errorf("%s. To generate a valid configuration file run `emailsender config > emailsender.generated.toml`", err)
	}

	logger, err := log.NewLogger(cmd.production(), config.LogLevel)
	if err != nil {
		return err
	}
	logger = logger.With("app", "emailsender")
	cmd.Logger = logger

	// Mark start-up in log.
	cmd.Logger.Info(fmt.Sprintf("emailsender starting, version %s, branch %s, commit %s", cmd.Version, cmd.Branch, cmd.Commit))
	cmd.Logger.Info(fmt.Sprintf("Go version %s, GOMAXPROCS set to %d", runtime.Version(), runtime.GOMAXPROCS(0)))

	// Write the PID file.
	if err = cmd.writePIDFile(options.PIDFile); err != nil {
		return fmt.Errorf("write pid file: %s", err)
	}

	// only run Bugsnag in production environments
	if cmd.production() {
		bugsnag.Configure(bugsnag.Configuration{
			APIKey:       config.BugsnagAPIKey,
			ReleaseStage: environment,
		})
	}
	cmd.worker = emailsender.NewEmailConsumer(&config.Worker, cmd.Logger)
	return cmd.worker.ConsumeJobs()
}

func (cmd *Command) Close() {
	if err := cmd.worker.Shutdown(); err != nil {
		cmd.Logger.Error(err)
	}
	close(cmd.Closed)
}

// ParseFlags parses the command line flags from args and returns an options set.
func (cmd *Command) ParseFlags(args ...string) (Options, error) {
	var options Options
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringVar(&options.ConfigPath, "config", "", "")
	fs.StringVar(&options.PIDFile, "pidfile", "", "")
	fs.Usage = func() {
		if _, err := fmt.Fprintln(cmd.Stderr, usage); err != nil {
			fmt.Println(err)
		}
	}
	if err := fs.Parse(args); err != nil {
		return Options{}, err
	}
	return options, nil
}

// writePIDFile writes the process ID to path.
func (cmd *Command) writePIDFile(path string) error {
	// Ignore if path is not set.
	if path == "" {
		return nil
	}

	// Ensure the required directory structure exists.
	err := os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		return fmt.Errorf("mkdir: %s", err)
	}

	// Retrieve the PID and write it.
	pid := strconv.Itoa(os.Getpid())
	if err := ioutil.WriteFile(path, []byte(pid), 0666); err != nil {
		return fmt.Errorf("write file: %s", err)
	}

	return nil
}

// ParseConfig parses the config at path.
// It returns a demo configuration if path is blank.
func (cmd *Command) ParseConfig(path string) (*Config, error) {
	// Use demo configuration if no config path is specified.
	if path == "" {
		cmd.Logger.Info("no configuration provided, using default settings")
		return NewDefaultConfig()
	}

	cmd.Logger.Info(fmt.Sprintf("Using configuration at: %s", path))

	config := NewConfig()
	if err := config.Load(path); err != nil {
		return nil, err
	}

	return config, nil
}

const usage = `Runs the emailsender.

Usage: emailsender run [flags]

    -config <path>
            Set the path to the configuration file.
            This defaults to the environment variable ENGINE_CONFIG_PATH,
            ~/.emailsender/emailsender.toml, or /etc/emailsender/emailsender.toml if a file
            is present at any of these locations.
            Disable the automatic loading of a configuration file using
            the null device (such as /dev/null).
    -pidfile <path>
            Write process ID to a file.
    -cpuprofile <path>
            Write CPU profiling information to a file.
    -memprofile <path>
            Write memory usage information to a file.
`

// Options represents the command line options that can be parsed.
type Options struct {
	ConfigPath string
	PIDFile    string
}

// GetConfigPath returns the config path from the options.
// It will return a path by searching in this order:
//   1. The CLI option in ConfigPath
//   2. The environment variable ENGINE_CONFIG_PATH
//   3. The first emailsender.toml file on the path:
//        - ~/.emailsender
//        - /etc/emailsender
func (opt *Options) GetConfigPath() string {
	if opt.ConfigPath != "" {
		if opt.ConfigPath == os.DevNull {
			return ""
		}
		return opt.ConfigPath
	} else if envVar := os.Getenv("ENGINE_CONFIG_PATH"); envVar != "" {
		return envVar
	}

	for _, path := range []string{
		os.ExpandEnv("${HOME}/.emailsender/emailsender.toml"),
		"/etc/emailsender/emailsender.toml",
	} {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}
