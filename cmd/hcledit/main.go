// The hcledit tool provides a CRUD interface to attributes within a Terraform
// file.
package main

import (
	"fmt"
	"os"

	"go.mercari.io/hcledit/cmd/hcledit/internal/command"
)

var version = "dev"

func main() {
	cmd := command.NewCmdRoot(version)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
