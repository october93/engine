package subcmd

import (
	"fmt"
	"os"

	"github.com/october93/engine/cmd"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "indexer",
	Short: "Indexer is the search engine utility tool",
	Long:  "Indexer is the search engine utility tool to interact with Engine's search engine.",
}

var (
	config     *cmd.Config
	configFile string
)

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&configFile, "config", "", "configuration file")
}

func initConfig() {
	// if configuration file is not set, ignore for now
	if configFile == "" {
		return
	}
	config = cmd.NewConfig()
	err := config.Load(configFile)
	if err != nil {
		fmt.Println("Can't read configuration:", err)
		os.Exit(1)
	}
}
