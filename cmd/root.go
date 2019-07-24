package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/october93/engine/kit/log"
	"github.com/spf13/cobra"
)

var (
	buildParameters BuildParameters
	cfgPath         string
)

var rootCmd = &cobra.Command{
	Use:   "engine",
	Short: "Engine powers the backend for the October app",
	Long:  "Engine powers the backend for the October app",
	Run: func(cmd *cobra.Command, args []string) {
		config := NewConfig()
		if err := config.Load(cfgPath); err != nil {
			exit(err)
		}
		// override server port if specified via environment variable
		if os.Getenv("PORT") != "" {
			port, err := strconv.Atoi(os.Getenv("PORT"))
			if err != nil {
				exit(err)
			}
			config.Server.Port = port
		}
		err := config.Validate()
		if err != nil {
			exit(err)
		}
		logger, err := log.NewLogger(isDeployed(buildParameters.Version), config.LogLevel)
		if err != nil {
			exit(err)
		}
		logger = logger.With("app", "engine")
		e := NewEngine(buildParameters, config, logger)
		err = e.Run()
		if err != nil {
			exit(err)
		}
		handleInterrupts(e, logger)
	},
}

func handleInterrupts(e *Engine, log log.Logger) {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	log.Info("Listening for signals")

	// Block until one of the signals above is received
	<-signalCh
	log.Info("Signal received, initializing clean shutdown...")
	go e.Close() //nolint

	// Block again until another signal is received, a shutdown timeout elapses,
	// or the Command is gracefully closed
	log.Info("Waiting for clean shutdown...")
	select {
	case <-signalCh:
		log.Info("second signal received, initializing hard shutdown")
	case <-time.After(time.Second * 30):
		log.Info("time limit reached, initializing hard shutdown")
	case <-e.Closed:
		log.Info("server shutdown completed")
	}
	// goodbye.
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgPath, "config", "", "config file (default is $HOME/config.toml)")
	rootCmd.AddCommand(cmdValidate)
}

type BuildParameters struct {
	Version string
	Branch  string
	Commit  string
}

func Execute(bp BuildParameters) {
	buildParameters = bp
	if err := rootCmd.Execute(); err != nil {
		exit(err)
	}
}

func exit(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func isDeployed(version string) bool {
	return version != "" && version != "unknown"
}
