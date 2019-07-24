package subcmd

import (
	"os"

	"github.com/BurntSushi/toml"
	enginecmd "github.com/october93/engine/cmd"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Display the default configuration",
	Long:  "Displays the default configuration.",
	Run: func(cmd *cobra.Command, args []string) {
		enc := toml.NewEncoder(os.Stdout)
		err := enc.Encode(enginecmd.NewConfig())
		if err != nil {
			exit(err)
		}
	},
}
