package main

import (
	"fmt"
	"os"

	"github.com/october93/engine/cmd/activityrecorder/subcmd"
)

func main() {
	if err := subcmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
