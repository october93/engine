package subcmd

import (
	"fmt"
	"os"

	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/worker/activityrecorder"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(recordCmd)
}

var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Records activities",
	Long:  "Starts recording all incoming activities",
	Run: func(cmd *cobra.Command, args []string) {
		logger, err := log.NewLogger(true, "info")
		if err != nil {
			exit(err)
		}
		ac := activityrecorder.NewActivityConsumer(&config.Worker, &config.Store, logger)
		err = ac.ConsumeJobs()
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
