package command

import (
	"os"

	"github.com/mitchellh/cli"
)

func Run(args []string) int {
	return run(args, commands())
}

func run(args []string, commands map[string]cli.CommandFactory) int {
	cli := &cli.CLI{
		Name:     "hcledit",
		Args:     args,
		Commands: commands,
		Version:  "0.1.0",

		Autocomplete:               true,
		AutocompleteNoDefaultFlags: true,

		HelpFunc:   cli.BasicHelpFunc("hcledit"),
		HelpWriter: os.Stderr,
	}

	exitcode, _ := cli.Run()
	return exitcode
}

func commands() map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"read": func() (cli.Command, error) {
			return &readCommand{}, nil
		},
	}
}
