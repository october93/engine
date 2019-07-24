package run

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

// PrintConfigCommand represents the command executed by "emailworker config".
type PrintConfigCommand struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// NewPrintConfigCommand return a new instance of PrintConfigCommand.
func NewPrintConfigCommand() *PrintConfigCommand {
	return &PrintConfigCommand{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// Run parses and prints the current config loaded.
func (cmd *PrintConfigCommand) Run(args ...string) error {
	// Parse command flags.
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	configPath := fs.String("config", "", "")
	fs.Usage = func() {
		if _, err := fmt.Fprintln(os.Stderr, printConfigUsage); err != nil {
			fmt.Println(err)
		}
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	// Parse config from path.
	opt := Options{ConfigPath: *configPath}
	config, err := cmd.parseConfig(opt.GetConfigPath())
	if err != nil {
		return fmt.Errorf("parse config: %s", err)
	}

	// Validate the configuration.
	if err = config.Validate(); err != nil {
		return fmt.Errorf("%s. To generate a valid configuration file run `emailworker config > emailworker.generated.toml`", err)
	}

	if err = toml.NewEncoder(cmd.Stdout).Encode(config); err != nil {
		return fmt.Errorf("error encoding toml: %s", err)
	}
	_, err = fmt.Fprint(cmd.Stdout, "\n")
	return err
}

// ParseConfig parses the config at path.
// Returns a demo configuration if path is blank.
func (cmd *PrintConfigCommand) parseConfig(path string) (*Config, error) {
	config, err := NewDefaultConfig()
	if err != nil {
		config = NewConfig()
	}

	if path == "" {
		return config, nil
	}

	if _, err = fmt.Fprintf(os.Stderr, "Merging with configuration at: %s\n", path); err != nil {
		fmt.Println(err)
	}

	if err := config.Load(path); err != nil {
		return nil, err
	}
	return config, nil
}

//TODO: (corylanou) confirm these will be the correct locations
var printConfigUsage = `Displays the default configuration.

Usage: emailworker config [flags]

    -config <path>
            Set the path to the initial configuration file.
            This defaults to the environment variable ENGINE_CONFIG_PATH,
            ~/.emailworker/emailworker.toml, or /etc/emailworker/emailworker.toml if a file
            is present at any of these locations.
            Disable the automatic loading of a configuration file using
            the null device (such as /dev/null).
`
