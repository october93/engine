package subcmd

import (
	"context"
	"fmt"
	"time"

	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/rpc"
	"github.com/october93/engine/rpc/client"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(connectionsCmd)
}

var connectionsCmd = &cobra.Command{
	Use:   "connections",
	Short: "Opens n WebSocket connections to the endpoint",
	Long:  "Stress tests the endpoint by opening n WebSocket connections.",
	Run: func(cmd *cobra.Command, args []string) {
		config := client.NewConfig()
		config.Address = endpoint
		l, err := log.NewLogger(false, log.Info)
		if err != nil {
			exit(err)
		}
		for i := 0; i < users; i++ {
			c, err := client.NewClient(config, l)
			if err != nil {
				exit(err)
			}
			req := rpc.GetCardsRequest{
				Params: rpc.GetCardsParams{
					PerPage: 10,
					Page:    2,
				},
			}
			_, err = c.GetCards(context.Background(), req)
			if err != nil {
				l.Error(err)
			}
		}
		l.Info(fmt.Sprintf("Opened all %d connections", users))
		time.Sleep(5 * time.Minute)
	},
}

type Connections struct {
	log   log.Logger
	users int
}

func NewConncetions(users int, log log.Logger) *Connections {
	return &Connections{
		log:   log,
		users: users,
	}
}

func (c *Connections) Run() {
	for i := 0; i < 100; i++ {
		go c.load(i % c.users)
	}
}

func (c *Connections) load(i int) {

}
