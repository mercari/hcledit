package main

import (
	"fmt"
	"os"

	"go.mercari.io/hcledit/cmd/hcledit/internal/command"
)

func main() {
	cmd := command.NewCmdRoot()

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
