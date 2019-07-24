// Package help is the help subcommand of the emailworker command.
package help

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Command displays help for command-line sub-commands.
type Command struct {
	Stdout io.Writer
}

// NewCommand returns a new instance of Command.
func NewCommand() *Command {
	return &Command{
		Stdout: os.Stdout,
	}
}

// Run executes the command.
func (cmd *Command) Run(args ...string) error {
	_, err := fmt.Fprintln(cmd.Stdout, strings.TrimSpace(usage))
	return err
}

const usage = `
Configure and start an emailworker.

Usage: emailworker [[command] [arguments]]

The commands are:

    config               display the default configuration
    help                 display this help message
    run                  run emailworker with existing configuration
    version              displays the emailworker version

"run" is the default command.

Use "emailworker [command] -help" for more information about a command.
`
