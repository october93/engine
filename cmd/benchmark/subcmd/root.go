package subcmd

import (
	"fmt"
	"os"
	"time"

	"github.com/october93/engine/rpc"
	"github.com/october93/engine/rpc/client"
	"github.com/october93/engine/kit/log"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "Benchmark is a tool to measure Engine's performance",
	Long:  "Benchmark is a tool to measure and find the limitation of Engine's performance.",
}

var (
	endpoint string
	users    int
	skip     int
)

const password = "secret"

func init() {
	RootCmd.PersistentFlags().StringVar(&endpoint, "endpoint", "", "endpoint to connect to")
	RootCmd.PersistentFlags().IntVar(&users, "users", 100, "number of users to create")
	RootCmd.PersistentFlags().IntVar(&skip, "skip", 0, "number of users to skip")
}

func newClient(endpint string, l log.Logger) (rpc.RPC, error) {
	config := client.NewConfig()
	config.Address = endpoint
	config.Timeout = 30 * time.Minute
	cl, err := client.NewClient(config, l)
	return cl, err
}

func exit(err error) {
	if _, err = fmt.Fprintln(os.Stderr, err); err != nil {
		fmt.Println(err)
	}
	os.Exit(1)
}

func username(i int) string {
	return fmt.Sprintf("user-%d", i+1)
}
