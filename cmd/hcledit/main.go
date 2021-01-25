package main

import (
	"os"

	"github.com/mercari/hcledit/cmd/hcledit/internal/command"
)

func main() {
	os.Exit(command.Run(os.Args[1:]))
}
