package subcmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/search"
	"github.com/october93/engine/store"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(indexCmd)
	RootCmd.AddCommand(indexUsersCmd)
	RootCmd.AddCommand(indexChannelsCmd)
}

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Indexes all objects",
	Long:  "Index all objects (users, channels) into the Algolia search engine",
	Run: func(cmd *cobra.Command, args []string) {
		if config == nil {
			exit(errors.New("missing configuration file, provide with --config [path]"))
		}
		if err := config.Validate(); err != nil {
			exit(err)
		}
		log, err := log.NewLogger(false, config.LogLevel)
		if err != nil {
			exit(err)
		}
		store, err := store.NewStore(&config.Store, log)

		if err != nil {
			exit(err)
		}
		indexer, err := search.NewIndexer(store, log, &config.Search)
		if err != nil {
			exit(err)
		}
		err = indexer.IndexAll(false)
		if err != nil {
			exit(err)
		}
	},
}

var indexUsersCmd = &cobra.Command{
	Use:   "indexusers",
	Short: "Indexes all users",
	Long:  "Index all channels into the Algolia search engine",
	Run: func(cmd *cobra.Command, args []string) {
		if config == nil {
			exit(errors.New("missing configuration file, provide with --config [path]"))
		}
		if err := config.Validate(); err != nil {
			exit(err)
		}
		log, err := log.NewLogger(false, config.LogLevel)
		if err != nil {
			exit(err)
		}
		store, err := store.NewStore(&config.Store, log)

		if err != nil {
			exit(err)
		}
		indexer, err := search.NewIndexer(store, log, &config.Search)
		if err != nil {
			exit(err)
		}
		err = indexer.IndexAllUsers(false)
		if err != nil {
			exit(err)
		}
	},
}

var indexChannelsCmd = &cobra.Command{
	Use:   "indexchannels",
	Short: "Indexes all channels",
	Long:  "Index all channels into the Algolia search engine",
	Run: func(cmd *cobra.Command, args []string) {
		if config == nil {
			exit(errors.New("missing configuration file, provide with --config [path]"))
		}
		if err := config.Validate(); err != nil {
			exit(err)
		}
		log, err := log.NewLogger(false, config.LogLevel)
		if err != nil {
			exit(err)
		}
		store, err := store.NewStore(&config.Store, log)

		if err != nil {
			exit(err)
		}
		indexer, err := search.NewIndexer(store, log, &config.Search)
		if err != nil {
			exit(err)
		}
		err = indexer.IndexAllUsers(false)
		if err != nil {
			exit(err)
		}
	},
}

func exit(err error) {
	if _, err = fmt.Fprintln(os.Stderr, err); err != nil {
		fmt.Println(err)
	}
	os.Exit(1)
}
