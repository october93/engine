package cmd

import (
	"github.com/spf13/cobra"
)

var cmdValidate = &cobra.Command{
	Use:   "validate",
	Short: "Validates the configuration file",
	Long:  `Validates the configuration file`,
	Run: func(cmd *cobra.Command, args []string) {
		config := NewConfig()
		if err := config.Load(cfgPath); err != nil {
			exit(err)
		}
		err := config.Validate()
		if err != nil {
			exit(err)
		}
	},
}
