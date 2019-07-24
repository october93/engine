package main

import (
	"math/rand"
	"time"

	"github.com/october93/engine/cmd"
)

const unknown = "unknown"

var (
	version string
	branch  string
	commit  string
)

func init() {
	if version == "" {
		version = unknown
	}
	if branch == "" {
		branch = unknown
	}
	if commit == "" {
		commit = unknown
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	cmd.Execute(cmd.BuildParameters{
		Version: version,
		Branch:  branch,
		Commit:  commit,
	})
}
