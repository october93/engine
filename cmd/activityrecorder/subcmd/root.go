package subcmd

import (
	"fmt"
	"os"

	"github.com/october93/engine/cmd"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "activityrecorder",
	Short: "Writes API activities into the database",
	Long:  "User issued API calls are written as activities into the database.",
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
