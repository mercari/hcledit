package main

import (
	"fmt"
	"os"

	"github.com/mercari/hcledit/cmd/hcledit/internal/command"
)

var version = "dev"

func main() {
	cmd := command.NewCmdRoot(version)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
