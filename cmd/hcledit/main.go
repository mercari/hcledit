package main

import (
	"os"

	"github.com/mercari/hcledit/cmd/hcledit/internal/command"
)

var version = "dev"

func main() {
	cmd := command.NewCmdRoot(version)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
