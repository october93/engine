package cmd

import (
	"fmt"

	"github.com/october93/engine/coinmanager"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/rpc/notifications"
	"github.com/october93/engine/rpc/protocol"
	"github.com/october93/engine/rpc/push"
	"github.com/october93/engine/store"
	"github.com/october93/engine/worker"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(coinUpdaterCmd)
	rootCmd.AddCommand(coinUpdaterPopularPostCmd)
}

type noConnections struct {
}

func (c *noConnections) WritersByUser(globalid.ID) []*protocol.PushWriter {
	return []*protocol.PushWriter{}
}

func (c *noConnections) Writers() []*protocol.PushWriter {
	return []*protocol.PushWriter{}
}

var coinUpdaterCmd = &cobra.Command{
	Use:   "coin-updater",
	Short: "Updates the token rewards and leaderboard position",
	Long:  `Updates the token rewards and leaderboard position`,
	Run: func(cmd *cobra.Command, args []string) {
		config := NewConfig()
		if err := config.Load(cfgPath); err != nil {
			exit(err)
		}
		err := config.Validate()
		if err != nil {
			exit(err)
		}
		log, err := log.NewLogger(isDeployed(buildParameters.Version), config.LogLevel)
		if err != nil {
			exit(err)
		}
		log = log.With("app", "coinupdater")
		store, err := store.NewStore(&config.Store, log)
		if err != nil {
			exit(err)
		}

		// create notifications builder
		imagePath := fmt.Sprintf("https://s3-%s.amazonaws.com/%s/%s", config.Server.RPC.S3Region, config.Server.RPC.S3Bucket, config.Server.RPC.SystemIconPath)
		ns := notifications.NewNotifications(store, imagePath, config.Server.RPC.UnitsPerCoin)

		nw, err := worker.NewNotifier(store, &config.Worker, log)
		if err != nil {
			exit(err)
		}
		p, err := push.NewPusher(&noConnections{}, store, &config.Push, log)
		if err != nil {
			exit(err)
		}

		// Coin manager
		cm := coinmanager.NewCoinManager(store, &config.CoinManager)

		coinUpdater := worker.NewCoinUpdater(store, log, ns, nw, p, config.Worker.CoinUpdater, cm)
		coinUpdater.Run()
	},
}

var coinUpdaterPopularPostCmd = &cobra.Command{
	Use:   "coin-updater-popular-posts",
	Short: "Finds and rewards users for popular posts",
	Long:  `Finds and rewards users for popular posts`,
	Run: func(cmd *cobra.Command, args []string) {
		config := NewConfig()
		if err := config.Load(cfgPath); err != nil {
			exit(err)
		}
		err := config.Validate()
		if err != nil {
			exit(err)
		}
		log, err := log.NewLogger(isDeployed(buildParameters.Version), config.LogLevel)
		if err != nil {
			exit(err)
		}
		log = log.With("app", "coinupdater")
		store, err := store.NewStore(&config.Store, log)
		if err != nil {
			exit(err)
		}

		// create notifications builder
		imagePath := fmt.Sprintf("https://s3-%s.amazonaws.com/%s/%s", config.Server.RPC.S3Region, config.Server.RPC.S3Bucket, config.Server.RPC.SystemIconPath)
		ns := notifications.NewNotifications(store, imagePath, config.Server.RPC.UnitsPerCoin)

		nw, err := worker.NewNotifier(store, &config.Worker, log)
		if err != nil {
			exit(err)
		}
		p, err := push.NewPusher(&noConnections{}, store, &config.Push, log)
		if err != nil {
			exit(err)
		}

		// Coin manager
		cm := coinmanager.NewCoinManager(store, &config.CoinManager)

		coinUpdater := worker.NewCoinUpdater(store, log, ns, nw, p, config.Worker.CoinUpdater, cm)
		coinUpdater.RunPopularPostUpdates()
	},
}
